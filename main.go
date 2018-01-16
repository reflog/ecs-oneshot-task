package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"strings"
	"time"
)

func logErrorAndFail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mapTasks(vs []*ecs.Task, f func(*ecs.Task) *string) []*string {
	vsm := make([]*string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

var (
	debug        = kingpin.Flag("debug", "Enable debug mode.").Bool()
	cluster      = kingpin.Flag("cluster", "Name of the ECS cluster").Short('c').Required().String()
	waitDuration = kingpin.Flag("wait", "How long to wait for task to finish").Short('t').Default("5m").Duration()
	taskName     = kingpin.Flag("task-name", "Name of the task to create in the cluster").Default("oneshot").Short('n').String()
	taskJson     = kingpin.Flag("task-json", "JSON file with task definition describing the container running the task").Required().Short('j').File()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	logLevel := aws.LogLevel(aws.LogOff)
	var err error

	if *debug {
		logLevel = aws.LogLevel(aws.LogDebugWithRequestErrors | aws.LogDebugWithHTTPBody)
	}

	sess, err := session.NewSession(&aws.Config{Credentials: credentials.NewEnvCredentials(), LogLevel: logLevel})
	logErrorAndFail(err)

	svc := ecs.New(sess)
	maxResults := int64(1)
	taskList, err := svc.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(*taskName),
		MaxResults:   &maxResults,
		Sort:         aws.String("DESC"),
		Status:       aws.String("ACTIVE"),
	})

	logErrorAndFail(err)

	if len(taskList.TaskDefinitionArns) > 0 {
		log.Fatalf("A task family with such name (%v) already exists! Sorry :(", *taskName)
	}

	def := ecs.RegisterTaskDefinitionInput{}

	buf := new(bytes.Buffer)
	buf.ReadFrom(*taskJson)
	err = json.Unmarshal(buf.Bytes(), &def)
	logErrorAndFail(err)
	def.Family = taskName

	log.Printf("Registering task definition...")
	registrationResult, err := svc.RegisterTaskDefinition(&def)
	logErrorAndFail(err)

	log.Printf("Registered %v successfully", *taskName)

	taskArn := aws.String(fmt.Sprintf("%s:%d", *taskName, *registrationResult.TaskDefinition.Revision))

	defer func() {
		_, err = svc.DeregisterTaskDefinition(&ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: taskArn,
		})
		logErrorAndFail(err)
		log.Printf("DeRegistered %v successfully", taskArn)
	}()

	clusterName := aws.String(*cluster)
	result, err := svc.RunTask(&ecs.RunTaskInput{
		Cluster:        clusterName,
		TaskDefinition: taskArn,
	})
	logErrorAndFail(err)

	failureCount := len(result.Failures)
	if failureCount > 0 {
		log.Printf("Failed to run task. %d failures detected:", failureCount)
		for i := 0; i < failureCount; i++ {
			log.Print(result.Failures[i])
		}
	} else {
		log.Printf("Waiting for task %v to finish...", taskArn)
		startWaitTime := time.Now()
		taskArns := mapTasks(result.Tasks, func(v *ecs.Task) *string { return v.TaskArn })
		waitParam := &ecs.DescribeTasksInput{
			Cluster: clusterName,
			Tasks:   taskArns,
		}

		for {
			if err = svc.WaitUntilTasksStopped(waitParam); err == nil {
				break
			}
			waitedTime := time.Now().Sub(startWaitTime)
			if waitedTime > *waitDuration {
				log.Fatalf("Aborting due to time out, task still running after %s or another error: %v", shortDur(waitedTime), err)
			}
		}

		log.Printf("Done! Removing the task definition...")
	}

}

# ecs-oneshot-task

Tiny Go CLI Program to run a one-shot (single execution) script on ECS. Think of it as a tiny AWS Lambda without time or memory constraints.

[![Build Status](https://travis-ci.org/reflog/ecs-oneshot-task.svg?branch=master)](https://travis-ci.org/reflog/ecs-oneshot-task)

# usage

```
usage: ecs-oneshot-task --cluster=CLUSTER --task-json=TASK-JSON [<flags>]

Flags:
      --help                   Show context-sensitive help (also try --help-long
                               and --help-man).
      --debug                  Enable debug mode.
  -c, --cluster=CLUSTER        Name of the ECS cluster
  -t, --wait=5m                How long to wait for task to finish
  -n, --task-name="oneshot"    Name of the task to create in the cluster
  -j, --task-json=TASK-JSON    JSON file with task definition describing the
                               container running the task
      --pass-aws-keys          Add AWS keys to task's environment.
  -p, --params=PARAMS ...      Parameter that can be used inside the JSON file
                               using Go templating
  -k, --aws-access-key-id=AWS-ACCESS-KEY-ID
                               AWS Access Key ID to use (overrides environment)
  -s, --aws-secret-key=AWS-SECRET-KEY
                               AWS Secret Access Key to use (overrides
                               environment)
  -r, --aws-region=AWS-REGION  AWS Region to user (overrides environment)
      --version                Show application version.

```

# usage example

In the `container` folder you can see a sample Docker definition for a container that can fetch a script from S3 on start and then execute it.
Using this program, we supply a ECS task definition JSON (i.e. what container to take, how much memory to give it and what environment variables to passO and the name of the cluster to run on and the task will be registered, executed, will wait until it's over and removed from definitions.

# tips

You can utilize `logConfiguration` inside task JSON to add logging. See AWS ECS documentation
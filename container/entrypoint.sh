#!/bin/sh
ret=`aws s3 ls $S3_PATH | wc -l`

if [ "$ret" = "0" ]; then
    echo "script file at path $S3_PATH not found!"
else
    echo "downloading script file at path $S3_PATH..."
    aws s3 cp $S3_PATH /tmp/oneshot.sh
    chmod +x /tmp/oneshot.sh
    echo "executing..."
    . /tmp/oneshot.sh
fi

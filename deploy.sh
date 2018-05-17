#!/usr/bin/env bash -e

if [ "$#" -ne 3 ]; then
    echo "Illegal number of parameters, usage: ./deploy.sh DEPLOYMENT_BUCKET DIST_FILE DEST_FUNCTION"
    echo "example: ./deploy.sh deployment-bucket-s3bucket-rkjl6q60hsw8 deployment-201805171544.zip DatadogLogs"
    exit 1
fi

BUCKET=${1}
FILE=${2}
FUNCTION_NAME=${3}

aws s3 cp ${FILE} s3://${BUCKET}/log-group-subscriber/ --profile cloudformation@flowlabdev

aws cloudformation deploy \
   --stack-name log-group-subscriber \
   --template-file cloudformation/template.yml \
   --parameter-overrides FunctionName=${FUNCTION_NAME} DeploymentBucket=${BUCKET} DeploymentFile=${FILE} \
   --capabilities CAPABILITY_IAM \
   --profile cloudformation@flowlabdev

## I am in progress ...

Build:
```sh
./build.sh
```

Deploy:
```sh
aws s3 cp deployment.zip s3://deployment-bucket-s3bucket-../log-group-subscriber/ --profile cloudformation@flowlabdev
aws cloudformation deploy \
    --stack-name log-group-subscriber \
    --template-file cloudformation/template.yml \
    --parameter-overrides FunctionName=DatadogLogs \
    --parameter-overrides DeploymentBucket=deployment-bucket-s3bucket-.. \
    --capabilities CAPABILITY_IAM \
    --profile cloudformation@flowlabdev
```
## AWS log group subscriber

Lambda function that creates subscription filters for other lambda function.

```


                         ---------------------------------------
   CloudWatch Event     |   log-subscriber-function             |
   (every 5 minutes)    |   1. Get all log groups               |
----------------------> |   2. Group all missing subscriptions  |
                        |   3. Create subscription filters      |
                        |                                       |
                         ---------------------------------------
```


To build run commend below. It compiles sources to `main` binary file and zips
it to deployment package `deployment.zip`
```sh
./build.sh
```

To deploy to AWS you need to provide two parameters:
1. is deployment s3 bucket name
2. Lambda function name that subscriptions will be created for

Upload `deployment.zip` to s3
```sh
aws s3 cp deployment.zip s3://deployment-bucket-s3bucket-../log-group-subscriber/ --profile cloudformation@flowlabdev
```

To deploy cloudformation run:
```sh
aws cloudformation deploy \
   --stack-name log-group-subscriber \
   --template-file cloudformation/template.yml \
   --parameter-overrides FunctionName=DatadogLogs DeploymentBucket=deployment-bucket-s3bucket-... \
   --capabilities CAPABILITY_IAM \
   --profile cloudformation@flowlabdev
```
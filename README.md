## AWS log group subscriber [![Build Status](https://travis-ci.org/flow-lab/log-group-subscriber.svg?branch=master)](https://travis-ci.org/flow-lab/log-group-subscriber)

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
it to deployment package `deployment-123456789.zip`
```sh
./build.sh
```

To deploy to AWS with cloudformation template use `deploy.sh` script
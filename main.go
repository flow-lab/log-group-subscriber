package main

import (
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"strings"
)

type LogGroup struct {
	LogGroupName *string
	FunctionArn  *string
}

func main() {
	lambda.Start(Handler)
}

func Handler(request events.CloudWatchEvent) (string, error) {
	log.Printf("Processing Lambda request %s\n", request)
	var functionArn = os.Getenv("DESTINATION_FUNCTION_ARN")

	sess := session.Must(session.NewSession())
	client := cloudwatchlogs.New(sess, &aws.Config{
		Region: aws.String(endpoints.EuWest1RegionID),
	})

	logGroups, err := GetLogGroups(client)
	check(err)

	missingSubscription := getLogGroupsWithMissingSubscription(logGroups.LogGroups, &functionArn, client)
	PutSubscriptionFilter(missingSubscription, client)

	return "Event processed", nil
}

func getLogGroupsWithMissingSubscription(groups []*cloudwatchlogs.LogGroup, functionArn *string, logs cloudwatchlogsiface.CloudWatchLogsAPI) []LogGroup {
	var result []LogGroup
	for _, element := range groups {
		if hasSubscriptionFilter(element, functionArn, logs) == false {
			l := LogGroup{
				LogGroupName: element.LogGroupName,
				FunctionArn:  functionArn,
			}
			log.Printf("%s is missing subscription for %s", *l.LogGroupName, *functionArn)
			result = append(result, l)
		}
	}
	return result
}

func hasSubscriptionFilter(logGroup *cloudwatchlogs.LogGroup, functionArn *string, logs cloudwatchlogsiface.CloudWatchLogsAPI) bool {
	var exists bool
	subscriptionFilters, err := DescribeSubscriptionFilters(logGroup.LogGroupName, logs)
	check(err)
	for _, subsFilter := range subscriptionFilters.SubscriptionFilters {
		if *functionArn == *subsFilter.DestinationArn {
			exists = true
			break
		}
	}
	return exists
}

func GetLogGroups(logs cloudwatchlogsiface.CloudWatchLogsAPI) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	input := cloudwatchlogs.DescribeLogGroupsInput{}
	return logs.DescribeLogGroups(&input)
}

func DescribeSubscriptionFilters(logGroupName *string, logs cloudwatchlogsiface.CloudWatchLogsAPI) (*cloudwatchlogs.DescribeSubscriptionFiltersOutput, error) {
	input := cloudwatchlogs.DescribeSubscriptionFiltersInput{
		LogGroupName: logGroupName,
	}
	return logs.DescribeSubscriptionFilters(&input)
}

func PutSubscriptionFilter(logGroups []LogGroup, logs cloudwatchlogsiface.CloudWatchLogsAPI) {
	level := ""
	for _, logGroup := range logGroups {
		filterName := fmt.Sprintf("%s-%s", strings.Replace(*logGroup.LogGroupName, "/", "", -1), "SubscriptionFilter")
		input := cloudwatchlogs.PutSubscriptionFilterInput{
			FilterName:     &filterName,
			LogGroupName:   logGroup.LogGroupName,
			DestinationArn: logGroup.FunctionArn,
			FilterPattern:  &level,
		}
		log.Printf("About to put subscription filter %s for %s", filterName, *logGroup.LogGroupName)
		_, err := logs.PutSubscriptionFilter(&input)
		check(err)

		fmt.Printf("PutSubscriptionFilter for %s", *logGroup.LogGroupName)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

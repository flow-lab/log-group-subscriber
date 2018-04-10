package main

import (
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/aws"
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
	log.Printf("processing event %s", request)
	var functionArn = os.Getenv("DESTINATION_FUNCTION_ARN")

	sess := session.Must(session.NewSession())
	client := cloudwatchlogs.New(sess, &aws.Config{})

	_, err := ProcessEvent(functionArn, client)
	if err != nil {
		panic(fmt.Errorf("unable to complete: %v", err))
	}

	return "event processed", nil
}

func ProcessEvent(functionArn string, logs cloudwatchlogsiface.CloudWatchLogsAPI) ([]string, error) {
	logGroups, err := GetLogGroups(logs)
	if err != nil {
		return nil, fmt.Errorf("get log groups: %v", err)
	}

	missingSubscription, err := getLogGroupsWithMissingSubscription(logGroups, &functionArn, logs)
	if err != nil {
		return nil, fmt.Errorf("get log with missing subscriptions: %v", err)
	}

	result, err := PutSubscriptionFilter(missingSubscription, logs)
	if err != nil {
		return nil, fmt.Errorf("get log with missing subscriptions: %v", err)
	}

	return result, nil
}

func GetLogGroups(logs cloudwatchlogsiface.CloudWatchLogsAPI) ([]*cloudwatchlogs.LogGroup, error) {
	var logGroups []*cloudwatchlogs.LogGroup
	input := cloudwatchlogs.DescribeLogGroupsInput{}
	err := logs.DescribeLogGroupsPages(&input, func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
		for _, logGroup := range page.LogGroups {
			logGroups = append(logGroups, logGroup)
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("describe log Groups: %v", err)
	}
	return logGroups, nil
}

func DescribeSubscriptionFilters(logGroupName *string, logs cloudwatchlogsiface.CloudWatchLogsAPI) (*cloudwatchlogs.DescribeSubscriptionFiltersOutput, error) {
	input := cloudwatchlogs.DescribeSubscriptionFiltersInput{
		LogGroupName: logGroupName,
	}
	return logs.DescribeSubscriptionFilters(&input)
}

func PutSubscriptionFilter(logGroups []LogGroup, logs cloudwatchlogsiface.CloudWatchLogsAPI) ([]string, error) {
	var result []string
	level := ""
	for _, logGroup := range logGroups {
		filterName := fmt.Sprintf("%s-%s", strings.Replace(*logGroup.LogGroupName, "/", "", -1), "SubscriptionFilter")
		input := cloudwatchlogs.PutSubscriptionFilterInput{
			FilterName:     &filterName,
			LogGroupName:   logGroup.LogGroupName,
			DestinationArn: logGroup.FunctionArn,
			FilterPattern:  &level,
		}
		log.Printf("put subscription filter %s for %s", filterName, *logGroup.LogGroupName)
		_, err := logs.PutSubscriptionFilter(&input)
		if err != nil {
			return result, fmt.Errorf("putSubscriptionFilter for %s: %v", *logGroup.LogGroupName, err)
		}
		result = append(result, *logGroup.LogGroupName)
		fmt.Printf("PutSubscriptionFilter for %s", *logGroup.LogGroupName)
	}
	return result, nil
}

func getLogGroupsWithMissingSubscription(groups []*cloudwatchlogs.LogGroup, functionArn *string, logs cloudwatchlogsiface.CloudWatchLogsAPI) ([]LogGroup, error) {
	var result []LogGroup
	for _, element := range groups {
		hasSubscriptionFilter, err := hasSubscriptionFilter(element, functionArn, logs)
		if err != nil {
			return nil, fmt.Errorf("getLogGroupsWithMissingSubscription: %v", err)
		}
		if hasSubscriptionFilter == false {
			logGroup := LogGroup{
				LogGroupName: element.LogGroupName,
				FunctionArn:  functionArn,
			}
			log.Printf("%s is missing subscription for %s", *logGroup.LogGroupName, *functionArn)
			result = append(result, logGroup)
		}
	}
	return result, nil
}

func hasSubscriptionFilter(logGroup *cloudwatchlogs.LogGroup, functionArn *string, logs cloudwatchlogsiface.CloudWatchLogsAPI) (bool, error) {
	subscriptionFilters, err := DescribeSubscriptionFilters(logGroup.LogGroupName, logs)
	if err != nil {
		return false, fmt.Errorf("describe subscription filters: %s", err)
	}
	for _, subsFilter := range subscriptionFilters.SubscriptionFilters {
		if *functionArn == *subsFilter.DestinationArn {
			return true, nil
		}
	}
	return false, nil
}

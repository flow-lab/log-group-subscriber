package main

import (
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

func Handler(request events.CloudWatchEvent) (string, error) {
	log.Printf("Processing Lambda request %s\n", request)

	sess := session.Must(session.NewSession())

	client := cloudwatchlogs.New(sess)

	logGroups, err := GetLogGroups(client)
	check(err)

	subscriptionFilters, err := DescribeSubscriptionFilters(client)
	check(err)

	fmt.Print(logGroups, subscriptionFilters)

	// TODO [grokrz]: impl

	return "Event processed", nil
}

func GetLogGroups(logs cloudwatchlogsiface.CloudWatchLogsAPI) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	input := cloudwatchlogs.DescribeLogGroupsInput{}
	return logs.DescribeLogGroups(&input)
}

func DescribeSubscriptionFilters(logs cloudwatchlogsiface.CloudWatchLogsAPI) (*cloudwatchlogs.DescribeSubscriptionFiltersOutput, error) {
	input := cloudwatchlogs.DescribeSubscriptionFiltersInput{}
	return logs.DescribeSubscriptionFilters(&input)
}

func PutSubscriptionFilter(logs cloudwatchlogsiface.CloudWatchLogsAPI) (*cloudwatchlogs.PutSubscriptionFilterOutput, error) {
	input := cloudwatchlogs.PutSubscriptionFilterInput{}
	return logs.PutSubscriptionFilter(&input)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	lambda.Start(Handler)
}

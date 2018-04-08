package main_test

import (
	"testing"
	"io/ioutil"
	"log-group-subscriber"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/stretchr/testify/assert"
)

//func TestHandler(t *testing.T) {
//	var scheduledEvent events.CloudWatchEvent
//	var inputJson = readJsonFromFile(t, "testdata/scheduled-event.json")
//
//	var err = json.Unmarshal(inputJson, &scheduledEvent)
//	check(t, err)
//
//	tests := []struct {
//		request events.CloudWatchEvent
//		expect  string
//		err     error
//	}{
//		{
//			request: scheduledEvent,
//			expect:  "Event processed",
//			err:     nil,
//		},
//	}
//
//	for _, test := range tests {
//		response, err := main.Handler(test.request)
//		assert.IsType(t, test.err, err)
//		assert.Equal(t, test.expect, response)
//	}
//}

func TestGetLogGroups(t *testing.T) {
	cwl := &mockCloudWatchLogsClient{}

	logGroups, err := main.GetLogGroups(cwl)

	check(t, err)
	assert.NotNil(t, logGroups)
	assert.Equal(t, 2, len(logGroups.LogGroups))
}

func TestDescribeSubscriptionFilters(t *testing.T) {
	cwl := &mockCloudWatchLogsClient{}

	subscriptionFilters, err := main.DescribeSubscriptionFilters(cwl)

	check(t, err)
	assert.Equal(t, 1, len(subscriptionFilters.SubscriptionFilters))
}

func TestPutSubscriptionFilter(t *testing.T) {
	cwl := &mockCloudWatchLogsClient{}

	res, err := main.PutSubscriptionFilter(cwl)

	check(t, err)
	assert.NotNil(t, res)
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
		panic(err)
	}
}

// Define a mock struct to be used in your unit tests of myFunc.
type mockCloudWatchLogsClient struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
}

func (m *mockCloudWatchLogsClient) DescribeLogGroups(input *cloudwatchlogs.DescribeLogGroupsInput) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	var inputJson = readFile("testdata/describeLogGroups-output.json")
	var describeLogGroupsOutput cloudwatchlogs.DescribeLogGroupsOutput
	err := json.Unmarshal(inputJson, &describeLogGroupsOutput)
	if err != nil {
		panic(err)
	}
	return &describeLogGroupsOutput, nil
}

func (m *mockCloudWatchLogsClient) DescribeSubscriptionFilters(*cloudwatchlogs.DescribeSubscriptionFiltersInput) (*cloudwatchlogs.DescribeSubscriptionFiltersOutput, error) {
	var inputJson = readFile("testdata/describeSubscriptionFilters-output.json")
	var describeSubscriptionFiltersOutput cloudwatchlogs.DescribeSubscriptionFiltersOutput
	err := json.Unmarshal(inputJson, &describeSubscriptionFiltersOutput)
	if err != nil {
		panic(err)
	}
	return &describeSubscriptionFiltersOutput, nil
}

func (m *mockCloudWatchLogsClient) PutSubscriptionFilter(input *cloudwatchlogs.PutSubscriptionFilterInput) (*cloudwatchlogs.PutSubscriptionFilterOutput, error) {
	return &cloudwatchlogs.PutSubscriptionFilterOutput{}, nil
}

func readFile(path string) ([]byte) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return f
}

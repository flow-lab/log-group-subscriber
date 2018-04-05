package main_test

import (
	"testing"
	"io/ioutil"
	"log-group-subscriber"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestHandler(t *testing.T) {
	var scheduledEvent events.CloudWatchEvent
	var inputJson = readJsonFromFile(t, "testdata/scheduled-event.json")

	var err = json.Unmarshal(inputJson, &scheduledEvent)
	check(t, err)

	tests := []struct {
		request events.CloudWatchEvent
		expect  string
		err     error
	}{
		{
			request: scheduledEvent,
			expect:  "Event processed",
			err:     nil,
		},
	}

	for _, test := range tests {
		response, err := main.Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}
}

func readJsonFromFile(t *testing.T, inputFile string) ([]byte) {
	inputJson, err := ioutil.ReadFile(inputFile)
	check(t, err)
	return inputJson
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
		panic(err)
	}
}

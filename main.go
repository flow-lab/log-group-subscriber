package main

import (
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.CloudWatchEvent) (string, error) {
	log.Printf("Processing Lambda request %s\n", request)
	// TODO [grokrz]: implement
	return "Event processed", nil
}

func main() {
	lambda.Start(Handler)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/yunuskilicdev/serverlessExample/common"
	"github.com/yunuskilicdev/serverlessExample/common/model"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		var request model.SendVerificationMailRequest
		json.Unmarshal([]byte(message.Body), &request)
		common.SendMail(request.Token)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

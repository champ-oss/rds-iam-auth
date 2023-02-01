package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	log "github.com/sirupsen/logrus"
	"os"
)

// global state
func init() {}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	_ = cfg.LoadConfig()

	for _, message := range sqsEvent.Records {
		log.Infof("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}

	log.Info("starting main run..")

	// Get list of RDS and Aurora

	// Look up password in SSM

	// Login to endpoint

	// Run SQL to enable IAM Auth

	return nil
}

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		// Support running the code locally
		_ = handler(nil, events.SQSEvent{})
	} else {
		lambda.Start(handler)
	}
}

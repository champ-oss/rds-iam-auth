package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/service/scheduler"
	"github.com/champ-oss/rds-iam-auth/service/worker"
	"os"
)

var config *cfg.Config
var schedulerService *scheduler.Service
var runnerService *worker.Service

func init() {
	config = cfg.LoadConfig()
	schedulerService = scheduler.NewService(config)
	runnerService = worker.NewService(config)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	if len(sqsEvent.Records) < 1 {
		return schedulerService.Run()
	}

	for _, message := range sqsEvent.Records {
		if err := runnerService.Run(message); err != nil {
			return err
		}
	}

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

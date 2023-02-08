package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/rds_client"
	"github.com/champ-oss/rds-iam-auth/pkg/sqs_client"
	"github.com/champ-oss/rds-iam-auth/pkg/ssm_client"
	"github.com/champ-oss/rds-iam-auth/service/scheduler"
	"github.com/champ-oss/rds-iam-auth/service/worker"
	log "github.com/sirupsen/logrus"
	"os"
)

var config *cfg.Config
var schedulerService *scheduler.Service
var workerService *worker.Service

func init() {
	config = cfg.LoadConfig()
	rdsClient := rds_client.NewRdsClient(config)
	sqsClient := sqs_client.NewSqsClient(config)
	ssmClient := ssm_client.NewSqsClient(config)
	schedulerService = scheduler.NewService(config, rdsClient, sqsClient)
	workerService = worker.NewService(config, rdsClient, ssmClient)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	if len(sqsEvent.Records) < 1 {
		return schedulerService.Run()
	}

	for _, message := range sqsEvent.Records {
		log.Warning("triggered from sqs message")
		if err := workerService.Run(message, nil); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		// Support running the code locally
		_ = handler(context.TODO(), events.SQSEvent{
			Records: []events.SQSMessage{
				{
					Body: "cluster|rds-iam-auth-20230202151930094200000014",
				},
				{
					Body: "instance|rds-iam-auth",
				},
			},
		})
	} else {
		lambda.Start(handler)
	}
}

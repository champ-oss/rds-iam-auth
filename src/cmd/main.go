package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
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

func handler(ctx context.Context, event json.RawMessage) error {
	log.Debugf("event: %s", event)

	if common.IsScheduledEvent(event) {
		return schedulerService.Run()

	} else if isEventBridgeRdsEvent, cloudwatchEvent := common.IsEventBridgeRdsEvent(event); isEventBridgeRdsEvent {
		return workerService.Run(nil, &cloudwatchEvent, nil)

	} else if isSqsEvent, sqsEvent := common.IsSqsEvent(event); isSqsEvent {
		for _, message := range sqsEvent.Records {
			if err := workerService.Run(&message, nil, nil); err != nil {
				return err
			}
		}
		return nil

	} else {
		return fmt.Errorf("unable to recognize lambda event")
	}
}

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		// Support running the code locally
		if err := handler(context.TODO(), nil); err != nil {
			panic(err)
		}
	} else {
		lambda.Start(handler)
	}
}

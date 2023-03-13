package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/rds_client"
	"github.com/champ-oss/rds-iam-auth/pkg/sqs_client"
	"github.com/champ-oss/rds-iam-auth/pkg/ssm_client"
	"github.com/champ-oss/rds-iam-auth/service/scheduler"
	"github.com/champ-oss/rds-iam-auth/service/worker"
	"os"
)

var config *cfg.Config
var schedulerService *scheduler.Service
var workerService *worker.Service

type Event struct {
	Records   []string
	Detail    string
	Resources []string
	Region    string
}

func init() {
	config = cfg.LoadConfig()
	rdsClient := rds_client.NewRdsClient(config)
	sqsClient := sqs_client.NewSqsClient(config)
	ssmClient := ssm_client.NewSqsClient(config)
	schedulerService = scheduler.NewService(config, rdsClient, sqsClient)
	workerService = worker.NewService(config, rdsClient, ssmClient)
}

func handler(ctx context.Context, event json.RawMessage) error {

	fmt.Println(event)

	//if len(sqsEvent.Records) < 1 {
	//	return schedulerService.Run()
	//}

	//for _, message := range sqsEvent.Records {
	//	log.Warning("triggered from sqs message")
	//	if err := workerService.Run(message, nil); err != nil {
	//		return err
	//	}
	//}

	return nil
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

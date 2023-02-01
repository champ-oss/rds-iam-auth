package scheduler

import (
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/rds_client"
	"github.com/champ-oss/rds-iam-auth/pkg/sqs_client"
)

type Service struct {
	config    *cfg.Config
	sqsClient *sqs_client.SqsClient
	rdsClient *rds_client.RdsClient
}

func NewService(config *cfg.Config) *Service {
	return &Service{
		config:    config,
		sqsClient: sqs_client.NewSqsClient(config.AwsRegion, config.QueueUrl),
		rdsClient: rds_client.NewRdsClient(config.AwsRegion, config.QueueUrl),
	}
}

func (s *Service) Run() error {
	_ = s.rdsClient.GetAllDatabases()
	return nil
}

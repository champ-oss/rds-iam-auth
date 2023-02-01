package scheduler

import (
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/sqs_client"
)

type Service struct {
	config    *cfg.Config
	sqsClient *sqs_client.SqsClient
}

func NewService(config *cfg.Config) *Service {
	return &Service{
		config:    config,
		sqsClient: sqs_client.NewSqsClient(config.AwsRegion, config.QueueUrl),
	}
}

func (s *Service) Run() error {
	return nil
}

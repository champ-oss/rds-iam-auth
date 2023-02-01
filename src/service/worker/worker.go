package worker

import (
	"github.com/aws/aws-lambda-go/events"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/rds_client"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	config    *cfg.Config
	rdsClient *rds_client.RdsClient
}

func NewService(config *cfg.Config) *Service {
	return &Service{
		config:    config,
		rdsClient: rds_client.NewRdsClient(config.AwsRegion, config.QueueUrl),
	}
}

func (s *Service) Run(message events.SQSMessage) error {
	log.Infof("sqs message body: %s", message.Body)
	return nil
}

package scheduler

import (
	"fmt"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
	"github.com/champ-oss/rds-iam-auth/pkg/rds_client"
	"github.com/champ-oss/rds-iam-auth/pkg/sqs_client"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	config    *cfg.Config
	sqsClient sqs_client.SqsClientInterface
	rdsClient rds_client.RdsClientInterface
}

func NewService(config *cfg.Config, rdsClient rds_client.RdsClientInterface, sqsClient sqs_client.SqsClientInterface) *Service {
	return &Service{
		config:    config,
		sqsClient: sqsClient,
		rdsClient: rdsClient,
	}
}

func (s *Service) Run() error {
	for _, database := range s.rdsClient.GetAllDBClusters() {
		message := fmt.Sprintf("%s%s%s", common.RdsTypeClusterKey, common.SqsMessageBodySeparator, database)
		if err := s.sqsClient.Send(message); err != nil {
			log.Error(err)
			return err
		}
	}

	for _, database := range s.rdsClient.GetAllDBInstances() {
		message := fmt.Sprintf("%s%s%s", common.RdsTypeInstanceKey, common.SqsMessageBodySeparator, database)
		if err := s.sqsClient.Send(message); err != nil {
			log.Error(err)
			return err
		}
	}

	log.Info("all databases have been scheduled using SQS")
	return nil
}

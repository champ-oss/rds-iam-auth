package worker

import (
	"github.com/aws/aws-lambda-go/events"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
	"github.com/champ-oss/rds-iam-auth/pkg/rds_client"
	log "github.com/sirupsen/logrus"
	"strings"
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
	rdsType, rdsIdentifier := parseSqsMessageBody(message)

	switch rdsType {
	case common.RdsTypeClusterKey:
		log.Infof("getting RDS cluster information for: %s", rdsIdentifier)

	case common.RdsTypeInstanceKey:
		log.Infof("getting RDS instance information for: %s", rdsIdentifier)

	default:
		log.Fatalf("unrecognized RDS type: %s", rdsType)
	}

	return nil
}

func parseSqsMessageBody(message events.SQSMessage) (rdsType string, rdsIdentifier string) {
	messageParts := strings.Split(message.Body, common.SqsMessageBodySeparator)
	if len(messageParts) != 2 {
		log.Fatalf("unable to parse sqs message: %s", message.Body)
	}
	rdsType = messageParts[0]
	rdsIdentifier = messageParts[1]
	return rdsType, rdsIdentifier
}

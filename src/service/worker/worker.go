package worker

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
	"github.com/champ-oss/rds-iam-auth/pkg/rds_client"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Service struct {
	config    *cfg.Config
	rdsClient rds_client.RdsClientInterface
}

func NewService(config *cfg.Config, rdsClient rds_client.RdsClientInterface) *Service {
	return &Service{
		config:    config,
		rdsClient: rdsClient,
	}
}

func (s *Service) Run(message events.SQSMessage) error {
	rdsType, rdsIdentifier, err := parseSqsMessage(message)
	if err != nil {
		return err
	}

	switch rdsType {
	case common.RdsTypeClusterKey:
		log.Infof("getting RDS cluster information for: %s", rdsIdentifier)

	case common.RdsTypeInstanceKey:
		log.Infof("getting RDS instance information for: %s", rdsIdentifier)

	default:
		return fmt.Errorf("unrecognized RDS type: %s", rdsType)
	}

	return nil
}

func parseSqsMessage(message events.SQSMessage) (rdsType string, rdsIdentifier string, err error) {
	log.Debugf("sqs message body: %s", message.Body)
	messageParts := strings.Split(message.Body, common.SqsMessageBodySeparator)
	if len(messageParts) != 2 {
		return "", "", fmt.Errorf("unable to parse sqs message: %s", message.Body)
	}
	rdsType = messageParts[0]
	rdsIdentifier = messageParts[1]
	return rdsType, rdsIdentifier, nil
}

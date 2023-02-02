package ssm_client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	log "github.com/sirupsen/logrus"
)

type SsmClientInterface interface {
	GetValue(name string) (string, error)
}

type SsmClient struct {
	queueUrl  string
	ssmClient *ssm.Client
}

func NewSqsClient(config *cfg.Config) *SsmClient {
	return &SsmClient{
		queueUrl:  config.QueueUrl,
		ssmClient: ssm.NewFromConfig(config.AwsConfig),
	}
}

func (s *SsmClient) GetValue(name string) (string, error) {
	log.Debugf("getting value from ssm parameter: %s", name)
	output, err := s.ssmClient.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *output.Parameter.Value, nil
}

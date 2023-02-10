package ssm_client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	log "github.com/sirupsen/logrus"
)

type SsmClientInterface interface {
	GetValue(name string) (string, error)
	SearchByName(name string) ([]string, error)
	SearchByTag(tagName, tagValue string) ([]string, error)
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

// SearchByName searches SSM for a parameter containing the provided name
func (s *SsmClient) SearchByName(name string) ([]string, error) {
	log.Debugf("searching ssm for %s", name)
	output, err := s.ssmClient.DescribeParameters(context.TODO(), &ssm.DescribeParametersInput{
		ParameterFilters: []types.ParameterStringFilter{
			{
				Key:    aws.String("Name"),
				Option: aws.String("Contains"),
				Values: []string{name},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var results []string
	for _, param := range output.Parameters {
		results = append(results, *param.Name)
	}
	return results, nil
}

// SearchByTag searches SSM for a parameter matching the provided tag key and value
func (s *SsmClient) SearchByTag(tagName, tagValue string) ([]string, error) {
	log.Debugf("searching ssm for tag %s=%s", tagName, tagValue)
	output, err := s.ssmClient.DescribeParameters(context.TODO(), &ssm.DescribeParametersInput{
		ParameterFilters: []types.ParameterStringFilter{
			{
				Key:    aws.String("tag:" + tagName),
				Option: aws.String("Equals"),
				Values: []string{tagValue},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var results []string
	for _, param := range output.Parameters {
		results = append(results, *param.Name)
	}
	return results, nil
}

// GetValue gets the decrypted value of an SSM parameter
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

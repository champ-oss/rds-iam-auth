package sqs_client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	log "github.com/sirupsen/logrus"
)

type SqsClientInterface interface {
	Send(messageBody string) error
}

type SqsClient struct {
	queueUrl  string
	sqsClient *sqs.Client
}

func NewSqsClient(config *cfg.Config) *SqsClient {
	return &SqsClient{
		queueUrl:  config.QueueUrl,
		sqsClient: sqs.NewFromConfig(config.AwsConfig),
	}
}

func (s *SqsClient) Send(messageBody string) error {
	log.Debugf("sending message: '%s' to queue: '%s'", messageBody, s.queueUrl)
	_, err := s.sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    aws.String(s.queueUrl),
	})
	return err
}

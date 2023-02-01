package sqs_client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
	log "github.com/sirupsen/logrus"
)

type SqsClient struct {
	queueUrl  string
	sqsClient *sqs.Client
}

func NewSqsClient(region string, queueUrl string) *SqsClient {
	return &SqsClient{
		queueUrl:  queueUrl,
		sqsClient: sqs.NewFromConfig(common.GetAWSConfig(region)),
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

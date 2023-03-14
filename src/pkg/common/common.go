package common

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	log "github.com/sirupsen/logrus"
	"strings"
)

const SqsMessageBodySeparator = "|"
const RdsTypeClusterKey = "cluster"
const RdsTypeInstanceKey = "instance"

type MySQLConnectionInfo struct {
	Endpoint       string
	Port           int32
	Username       string
	Password       string
	Database       string
	SecurityGroups []string
}

// ParseSqsMessage parses the RDS type and RDS identifier from the incoming SQS message body
func ParseSqsMessage(message events.SQSMessage) (rdsType string, rdsIdentifier string, err error) {
	log.Debugf("sqs message body: %s", message.Body)
	messageParts := strings.Split(message.Body, SqsMessageBodySeparator)
	if len(messageParts) != 2 {
		return "", "", fmt.Errorf("unable to parse sqs message: %s", message.Body)
	}
	rdsType = messageParts[0]
	rdsIdentifier = messageParts[1]
	return rdsType, rdsIdentifier, nil
}

// GetSecurityGroupIds parses the security groups into a slice of strings
func GetSecurityGroupIds(vpcSecurityGroups []types.VpcSecurityGroupMembership) []string {
	var securityGroups []string
	for _, sg := range vpcSecurityGroups {
		securityGroups = append(securityGroups, *sg.VpcSecurityGroupId)
	}
	return securityGroups
}

// IsSqsEvent parses a lambda event to detect if it came from an SQS message
func IsSqsEvent(event []byte) (bool, events.SQSEvent) {
	sqsEvent := events.SQSEvent{}
	_ = json.Unmarshal(event, &sqsEvent)
	if len(sqsEvent.Records) > 0 {
		log.Info("detected SQS event")
		return true, sqsEvent
	}
	return false, sqsEvent
}

// IsEventBridgeRdsEvent parses a lambda event to detect if it came from an EventBridge RDS event
func IsEventBridgeRdsEvent(event []byte) (bool, events.CloudWatchEvent) {
	cloudwatchEvent := events.CloudWatchEvent{}
	_ = json.Unmarshal(event, &cloudwatchEvent)
	if cloudwatchEvent.Source == "aws.rds" {
		log.Info("detected EventBridge RDS event")
		return true, cloudwatchEvent
	}
	return false, cloudwatchEvent
}

// IsScheduledEvent parses a lambda event to detect if it came from a CloudWatch scheduled event
func IsScheduledEvent(event []byte) bool {
	cloudwatchEvent := events.CloudWatchEvent{}
	_ = json.Unmarshal(event, &cloudwatchEvent)
	if cloudwatchEvent.Source == "aws.events" {
		log.Info("detected scheduled event")
		return true
	}
	return false
}

package common

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseSqsMessage_with_valid_message(t *testing.T) {
	rdsType, rdsIdentifier, err := ParseSqsMessage(&events.SQSMessage{Body: "cluster|abc123"})
	assert.Equal(t, "cluster", rdsType)
	assert.Equal(t, "abc123", rdsIdentifier)
	assert.NoError(t, err)
}

func Test_ParseSqsMessage_with_invalid_message(t *testing.T) {
	rdsType, rdsIdentifier, err := ParseSqsMessage(&events.SQSMessage{Body: "foo"})
	assert.Equal(t, "", rdsType)
	assert.Equal(t, "", rdsIdentifier)
	assert.Equal(t, err.Error(), "unable to parse sqs message: foo")
}

func Test_ParseEventBridgeRdsEvent_with_cluster_message(t *testing.T) {
	event := &events.CloudWatchEvent{
		DetailType: "RDS DB Cluster Event",
		Resources:  []string{"arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth"},
	}
	rdsType, rdsIdentifier, err := ParseEventBridgeRdsEvent(event)
	assert.Equal(t, "cluster", rdsType)
	assert.Equal(t, "rds-iam-auth", rdsIdentifier)
	assert.NoError(t, err)
}

func Test_ParseEventBridgeRdsEvent_with_instance_message(t *testing.T) {
	event := &events.CloudWatchEvent{
		DetailType: "RDS DB Instance Event",
		Resources:  []string{"arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth"},
	}
	rdsType, rdsIdentifier, err := ParseEventBridgeRdsEvent(event)
	assert.Equal(t, "instance", rdsType)
	assert.Equal(t, "rds-iam-auth", rdsIdentifier)
	assert.NoError(t, err)
}

func Test_ParseEventBridgeRdsEvent_with_invalid_detail_type(t *testing.T) {
	event := &events.CloudWatchEvent{
		DetailType: "foo",
		Resources:  []string{"arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth"},
	}
	rdsType, rdsIdentifier, err := ParseEventBridgeRdsEvent(event)
	assert.Equal(t, "", rdsType)
	assert.Equal(t, "", rdsIdentifier)
	assert.Equal(t, err.Error(), "unable to parse event detail type: foo")
}

func Test_ParseEventBridgeRdsEvent_with_invalid_resources(t *testing.T) {
	event := &events.CloudWatchEvent{
		DetailType: "RDS DB Instance Event",
		Resources:  []string{"foo"},
	}
	rdsType, rdsIdentifier, err := ParseEventBridgeRdsEvent(event)
	assert.Equal(t, "", rdsType)
	assert.Equal(t, "", rdsIdentifier)
	assert.Equal(t, err.Error(), "unable to parse event resources: [foo]")
}

func Test_GetSecurityGroupIds(t *testing.T) {
	vpcSecurityGroups := []types.VpcSecurityGroupMembership{
		{
			VpcSecurityGroupId: aws.String("abc123"),
		},
	}
	assert.Equal(t, []string{"abc123"}, GetSecurityGroupIds(vpcSecurityGroups))
}

func Test_IsSqsEvent_with_records(t *testing.T) {
	event := `
	{
		"Records": [
			{
				"messageId": "059f36b4-87a3-44ab-83d2-661975830a7d",
				"receiptHandle": "AQEBwJnKyrHigUMZj6rYigCgxlaS3SLy0a...",
				"body": "Test message.",
				"attributes": {
					"ApproximateReceiveCount": "1",
					"SentTimestamp": "1545082649183",
					"SenderId": "AIDAIENQZJOLO23YVJ4VO",
					"ApproximateFirstReceiveTimestamp": "1545082649185"
				},
				"messageAttributes": {},
				"md5OfBody": "e4e68fb7bd0e697a0ae8f1bb342846b3",
				"eventSource": "aws:sqs",
				"eventSourceARN": "arn:aws:sqs:us-east-2:123456789012:my-queue",
				"awsRegion": "us-east-2"
			}
		]
	}
    `
	isSqsEvent, sqsEvent := IsSqsEvent([]byte(event))
	assert.True(t, isSqsEvent)
	assert.Equal(t, 1, len(sqsEvent.Records))
}

func Test_IsSqsEvent_with_no_records(t *testing.T) {
	event := `
	{
		"Records": []
	}
    `
	isSqsEvent, sqsEvent := IsSqsEvent([]byte(event))
	assert.False(t, isSqsEvent)
	assert.Equal(t, 0, len(sqsEvent.Records))
}

func Test_IsSqsEvent_with_invalid_event(t *testing.T) {
	event := `
	{
		"foo": "bar"
	}
    `
	isSqsEvent, sqsEvent := IsSqsEvent([]byte(event))
	assert.False(t, isSqsEvent)
	assert.Equal(t, 0, len(sqsEvent.Records))
}

func Test_IsEventBridgeRdsEvent_with_rds_source(t *testing.T) {
	event := `
	{
	  "version": "0",
	  "id": "9e2d5576-6dea-0ac1-9d7d-5b4ff263397d",
	  "detail-type": "RDS DB Instance Event",
	  "source": "aws.rds",
	  "account": "1234567890",
	  "time": "2023-03-13T21:55:03Z",
	  "region": "us-east-2",
	  "resources": [
		"arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth"
	  ],
	  "detail": {
		"EventCategories": [
		  "availability"
		],
		"SourceType": "DB_INSTANCE",
		"SourceArn": "arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth",
		"Date": "2023-03-13T21:55:03.630Z",
		"Message": "DB instance shutdown",
		"SourceIdentifier": "rds-iam-auth",
		"EventID": "RDS-EVENT-0004"
	  }
	}
    `
	isEventBridgeRdsEvent, cloudwatchEvent := IsEventBridgeRdsEvent([]byte(event))
	assert.True(t, isEventBridgeRdsEvent)
	assert.Equal(t, "aws.rds", cloudwatchEvent.Source)
}

func Test_IsEventBridgeRdsEvent_with_unrecognized_source(t *testing.T) {
	event := `
	{
	  "version": "0",
	  "id": "9e2d5576-6dea-0ac1-9d7d-5b4ff263397d",
	  "detail-type": "RDS DB Instance Event",
	  "source": "foo",
	  "account": "1234567890",
	  "time": "2023-03-13T21:55:03Z",
	  "region": "us-east-2",
	  "resources": [
		"arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth"
	  ],
	  "detail": {
		"EventCategories": [
		  "availability"
		],
		"SourceType": "DB_INSTANCE",
		"SourceArn": "arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth",
		"Date": "2023-03-13T21:55:03.630Z",
		"Message": "DB instance shutdown",
		"SourceIdentifier": "rds-iam-auth",
		"EventID": "RDS-EVENT-0004"
	  }
	}
    `
	isEventBridgeRdsEvent, cloudwatchEvent := IsEventBridgeRdsEvent([]byte(event))
	assert.False(t, isEventBridgeRdsEvent)
	assert.Equal(t, "foo", cloudwatchEvent.Source)
}

func Test_IsScheduledEvent_with_events_source(t *testing.T) {
	event := `
	{
	  "version": "0",
	  "id": "9e2d5576-6dea-0ac1-9d7d-5b4ff263397d",
	  "detail-type": "RDS DB Instance Event",
	  "source": "aws.events",
	  "account": "1234567890",
	  "time": "2023-03-13T21:55:03Z",
	  "region": "us-east-2",
	  "resources": [
		"arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth"
	  ],
	  "detail": {
		"EventCategories": [
		  "availability"
		],
		"SourceType": "DB_INSTANCE",
		"SourceArn": "arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth",
		"Date": "2023-03-13T21:55:03.630Z",
		"Message": "DB instance shutdown",
		"SourceIdentifier": "rds-iam-auth",
		"EventID": "RDS-EVENT-0004"
	  }
	}
    `
	isScheduledEvent := IsScheduledEvent([]byte(event))
	assert.True(t, isScheduledEvent)
}

func Test_IsScheduledEvent_with_unrecognized_source(t *testing.T) {
	event := `
	{
	  "version": "0",
	  "id": "9e2d5576-6dea-0ac1-9d7d-5b4ff263397d",
	  "detail-type": "RDS DB Instance Event",
	  "source": "foo",
	  "account": "1234567890",
	  "time": "2023-03-13T21:55:03Z",
	  "region": "us-east-2",
	  "resources": [
		"arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth"
	  ],
	  "detail": {
		"EventCategories": [
		  "availability"
		],
		"SourceType": "DB_INSTANCE",
		"SourceArn": "arn:aws:rds:us-east-2:1234567890:db:rds-iam-auth",
		"Date": "2023-03-13T21:55:03.630Z",
		"Message": "DB instance shutdown",
		"SourceIdentifier": "rds-iam-auth",
		"EventID": "RDS-EVENT-0004"
	  }
	}
    `
	isScheduledEvent := IsScheduledEvent([]byte(event))
	assert.False(t, isScheduledEvent)
}

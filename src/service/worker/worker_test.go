package worker

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/champ-oss/rds-iam-auth/mocks/mock_rds_client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test_NewService_no_error tests creating a new service
func Test_NewService_no_error(t *testing.T) {
	svc := NewService(nil, nil)
	assert.NotNil(t, svc)
}

// Test_Run__with_instance_no_error tests running with an RDS cluster value
func Test_Run__with_cluster_no_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rdsClient := mock_rds_client.NewMockRdsClientInterface(ctrl)
	svc := Service{nil, rdsClient}
	message := events.SQSMessage{Body: "cluster|cluster1"}
	assert.NoError(t, svc.Run(message))
}

// Test_Run__with_instance_no_error tests running with an RDS instance value
func Test_Run__with_instance_no_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rdsClient := mock_rds_client.NewMockRdsClientInterface(ctrl)
	svc := Service{nil, rdsClient}
	message := events.SQSMessage{Body: "instance|instance1"}
	assert.NoError(t, svc.Run(message))
}

// Test_Run_parsing_error tests passing a message body that cannot be parsed
func Test_Run_parsing_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rdsClient := mock_rds_client.NewMockRdsClientInterface(ctrl)
	svc := Service{nil, rdsClient}

	// test invalid SQS message body "foo"
	message := events.SQSMessage{Body: "foo"}
	assert.ErrorContains(t, svc.Run(message), "unable to parse sqs message: foo")
}

// Test_Run_unrecognized_type_error tests passing an unsupported RDS type
func Test_Run_unrecognized_type_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rdsClient := mock_rds_client.NewMockRdsClientInterface(ctrl)
	svc := Service{nil, rdsClient}

	// test invalid RDS type "foo"
	message := events.SQSMessage{Body: "foo|cluster1"}
	assert.ErrorContains(t, svc.Run(message), "unrecognized RDS type: foo")
}

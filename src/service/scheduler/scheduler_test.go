package scheduler

import (
	"fmt"
	"github.com/champ-oss/rds-iam-auth/mocks/mock_rds_client"
	"github.com/champ-oss/rds-iam-auth/mocks/mock_sqs_client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewService_no_error(t *testing.T) {
	svc := NewService(nil, nil, nil)
	assert.NotNil(t, svc)
}

func Test_Run_no_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sqsClient := mock_sqs_client.NewMockSqsClientInterface(ctrl)
	rdsClient := mock_rds_client.NewMockRdsClientInterface(ctrl)

	rdsClient.EXPECT().GetAllDBClusters().Return([]string{"cluster1", "cluster2"})
	sqsClient.EXPECT().Send("cluster|cluster1").Return(nil)
	sqsClient.EXPECT().Send("cluster|cluster2").Return(nil)

	rdsClient.EXPECT().GetAllDBInstances().Return([]string{"instance1", "instance2"})
	sqsClient.EXPECT().Send("instance|instance1").Return(nil)
	sqsClient.EXPECT().Send("instance|instance2").Return(nil)

	svc := Service{nil, sqsClient, rdsClient}
	assert.NoError(t, svc.Run(nil))
}

func Test_Run_with_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sqsClient := mock_sqs_client.NewMockSqsClientInterface(ctrl)
	rdsClient := mock_rds_client.NewMockRdsClientInterface(ctrl)

	rdsClient.EXPECT().GetAllDBClusters().Return([]string{"cluster1", "cluster2"})
	sqsClient.EXPECT().Send("cluster|cluster1").Return(fmt.Errorf("some error"))

	svc := Service{nil, sqsClient, rdsClient}
	assert.Errorf(t, svc.Run(nil), "some error")
}

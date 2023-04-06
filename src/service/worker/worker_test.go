package worker

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/mocks/mock_mysql_client"
	"github.com/champ-oss/rds-iam-auth/mocks/mock_rds_client"
	"github.com/champ-oss/rds-iam-auth/mocks/mock_ssm_client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

// setUpMockService sets up the service for testing
func setUpMockService(t *testing.T) (*Service, *mock_rds_client.MockRdsClientInterface, *mock_ssm_client.MockSsmClientInterface, *mock_mysql_client.MockMysqlClientInterface) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		SsmSearchPatterns:  []string{"%s-password"},
		DbIamReadUsername:  "readUser",
		DbIamAdminUsername: "adminUser",
	}
	rdsClient := mock_rds_client.NewMockRdsClientInterface(ctrl)
	ssmClient := mock_ssm_client.NewMockSsmClientInterface(ctrl)
	mysqlClient := mock_mysql_client.NewMockMysqlClientInterface(ctrl)
	return &Service{&cfg, rdsClient, ssmClient}, rdsClient, ssmClient, mysqlClient
}

// Test_NewService_no_error tests creating a new service
func Test_NewService_no_error(t *testing.T) {
	svc := NewService(nil, nil, nil)
	assert.NotNil(t, svc)
}

// Test_Run_with_cluster_no_error tests running successfully with an RDS cluster value
func Test_Run_with_cluster_no_error(t *testing.T) {
	svc, rdsClient, ssmClient, mysqlClient := setUpMockService(t)

	rdsClient.EXPECT().GetDBCluster("cluster1").Return(&types.DBCluster{
		Endpoint:       aws.String("endpoint1"),
		Port:           aws.Int32(1111),
		MasterUsername: aws.String("user"),
		DatabaseName:   aws.String("this"),
		VpcSecurityGroups: []types.VpcSecurityGroupMembership{
			{
				VpcSecurityGroupId: aws.String("sg1"),
			},
		},
	}, nil)

	ssmClient.EXPECT().SearchByTag("cluster_identifier", "cluster1").Return([]string{"cluster1-password"}, nil)
	ssmClient.EXPECT().SearchByTag("identifier", "cluster1").Return([]string{}, nil)
	ssmClient.EXPECT().GetValue("cluster1-password").Return("password1", nil)

	mysqlClient.EXPECT().Query("CREATE USER IF NOT EXISTS 'readUser'@'%' IDENTIFIED WITH AWSAuthenticationPlugin as 'RDS'").Return(nil)
	mysqlClient.EXPECT().Query("GRANT SELECT ON *.* TO readUser").Return(nil)
	mysqlClient.EXPECT().Query("CREATE USER IF NOT EXISTS 'adminUser'@'%' IDENTIFIED WITH AWSAuthenticationPlugin as 'RDS'").Return(nil)
	mysqlClient.EXPECT().Query("GRANT ALL PRIVILEGES ON `%`.* TO adminUser").Return(nil)
	mysqlClient.EXPECT().Query("FLUSH PRIVILEGES").Return(nil)
	mysqlClient.EXPECT().Query("SELECT Host, User FROM user").Return(nil)
	mysqlClient.EXPECT().CloseDb()

	message := events.SQSMessage{Body: "cluster|cluster1"}
	assert.NoError(t, svc.Run(&message, mysqlClient))
}

// Test_Run_with_error_finding_cluster tests being unable to find the RDS cluster
func Test_Run_with_error_finding_cluster(t *testing.T) {
	svc, rdsClient, _, _ := setUpMockService(t)
	rdsClient.EXPECT().GetDBCluster("cluster1").Return(nil, fmt.Errorf("unable to find"))

	message := events.SQSMessage{Body: "cluster|cluster1"}
	assert.ErrorContains(t, svc.Run(&message, nil), "unable to find")
}

// Test_Run_with_instance_no_error tests running successfully with an RDS instance value
func Test_Run_with_instance_no_error(t *testing.T) {
	svc, rdsClient, ssmClient, mysqlClient := setUpMockService(t)

	rdsClient.EXPECT().GetDBInstance("instance1").Return(&types.DBInstance{
		Endpoint: &types.Endpoint{
			Address: aws.String("endpoint1"),
			Port:    1111,
		},
		MasterUsername: aws.String("user"),
		DBName:         aws.String("this"),
	}, nil)

	ssmClient.EXPECT().SearchByTag("cluster_identifier", "instance1").Return([]string{}, nil)
	ssmClient.EXPECT().SearchByTag("identifier", "instance1").Return([]string{"instance1-password"}, nil)
	ssmClient.EXPECT().GetValue("instance1-password").Return("password1", nil)

	mysqlClient.EXPECT().Query("CREATE USER IF NOT EXISTS 'readUser'@'%' IDENTIFIED WITH AWSAuthenticationPlugin as 'RDS'").Return(nil)
	mysqlClient.EXPECT().Query("GRANT SELECT ON *.* TO readUser").Return(nil)
	mysqlClient.EXPECT().Query("CREATE USER IF NOT EXISTS 'adminUser'@'%' IDENTIFIED WITH AWSAuthenticationPlugin as 'RDS'").Return(nil)
	mysqlClient.EXPECT().Query("GRANT ALL PRIVILEGES ON `%`.* TO adminUser").Return(nil)
	mysqlClient.EXPECT().Query("FLUSH PRIVILEGES").Return(nil)
	mysqlClient.EXPECT().Query("SELECT Host, User FROM user").Return(nil)
	mysqlClient.EXPECT().CloseDb()

	message := events.SQSMessage{Body: "instance|instance1"}
	assert.NoError(t, svc.Run(&message, mysqlClient))
}

// Test_Run_with_error_finding_instance tests being unable to find the RDS instance
func Test_Run_with_error_finding_instance(t *testing.T) {
	svc, rdsClient, _, _ := setUpMockService(t)
	rdsClient.EXPECT().GetDBInstance("instance1").Return(nil, fmt.Errorf("unable to find"))

	message := events.SQSMessage{Body: "instance|instance1"}
	assert.ErrorContains(t, svc.Run(&message, nil), "unable to find")
}

// Test_Run_parsing_error tests passing a message body that cannot be parsed
func Test_Run_parsing_error(t *testing.T) {
	svc, _, _, _ := setUpMockService(t)

	// test invalid SQS message body "foo"
	message := events.SQSMessage{Body: "foo"}
	assert.ErrorContains(t, svc.Run(&message, nil), "unable to parse sqs message: foo")
}

// Test_Run_unrecognized_type_error tests passing an unsupported RDS type
func Test_Run_unrecognized_type_error(t *testing.T) {
	svc, _, _, _ := setUpMockService(t)

	// test invalid RDS type "foo"
	message := events.SQSMessage{Body: "foo|cluster1"}
	assert.ErrorContains(t, svc.Run(&message, nil), "unrecognized RDS type: foo")
}

// Test_Run_error_finding_password tests being unable to find the password in SSM
func Test_Run_error_finding_password(t *testing.T) {
	svc, rdsClient, ssmClient, _ := setUpMockService(t)

	rdsClient.EXPECT().GetDBCluster("cluster1").Return(&types.DBCluster{
		Endpoint:       aws.String("endpoint1"),
		Port:           aws.Int32(1111),
		MasterUsername: aws.String("user"),
		DatabaseName:   aws.String("this"),
	}, nil)

	ssmClient.EXPECT().SearchByTag("cluster_identifier", "cluster1").Return([]string{}, nil)
	ssmClient.EXPECT().SearchByTag("identifier", "cluster1").Return([]string{}, nil)
	ssmClient.EXPECT().SearchByName("cluster1-password").Return([]string{}, nil)

	message := events.SQSMessage{Body: "cluster|cluster1"}
	assert.ErrorContains(t, svc.Run(&message, nil), "unable to find")
}

// Test_Run_with_error_connecting_mysql tests an error connecting to the mysql server
func Test_Run_with_error_connecting_mysql(t *testing.T) {
	svc, rdsClient, ssmClient, _ := setUpMockService(t)

	rdsClient.EXPECT().GetDBCluster("cluster1").Return(&types.DBCluster{
		Endpoint:       aws.String("localhost"),
		Port:           aws.Int32(65000),
		MasterUsername: aws.String("user"),
		DatabaseName:   aws.String("this"),
		VpcSecurityGroups: []types.VpcSecurityGroupMembership{
			{
				VpcSecurityGroupId: aws.String("sg1"),
			},
		},
	}, nil)

	ssmClient.EXPECT().SearchByTag("cluster_identifier", "cluster1").Return([]string{}, nil)
	ssmClient.EXPECT().SearchByTag("identifier", "cluster1").Return([]string{"cluster1-password"}, nil)
	ssmClient.EXPECT().GetValue("cluster1-password").Return("password1", nil)

	message := events.SQSMessage{Body: "cluster|cluster1"}
	assert.Errorf(t, svc.Run(&message, nil), "some error")
}

// Test_Run_with_error_running_mysql_query tests with an error executing a mysql query
func Test_Run_with_error_running_mysql_query(t *testing.T) {
	svc, rdsClient, ssmClient, mysqlClient := setUpMockService(t)

	rdsClient.EXPECT().GetDBCluster("cluster1").Return(&types.DBCluster{
		Endpoint:       aws.String("endpoint1"),
		Port:           aws.Int32(1111),
		MasterUsername: aws.String("user"),
		DatabaseName:   aws.String("this"),
		VpcSecurityGroups: []types.VpcSecurityGroupMembership{
			{
				VpcSecurityGroupId: aws.String("sg1"),
			},
		},
	}, nil)

	ssmClient.EXPECT().SearchByTag("cluster_identifier", "cluster1").Return([]string{}, nil)
	ssmClient.EXPECT().SearchByTag("identifier", "cluster1").Return([]string{"cluster1-password"}, nil)
	ssmClient.EXPECT().GetValue("cluster1-password").Return("password1", nil)

	mysqlClient.EXPECT().Query("CREATE USER IF NOT EXISTS 'readUser'@'%' IDENTIFIED WITH AWSAuthenticationPlugin as 'RDS'").Return(fmt.Errorf("some error"))
	mysqlClient.EXPECT().CloseDb()

	message := events.SQSMessage{Body: "cluster|cluster1"}
	assert.Errorf(t, svc.Run(&message, mysqlClient), "some error")
}

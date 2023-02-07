package worker

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
	"github.com/champ-oss/rds-iam-auth/pkg/mysql_client"
	"github.com/champ-oss/rds-iam-auth/pkg/rds_client"
	"github.com/champ-oss/rds-iam-auth/pkg/ssm_client"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Service struct {
	config      *cfg.Config
	rdsClient   rds_client.RdsClientInterface
	ssmClient   ssm_client.SsmClientInterface
	mysqlClient mysql_client.MysqlClientInterface
}

// NewService creates a new instance of this service
func NewService(config *cfg.Config, rdsClient rds_client.RdsClientInterface, ssmClient ssm_client.SsmClientInterface, mysqlClient mysql_client.MysqlClientInterface) *Service {
	return &Service{
		config:      config,
		rdsClient:   rdsClient,
		ssmClient:   ssmClient,
		mysqlClient: mysqlClient,
	}
}

// Run is the entrypoint for this service
func (s *Service) Run(message events.SQSMessage) error {
	rdsType, rdsIdentifier, err := parseSqsMessage(message)
	if err != nil {
		return err
	}

	var mySQLConnectionInfo common.MySQLConnectionInfo

	switch rdsType {
	case common.RdsTypeClusterKey:
		mySQLConnectionInfo, err = s.getDBClusterInfo(rdsIdentifier)
		if err != nil {
			return err
		}

	case common.RdsTypeInstanceKey:
		mySQLConnectionInfo, err = s.getDBInstanceInfo(rdsIdentifier)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unrecognized RDS type: %s", rdsType)
	}

	mySQLConnectionInfo.Password, err = s.findPassword(rdsIdentifier)
	if err != nil {
		return err
	}

	_, err = s.mysqlClient.Connect(mySQLConnectionInfo)
	if err != nil {
		return err
	}

	return nil
}

// getDBClusterInfo retrieves connection information for the RDS cluster
func (s *Service) getDBClusterInfo(rdsIdentifier string) (common.MySQLConnectionInfo, error) {
	log.Infof("getting RDS cluster information for: %s", rdsIdentifier)
	cluster, err := s.rdsClient.GetDBCluster(rdsIdentifier)
	if err != nil {
		return common.MySQLConnectionInfo{}, err
	}

	mySQLConnectionInfo := common.MySQLConnectionInfo{
		Endpoint:       *cluster.Endpoint,
		Port:           *cluster.Port,
		Username:       *cluster.MasterUsername,
		Database:       *cluster.DatabaseName,
		SecurityGroups: getSecurityGroupIds(cluster.VpcSecurityGroups),
	}
	log.Debugf("%+v", mySQLConnectionInfo)
	return mySQLConnectionInfo, nil
}

// getDBInstanceInfo retrieves connection information for the RDS instance
func (s *Service) getDBInstanceInfo(rdsIdentifier string) (common.MySQLConnectionInfo, error) {
	log.Infof("getting RDS instance information for: %s", rdsIdentifier)
	instance, err := s.rdsClient.GetDBInstance(rdsIdentifier)
	if err != nil {
		return common.MySQLConnectionInfo{}, err
	}

	mySQLConnectionInfo := common.MySQLConnectionInfo{
		Endpoint:       *instance.Endpoint.Address,
		Port:           instance.Endpoint.Port,
		Username:       *instance.MasterUsername,
		Database:       *instance.DBName,
		SecurityGroups: getSecurityGroupIds(instance.VpcSecurityGroups),
	}
	log.Debugf("%+v", mySQLConnectionInfo)
	return mySQLConnectionInfo, nil
}

// findPassword attempts to locate the RDS password in SSM using the configured search patterns
func (s *Service) findPassword(rdsIdentifier string) (string, error) {
	log.Infof("attempting to find password in SSM for RDS database: %s", rdsIdentifier)
	for _, pattern := range s.config.SsmSearchPatterns {
		// Example of search pattern: "/rds-iam-auth/mysql/%s/password"
		result, _ := s.ssmClient.GetValue(fmt.Sprintf(pattern, rdsIdentifier))
		if result != "" {
			log.Info("password found in ssm")
			return result, nil
		}
	}
	return "", fmt.Errorf("unable to find password in SSM")
}

// parseSqsMessage parses the RDS type and RDS identifier from the incoming SQS message body
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

// getSecurityGroupIds parses the security groups into a slice of strings
func getSecurityGroupIds(vpcSecurityGroups []types.VpcSecurityGroupMembership) []string {
	var securityGroups []string
	for _, sg := range vpcSecurityGroups {
		securityGroups = append(securityGroups, *sg.VpcSecurityGroupId)
	}
	return securityGroups
}

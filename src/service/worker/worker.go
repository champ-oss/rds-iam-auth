package worker

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
	"github.com/champ-oss/rds-iam-auth/pkg/mysql_client"
	"github.com/champ-oss/rds-iam-auth/pkg/rds_client"
	"github.com/champ-oss/rds-iam-auth/pkg/ssm_client"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	config    *cfg.Config
	rdsClient rds_client.RdsClientInterface
	ssmClient ssm_client.SsmClientInterface
}

// NewService creates a new instance of this service
func NewService(config *cfg.Config, rdsClient rds_client.RdsClientInterface, ssmClient ssm_client.SsmClientInterface) *Service {
	return &Service{
		config:    config,
		rdsClient: rdsClient,
		ssmClient: ssmClient,
	}
}

// Run is the entrypoint for this service
func (s *Service) Run(message events.SQSMessage, mysqlClient mysql_client.MysqlClientInterface) error {
	rdsType, rdsIdentifier, err := common.ParseSqsMessage(message)
	if err != nil {
		return err
	}

	mySQLConnectionInfo, err := s.getConnectionInfo(rdsType, rdsIdentifier)
	if err != nil {
		return err
	}

	mySQLConnectionInfo.Password, err = s.findPassword(rdsIdentifier)
	if err != nil {
		return err
	}

	return s.createMysqlIamUsers(mysqlClient, mySQLConnectionInfo)
}

// getConnectionInfo gets connection information and returns common.MySQLConnectionInfo
func (s *Service) getConnectionInfo(rdsType, rdsIdentifier string) (common.MySQLConnectionInfo, error) {
	switch rdsType {
	case common.RdsTypeClusterKey:
		mySQLConnectionInfo, err := s.getDBClusterInfo(rdsIdentifier)
		return mySQLConnectionInfo, err

	case common.RdsTypeInstanceKey:
		mySQLConnectionInfo, err := s.getDBInstanceInfo(rdsIdentifier)
		return mySQLConnectionInfo, err

	default:
		return common.MySQLConnectionInfo{}, fmt.Errorf("unrecognized RDS type: %s", rdsType)
	}
}

// getDBClusterInfo retrieves connection information for the RDS cluster
func (s *Service) getDBClusterInfo(rdsIdentifier string) (common.MySQLConnectionInfo, error) {
	cluster, err := s.rdsClient.GetDBCluster(rdsIdentifier)
	if err != nil {
		return common.MySQLConnectionInfo{}, err
	}

	mySQLConnectionInfo := common.MySQLConnectionInfo{
		Endpoint:       *cluster.Endpoint,
		Port:           *cluster.Port,
		Username:       *cluster.MasterUsername,
		Database:       s.config.DefaultDatabase,
		SecurityGroups: common.GetSecurityGroupIds(cluster.VpcSecurityGroups),
	}
	log.Debugf("%+v", mySQLConnectionInfo)
	return mySQLConnectionInfo, nil
}

// getDBInstanceInfo retrieves connection information for the RDS instance
func (s *Service) getDBInstanceInfo(rdsIdentifier string) (common.MySQLConnectionInfo, error) {
	instance, err := s.rdsClient.GetDBInstance(rdsIdentifier)
	if err != nil {
		return common.MySQLConnectionInfo{}, err
	}

	mySQLConnectionInfo := common.MySQLConnectionInfo{
		Endpoint:       *instance.Endpoint.Address,
		Port:           instance.Endpoint.Port,
		Username:       *instance.MasterUsername,
		Database:       s.config.DefaultDatabase,
		SecurityGroups: common.GetSecurityGroupIds(instance.VpcSecurityGroups),
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

// createMysqlIamUsers executes the SQL queries to set up read-only and admin users for IAM authentication
func (s *Service) createMysqlIamUsers(mysqlClient mysql_client.MysqlClientInterface, mySQLConnectionInfo common.MySQLConnectionInfo) error {
	if mysqlClient == nil {
		var err error
		mysqlClient, err = mysql_client.NewMysqlClient(s.config, mySQLConnectionInfo)
		if err != nil {
			return err
		}
	}
	defer mysqlClient.CloseDb()

	log.Infof("creating read only user: %s", s.config.DbIamReadUsername)
	if err := mysqlClient.Query("CREATE USER IF NOT EXISTS '" + s.config.DbIamReadUsername + "'@'%' IDENTIFIED WITH AWSAuthenticationPlugin as 'RDS'"); err != nil {
		return err
	}

	log.Info("setting read only user permissions")
	if err := mysqlClient.Query("GRANT SELECT ON *.* TO " + s.config.DbIamReadUsername); err != nil {
		return err
	}

	log.Infof("creating admin user: %s", s.config.DbIamAdminUsername)
	if err := mysqlClient.Query("CREATE USER IF NOT EXISTS '" + s.config.DbIamAdminUsername + "'@'%' IDENTIFIED WITH AWSAuthenticationPlugin as 'RDS'"); err != nil {
		return err
	}

	log.Info("setting admin user permissions")
	if err := mysqlClient.Query("GRANT ALL PRIVILEGES ON `%`.* TO " + s.config.DbIamAdminUsername); err != nil {
		return err
	}

	log.Info("flushing privileges")
	if err := mysqlClient.Query("FLUSH PRIVILEGES"); err != nil {
		return err
	}

	log.Info("checking users")
	if err := mysqlClient.Query("SELECT Host, User FROM user"); err != nil {
		return err
	}

	return nil
}

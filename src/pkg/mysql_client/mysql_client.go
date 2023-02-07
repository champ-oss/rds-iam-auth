package mysql_client

import (
	"database/sql"
	"fmt"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"strings"
)

type MysqlClientInterface interface {
	Connect(mySQLConnectionInfo common.MySQLConnectionInfo) (*sql.DB, error)
}

type MysqlClient struct {
}

func NewMysqlClient(config *cfg.Config) *MysqlClient {
	return &MysqlClient{}
}

func (m *MysqlClient) Connect(mySQLConnectionInfo common.MySQLConnectionInfo) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify&allowCleartextPasswords=true",
		mySQLConnectionInfo.Username, mySQLConnectionInfo.Password, mySQLConnectionInfo.Endpoint, mySQLConnectionInfo.Port, mySQLConnectionInfo.Database)

	log.Infof("connecting to MySQL server: %s", strings.ReplaceAll(dsn, mySQLConnectionInfo.Password, "***"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Info("connected successfully")
	return db, err
}

package mysql_client

import (
	"database/sql"
	"fmt"
	cfg "github.com/champ-oss/rds-iam-auth/config"
	"github.com/champ-oss/rds-iam-auth/pkg/common"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type MysqlClientInterface interface {
	CloseDb()
	Query(sql string) error
}

type MysqlClient struct {
	config *cfg.Config
	db     *sql.DB
}

func NewMysqlClient(config *cfg.Config, mySQLConnectionInfo common.MySQLConnectionInfo) (*MysqlClient, error) {
	db, err := connect(mySQLConnectionInfo)
	if err != nil {
		return nil, err
	}

	return &MysqlClient{
		config: config,
		db:     db,
	}, nil
}

// connect creates a connection to the mysql server
func connect(mySQLConnectionInfo common.MySQLConnectionInfo) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify&allowCleartextPasswords=true",
		mySQLConnectionInfo.Username, mySQLConnectionInfo.Password, mySQLConnectionInfo.Endpoint, mySQLConnectionInfo.Port, mySQLConnectionInfo.Database)

	log.Infof("connecting to MySQL server: %s", strings.ReplaceAll(dsn, mySQLConnectionInfo.Password, "***"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Info("connected successfully")
	return db, err
}

// Query executes the given sql query and returns an error
func (m *MysqlClient) Query(sql string) error {
	log.Debug(sql)
	rows, err := m.db.Query(sql)
	if err != nil {
		return err
	}
	defer closeRows(rows)

	for rows.Next() {
		var results []byte
		_ = rows.Scan(&results)
		log.Debugf("query result: %s", results)
	}
	return nil
}

// CloseDb closes the DB connection
func (m *MysqlClient) CloseDb() {
	if err := m.db.Close(); err != nil {
		log.Errorf("unable to close db connection: %s", err)
	}
}

// closeRows closes rows
func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Errorf("unable to close db rows: %s", err)
	}
}

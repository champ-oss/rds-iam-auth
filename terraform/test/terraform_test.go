package test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gruntwork-io/terratest/modules/terraform"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestTerraform(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/complete",
		BackendConfig: map[string]interface{}{
			"bucket": os.Getenv("TF_STATE_BUCKET"),
			"key":    os.Getenv("TF_VAR_git"),
		},
		Vars: map[string]interface{}{},
	}
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	dbUser := "db_iam_read"
	dbName := "this"
	region := terraform.Output(t, terraformOptions, "region")
	testAuroraEndpoint := terraform.Output(t, terraformOptions, "test_aurora_endpoint") + ":3306"
	testMysqlEndpoint := terraform.Output(t, terraformOptions, "test_mysql_endpoint") + ":3306"

	checkDatabaseConnection(testAuroraEndpoint, region, dbUser, dbName)
	checkDatabaseConnection(testMysqlEndpoint, region, dbUser, dbName)
}

// getAWSConfig Logs in to AWS and return a config
func getAWSConfig(region string) aws.Config {
	log.Infof("Getting AWS Config using region: %s", region)
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Loaded AWS configuration successfully")
	return awsConfig
}

func checkDatabaseConnection(dbEndpoint, region, dbUser, dbName string) {
	authenticationToken, err := auth.BuildAuthToken(context.TODO(), dbEndpoint, region, dbUser, getAWSConfig(region).Credentials)
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true", dbUser, authenticationToken, dbEndpoint, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

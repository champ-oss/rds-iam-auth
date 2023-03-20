package test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gruntwork-io/terratest/modules/terraform"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
	defer destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	dbName := "mysql"
	region := terraform.Output(t, terraformOptions, "region")
	functionName := terraform.Output(t, terraformOptions, "function_name")

	testAuroraEndpoint := terraform.Output(t, terraformOptions, "test_aurora_endpoint") + ":3306"
	testMysqlEndpoint := terraform.Output(t, terraformOptions, "test_mysql_endpoint") + ":3306"
	dbIamReadUsername := terraform.Output(t, terraformOptions, "db_iam_read_username")
	dbIamAdminUsername := terraform.Output(t, terraformOptions, "db_iam_admin_username")

	assert.NoError(t, checkDatabaseConnection(testAuroraEndpoint, region, dbIamReadUsername, dbName))
	assert.NoError(t, checkDatabaseConnection(testAuroraEndpoint, region, dbIamAdminUsername, dbName))

	assert.NoError(t, checkDatabaseConnection(testMysqlEndpoint, region, dbIamReadUsername, dbName))
	assert.NoError(t, checkDatabaseConnection(testMysqlEndpoint, region, dbIamAdminUsername, dbName))

	log.Info("invoking the lambda to check for errors")
	output, err := invokeLambda(region, functionName)
	assert.NoError(t, err)
	log.Infof(output)
}

func destroy(t *testing.T, options *terraform.Options) {
	targetedOptions := options
	targetedOptions.Targets = []string{
		"module.aurora",
		"module.mysql",
	}
	terraform.Destroy(t, targetedOptions)
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

// invokeLambda calls an AWS lambda function and waits for the result
func invokeLambda(region, functionName string) (string, error) {
	client := lambda.NewFromConfig(getAWSConfig(region))
	log.Infof("invoking lambda %s", functionName)
	output, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: "RequestResponse",
		LogType:        "Tail",
	})
	log.Info(output.StatusCode)
	return *output.LogResult, err
}

// checkDatabaseConnection logs into a MySQL database using IAM credentials
func checkDatabaseConnection(dbEndpoint, region, dbUser, dbName string) error {
	log.Infof("getting IAM auth token for RDS endpoint: %s", dbEndpoint)
	authenticationToken, err := auth.BuildAuthToken(context.TODO(), dbEndpoint, region, dbUser, getAWSConfig(region).Credentials)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=skip-verify&allowCleartextPasswords=true", dbUser, authenticationToken, dbEndpoint, dbName)

	log.Infof("connecting to MySQL endpoint: %s", dbEndpoint)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}
	log.Info("connected successfully")
	return nil
}

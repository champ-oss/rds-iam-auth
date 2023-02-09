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

	dbName := "mysql"
	region := terraform.Output(t, terraformOptions, "region")
	functionName := terraform.Output(t, terraformOptions, "function_name")
	testAuroraEndpoint := terraform.Output(t, terraformOptions, "test_aurora_endpoint") + ":3306"
	testAuroraMasterUsername := terraform.Output(t, terraformOptions, "test_aurora_master_username")
	testAuroraMasterPassword := terraform.Output(t, terraformOptions, "test_aurora_master_password")
	testMysqlEndpoint := terraform.Output(t, terraformOptions, "test_mysql_endpoint") + ":3306"
	testMysqlMasterUsername := terraform.Output(t, terraformOptions, "test_mysql_master_username")
	testMysqlMasterPassword := terraform.Output(t, terraformOptions, "test_mysql_master_password")
	dbIamReadUsername := terraform.Output(t, terraformOptions, "db_iam_read_username")
	dbIamAdminUsername := terraform.Output(t, terraformOptions, "db_iam_admin_username")

	// Drop the IAM users to reset the test for the next run
	defer dropUsers(testAuroraEndpoint, testAuroraMasterUsername, testAuroraMasterPassword, dbName, []string{dbIamReadUsername, dbIamAdminUsername})
	defer dropUsers(testMysqlEndpoint, testMysqlMasterUsername, testMysqlMasterPassword, dbName, []string{dbIamReadUsername, dbIamAdminUsername})

	invokeLambda(region, functionName)

	checkDatabaseConnection(testAuroraEndpoint, region, dbIamReadUsername, dbName)
	checkDatabaseConnection(testAuroraEndpoint, region, dbIamAdminUsername, dbName)

	checkDatabaseConnection(testMysqlEndpoint, region, dbIamReadUsername, dbName)
	checkDatabaseConnection(testMysqlEndpoint, region, dbIamAdminUsername, dbName)
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
func invokeLambda(region, functionName string) {
	client := lambda.NewFromConfig(getAWSConfig(region))
	log.Infof("invoking lambda %s", functionName)
	output, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: "RequestResponse",
		LogType:        "Tail",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(output.StatusCode)
}

// checkDatabaseConnection logs into a MySQL database using IAM credentials
func checkDatabaseConnection(dbEndpoint, region, dbUser, dbName string) {
	log.Infof("getting IAM auth token for RDS endpoint: %s", dbEndpoint)
	authenticationToken, err := auth.BuildAuthToken(context.TODO(), dbEndpoint, region, dbUser, getAWSConfig(region).Credentials)
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=skip-verify&allowCleartextPasswords=true", dbUser, authenticationToken, dbEndpoint, dbName)

	log.Infof("connecting to MySQL endpoint: %s", dbEndpoint)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("connected successfully")
}

// dropUsers deletes the given usernames from the MySQL server
func dropUsers(dbEndpoint, loginUser, loginPassword, dbName string, dropUsers []string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=skip-verify&allowCleartextPasswords=true", loginUser, loginPassword, dbEndpoint, dbName)

	log.Infof("connecting to MySQL endpoint: %s", dbEndpoint)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("connected successfully")

	for _, user := range dropUsers {
		log.Infof("dropping user: %s", user)
		_, err = db.Query("DROP USER IF EXISTS " + user)
		if err != nil {
			log.Fatal(err)
		}
	}
}

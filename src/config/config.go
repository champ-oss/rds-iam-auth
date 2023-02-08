package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Config struct {
	Debug              bool
	AwsRegion          string
	QueueUrl           string
	AwsConfig          aws.Config
	SsmSearchPatterns  []string
	DbIamReadUsername  string
	DbIamAdminUsername string
	DefaultDatabase    string
}

// LoadConfig loads configuration values from environment variables
func LoadConfig() *Config {
	cfg := Config{
		Debug:     parseBool("DEBUG", true),
		AwsRegion: parseString("AWS_REGION", "us-east-2"),
		QueueUrl:  parseString("QUEUE_URL", ""),
		SsmSearchPatterns: []string{
			"%s-mysql",
			"/rds-iam-auth/mysql/%s/password",
		},
		DbIamReadUsername:  parseString("DB_IAM_READ_USERNAME", "db_iam_read"),
		DbIamAdminUsername: parseString("DB_IAM_ADMIN_USERNAME", "db_iam_admin"),
		DefaultDatabase:    parseString("DEFAULT_DATABASE", "mysql"),
	}

	cfg.AwsConfig = getAWSConfig(cfg.AwsRegion)
	setLogging(cfg.Debug)
	return &cfg
}

func setLogging(debug bool) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debugging mode enabled")
	}
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

// parseBool parses an environment variable as a boolean value
func parseBool(key string, fallback bool) bool {
	if value := os.Getenv(key); strings.ToLower(value) == "true" {
		return true
	}
	if value := os.Getenv(key); strings.ToLower(value) == "false" {
		return false
	}
	return fallback
}

// parseString parses an environment variable as a string value
func parseString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

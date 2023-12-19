set -e

aws rds generate-db-auth-token --hostname $TEST_AURORA_ENDPOINT --port 3306 --username $DB_IAM_READ_USERNAME
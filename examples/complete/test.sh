set -e

# test aurora read only user
AUTH_TOKEN=$(aws rds generate-db-auth-token --hostname $TEST_AURORA_ENDPOINT --port 3306 --username $DB_IAM_READ_USERNAME)
mysql --host=$TEST_AURORA_ENDPOINT --port=3306 --enable-cleartext-plugin --user=$DB_IAM_READ_USERNAME --password=$AUTH_TOKEN

# test aurora admin user
AUTH_TOKEN=$(aws rds generate-db-auth-token --hostname $TEST_AURORA_ENDPOINT --port 3306 --username $DB_IAM_ADMIN_USERNAME)
mysql --host=$TEST_AURORA_ENDPOINT --port=3306 --enable-cleartext-plugin --user=$DB_IAM_ADMIN_USERNAME --password=$AUTH_TOKEN

# test mysql read only user
AUTH_TOKEN=$(aws rds generate-db-auth-token --hostname $TEST_MYSQL_ENDPOINT --port 3306 --username $DB_IAM_READ_USERNAME)
mysql --host=$TEST_MYSQL_ENDPOINT --port=3306 --enable-cleartext-plugin --user=$DB_IAM_READ_USERNAME --password=$AUTH_TOKEN

# test mysql admin user
AUTH_TOKEN=$(aws rds generate-db-auth-token --hostname $TEST_MYSQL_ENDPOINT --port 3306 --username $DB_IAM_ADMIN_USERNAME)
mysql --host=$TEST_MYSQL_ENDPOINT --port=3306 --enable-cleartext-plugin --user=$DB_IAM_ADMIN_USERNAME --password=$AUTH_TOKEN
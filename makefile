run:
	cd src && go run cmd/main.go

test:
	cd src && go test ./...

coverage:
	cd src && go test -json -coverprofile=cover.out ./... > result.json
	cd src && go tool cover -func cover.out
	cd src && go tool cover -html=cover.out

fmt:
	cd src && go fmt ./...
	terraform fmt -recursive -diff

tidy:
	cd src && go mod tidy

download:
	cd src && go mod download

docker:
	docker build -t rds-iam-auth .
	docker run -p 9000:8080 rds-iam-auth

mocks:
	cd src && go install github.com/golang/mock/mockgen@latest
	cd src && mockgen -source pkg/rds_client/rds_client.go -destination mocks/mock_rds_client/mock.go -package mock_rds_client
	cd src && mockgen -source pkg/sqs_client/sqs_client.go -destination mocks/mock_sqs_client/mock.go -package mock_sqs_client
	cd src && mockgen -source pkg/ssm_client/ssm_client.go -destination mocks/mock_ssm_client/mock.go -package mock_ssm_client
	cd src && mockgen -source pkg/mysql_client/mysql_client.go -destination mocks/mock_mysql_client/mock.go -package mock_mysql_client

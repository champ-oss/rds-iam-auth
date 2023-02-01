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
	cd terraform/test && go mod tidy

download:
	cd src && go mod download

docker:
	docker build -t rds-iam-auth .
	docker run -p 9000:8080 rds-iam-auth
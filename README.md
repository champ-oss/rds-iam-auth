# rds-iam-auth

An application which enables AWS IAM authentication on all databases in an AWS VPC.

[![build](https://github.com/champ-oss/rds-iam-auth/actions/workflows/build.yml/badge.svg)](https://github.com/champ-oss/rds-iam-auth/actions/workflows/build.yml)
[![test](https://github.com/champ-oss/rds-iam-auth/actions/workflows/test.yml/badge.svg)](https://github.com/champ-oss/rds-iam-auth/actions/workflows/test.yml)
[![release](https://github.com/champ-oss/rds-iam-auth/actions/workflows/release.yml/badge.svg)](https://github.com/champ-oss/rds-iam-auth/actions/workflows/release.yml)


## Why is this needed?
AWS allows you to enable IAM authentication on RDS clusters and instances, however you must log in to the database itself and run several
SQL commands in order to finish the IAM setup process. This can be a tedious process when you have many databases so this application
automates the setup process for every database it is able to find within a VPC.


## How does it work?
The application first scans for all RDS clusters and instances in the VPC which the application is connected to. Then it uses multiple methods to try
to discover the master password as an SSM parameter. If the password is found, the application connects to the database and runs the SQL commands
necessary to enable IAM authentication.


## Usage
This project is set up to be easily deployed into AWS using Terraform. Please see the terraform/examples/complete folder for an example of everything needed to deploy the application.

AWS services being used:
- AWS Lambda
- AWS SQS


## Troubleshooting
All application logs are exported to CloudWatch Logs. By default, the log group will be: `/aws/lambda/rds-iam-auth-lambda`

## Go Code Structure

#### `src/cmd`
Contains the entry point for the entire application, processes the configuration parameters, and also sets up and injects the dependencies to use.
The application consists of a single Go module. However, the application may run as either the `scheduler` or the `worker` depending on the context in which it is launched.

- scheduler - responsible for gathering a list of all databases and sending the list to AWS SQS
- worker - receives an SQS message and operates on one database at a time

#### `src/config`
Parses environment variables and configures the application at runtime.

#### `src/mocks`
Contains code for testing which is entirely generated using mockgen, which is part of [gomock](https://github.com/golang/mock). Use the command `make mocks` to update the mock code.

#### `src/pkg`
Contains library and utility code for communicating with RDS and SQS using the AWS SDK for Go.

#### `src/service`
Contains the business logic for the `scheduler` and `worker`.



## Terraform Code Structure
#### `terraform/`
Contains the Terraform module code for deploying the application.

#### `terraform/examples/complete`
Contains a working example of deploying to an AWS environment. This example is also used for integration testing the application.

#### `terraform/test`
Contains Go test code for integration testing the application. See the [Integration Testing](#integration-testing) section for more information.


## Integration Testing
This application is fully tested in a live AWS environment, using the [`test` workflow](https://github.com/champ-oss/rds-iam-auth/blob/main/.github/workflows/test.yml).
This workflow runs on every commit and helps ensure that changes do not break the functionality of the application.



## Bugs / Issues / Features
Please use the built-in [GitHub issues tracker](https://github.com/champ-oss/rds-iam-auth/issues).


## Contributing
We welcome any and all contributions! Please see [GitHub issues tracker](https://github.com/champ-oss/rds-iam-auth/issues) for current bugs or enhancements, which may be a good place to start if you would like to contribute.
If you decide to work on an issue, please assign it to yourself. We are happy to review pull requests!


### Setting up a development environment
- Go 1.19 is currently being used for development.
- You can use `make download` to install all the dependencies.
- Use `make run` to test the application locally. This will start a local mysql database as well.
- Use `make test` to run all the unit tests
- Use `make coverage` to run all the unit tests and check code coverage. The coverage report will be opened in your browser.
- Use `make mocks` to generate/update mock test files. This will be needed when updating unit tests.
- Use `make fmt` to properly format the Go and Terraform code.
- Use `make tidy` to run tidy up Go dependencies
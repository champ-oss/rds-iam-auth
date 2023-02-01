terraform {
  backend "s3" {}
}

provider "aws" {
  region = "us-east-2"
}

module "this" {
  source = "../../"
}
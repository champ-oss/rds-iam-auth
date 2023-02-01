provider "aws" {
  region = "us-east-2"
}

module "this" {
  source = "../../"
}
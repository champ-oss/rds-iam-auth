terraform {
  backend "s3" {}
}

provider "aws" {
  region = "us-east-2"
}

locals {
  git = "rds-iam-auth"
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = local.git
  }
}

data "aws_vpcs" "this" {
  tags = {
    purpose = "vega"
  }
}

data "aws_subnets" "private" {
  tags = {
    purpose = "vega"
    Type    = "Private"
  }

  filter {
    name   = "vpc-id"
    values = [data.aws_vpcs.this.ids[0]]
  }
}

resource "aws_security_group" "test" {
  name_prefix = "test-aurora-"
  vpc_id      = data.aws_vpcs.this.ids[0]
}

module "aurora" {
  source                    = "github.com/champ-oss/terraform-aws-aurora.git?ref=v1.0.29-f57eb21"
  cluster_identifier_prefix = local.git
  git                       = local.git
  protect                   = false
  skip_final_snapshot       = true
  vpc_id                    = data.aws_vpcs.this.ids[0]
  private_subnet_ids        = data.aws_subnets.private.ids
  source_security_group_id  = aws_security_group.test.id
  tags                      = local.tags
}

module "mysql" {
  source                   = "git::git@github.com:champ-oss/terraform-aws-mysql.git?ref=v1.0.162-468d0e0"
  vpc_id                   = data.aws_vpcs.this.ids[0]
  private_subnet_ids       = data.aws_subnets.private.ids
  source_security_group_id = aws_security_group.test.id
  name_prefix              = local.git
  git                      = local.git
  skip_final_snapshot      = true
  protect                  = false
  tags                     = local.tags
  name                     = "test"
}

module "this" {
  source = "../../"
}
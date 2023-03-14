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

data "aws_subnets" "public" {
  tags = {
    purpose = "vega"
    Type    = "Public"
  }

  filter {
    name   = "vpc-id"
    values = [data.aws_vpcs.this.ids[0]]
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
  source                    = "github.com/champ-oss/terraform-aws-aurora.git?ref=v1.0.35-db258d6"
  cluster_identifier_prefix = local.git
  git                       = local.git
  protect                   = false
  skip_final_snapshot       = true
  vpc_id                    = data.aws_vpcs.this.ids[0]
  private_subnet_ids        = data.aws_subnets.public.ids
  source_security_group_id  = aws_security_group.test.id
  tags                      = local.tags
  publicly_accessible       = true
  cidr_blocks               = ["0.0.0.0/0"]
}

module "mysql" {
  source                   = "github.com/champ-oss/terraform-aws-mysql.git?ref=v1.0.165-29d9cd6"
  vpc_id                   = data.aws_vpcs.this.ids[0]
  private_subnet_ids       = data.aws_subnets.public.ids
  source_security_group_id = aws_security_group.test.id
  name_prefix              = local.git
  git                      = local.git
  skip_final_snapshot      = true
  protect                  = false
  tags                     = local.tags
  name                     = "test"
  publicly_accessible      = true
  cidr_blocks              = ["0.0.0.0/0"]
}

module "this" {
  source             = "../../"
  vpc_id             = data.aws_vpcs.this.ids[0]
  private_subnet_ids = data.aws_subnets.private.ids
}
locals {
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = var.git
  }
}

module "hash" {
  source   = "github.com/champ-oss/terraform-git-hash.git?ref=v1.0.5-d405e8d"
  path     = "${path.module}/.."
  fallback = ""
}

module "lambda" {
  source = "github.com/champ-oss/terraform-aws-lambda.git?ref=v1.0.111-919a6e1"
  git    = var.git
  name   = "lambda"
  #  vpc_id                          = var.enable_vpc ? var.vpc_id : null
  #  private_subnet_ids              = var.enable_vpc ? var.private_subnet_ids : null
  #  enable_vpc                      = var.enable_vpc
  sync_image       = true
  sync_source_repo = "champtitles/rds-iam-auth"
  ecr_name         = "${var.git}-lambda"
  ecr_tag          = module.hash.hash
  tags             = merge(local.tags, var.tags)
  environment      = {}
}
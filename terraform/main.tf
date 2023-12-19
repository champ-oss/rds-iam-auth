locals {
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = var.git
  }
}

data "aws_region" "this" {}
data "aws_caller_identity" "this" {}

module "hash" {
  source   = "github.com/champ-oss/terraform-git-hash.git?ref=v1.0.14-02da677"
  path     = "${path.module}/.."
  fallback = ""
}

locals {
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = var.git
  }
}

data "aws_region" "this" {}

module "hash" {
  source   = "github.com/champ-oss/terraform-git-hash.git?ref=v1.0.10-32785f1"
  path     = "${path.module}/.."
  fallback = ""
}

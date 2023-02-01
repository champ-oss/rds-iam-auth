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

resource "aws_sqs_queue" "this" {
  name                       = var.git
  visibility_timeout_seconds = 30
  tags                       = merge(local.tags, var.tags)
  sqs_managed_sse_enabled    = true
}

resource "aws_lambda_event_source_mapping" "this" {
  event_source_arn = aws_sqs_queue.this.arn
  function_name    = module.lambda.arn
  enabled          = true
  batch_size       = 1
}

data "aws_iam_policy_document" "this" {
  statement {
    actions = [
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "this" {
  name_prefix = var.git
  policy      = data.aws_iam_policy_document.this.json
}

resource "aws_iam_role_policy_attachment" "invoker" {
  policy_arn = aws_iam_policy.this.arn
  role       = module.lambda.role_name
}
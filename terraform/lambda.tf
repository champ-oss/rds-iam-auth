module "lambda" {
  source              = "github.com/champ-oss/terraform-aws-lambda.git?ref=v1.0.111-919a6e1"
  git                 = var.git
  name                = "lambda"
  enable_vpc          = true
  vpc_id              = var.vpc_id
  private_subnet_ids  = var.private_subnet_ids
  sync_image          = true
  sync_source_repo    = "champtitles/rds-iam-auth"
  ecr_name            = "${var.git}-lambda"
  ecr_tag             = module.hash.hash
  enable_cw_event     = true
  schedule_expression = "rate(1 hour)"
  tags                = merge(local.tags, var.tags)
  environment = {
    QUEUE_URL = aws_sqs_queue.this.url
  }
}

resource "aws_lambda_event_source_mapping" "this" {
  event_source_arn = aws_sqs_queue.this.arn
  function_name    = module.lambda.arn
  enabled          = true
  batch_size       = 1
}
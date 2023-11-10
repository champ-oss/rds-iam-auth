module "lambda" {
  source              = "github.com/champ-oss/terraform-aws-lambda.git?ref=v1.0.133-c385eba"
  git                 = var.git
  name                = "rds-iam-auth"
  description         = "https://github.com/champ-oss/rds-iam-auth"
  enable_vpc          = true
  vpc_id              = var.vpc_id
  private_subnet_ids  = var.private_subnet_ids
  sync_image          = true
  sync_source_repo    = "champtitles/rds-iam-auth"
  ecr_name            = "${var.git}-lambda"
  ecr_tag             = module.hash.hash
  enable_cw_event     = true
  schedule_expression = var.schedule_expression
  tags                = merge(local.tags, var.tags)
  environment = {
    QUEUE_URL             = aws_sqs_queue.this.url
    DB_IAM_READ_USERNAME  = var.db_iam_read_username
    DB_IAM_ADMIN_USERNAME = var.db_iam_admin_username
    SSM_SEARCH_PATTERNS   = join(",", var.ssm_search_patterns)
  }
}

resource "aws_lambda_event_source_mapping" "this" {
  event_source_arn = aws_sqs_queue.this.arn
  function_name    = module.lambda.arn
  enabled          = true
  batch_size       = 1
}

resource "aws_lambda_permission" "this" {
  statement_id_prefix = "${var.git}-rds-iam-auth-events-"
  action              = "lambda:InvokeFunction"
  function_name       = module.lambda.function_name
  principal           = "events.amazonaws.com"
}
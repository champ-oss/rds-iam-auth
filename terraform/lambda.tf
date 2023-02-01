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
  environment = {
    AWS_REGION = data.aws_region.this.name
    QUEUE_URL  = aws_sqs_queue.this.url
  }
}

resource "aws_lambda_event_source_mapping" "this" {
  event_source_arn = aws_sqs_queue.this.arn
  function_name    = module.lambda.arn
  enabled          = true
  batch_size       = 1
}
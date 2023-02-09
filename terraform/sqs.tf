resource "aws_sqs_queue" "this" {
  name                       = "${var.git}-rds-iam-auth"
  visibility_timeout_seconds = var.retry_delay_seconds
  message_retention_seconds  = var.retry_timeout_minutes * 60
  tags                       = merge(local.tags, var.tags)
  sqs_managed_sse_enabled    = true
}
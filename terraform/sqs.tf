resource "aws_sqs_queue" "this" {
  name                       = var.git
  visibility_timeout_seconds = 30
  message_retention_seconds  = 4 * 60 * 60
  tags                       = merge(local.tags, var.tags)
  sqs_managed_sse_enabled    = true
}
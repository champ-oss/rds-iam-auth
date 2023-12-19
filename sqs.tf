resource "aws_sqs_queue" "this" {
  name                       = "${var.git}-rds-iam-auth"
  visibility_timeout_seconds = var.retry_delay_seconds
  tags                       = merge(local.tags, var.tags)
  sqs_managed_sse_enabled    = true
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.deadletter.arn
    maxReceiveCount     = var.max_receive_count
  })
}

resource "aws_sqs_queue" "deadletter" {
  name                    = "${var.git}-rds-iam-auth-deadletter"
  sqs_managed_sse_enabled = true
  tags                    = merge(local.tags, var.tags)
  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns = [
      "arn:aws:sqs:${data.aws_region.this.name}:${data.aws_caller_identity.this.account_id}:${var.git}-rds-iam-auth"
    ]
  })
}
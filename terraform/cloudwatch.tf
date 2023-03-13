resource "aws_cloudwatch_event_rule" "this" {
  name_prefix = "${var.git}-rds-iam-auth-"
  is_enabled  = var.event_triggers_enabled
  event_pattern = jsonencode({
    "source" : [
      "aws.rds"
    ],
    detail-type = [
      "RDS DB Snapshot Event",
      "RDS DB Instance Event",
      "RDS DB Cluster Event",
      "RDS DB Snapshot Event"
    ]
  })
}

resource "aws_cloudwatch_event_target" "this" {
  arn  = module.lambda.arn
  rule = aws_cloudwatch_event_rule.this.id
}
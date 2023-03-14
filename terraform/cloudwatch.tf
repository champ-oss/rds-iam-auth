resource "aws_cloudwatch_event_rule" "this" {
  name_prefix = "${var.git}-rds-iam-auth-events-"
  is_enabled  = var.event_triggers_enabled

  # https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-cloudwatch-events.sample.html
  event_pattern = jsonencode({
    "source" : [
      "aws.rds"
    ],
    detail-type = [
      "RDS DB Snapshot Event",
      "RDS DB Instance Event",
      "RDS DB Cluster Event"
    ],
    "detail" : {
      "EventID" : [
        "RDS-EVENT-0005", # DB instance created.
        "RDS-EVENT-0016", # The master password for the DB instance has been reset.
        "RDS-EVENT-0043", # Restored from snapshot
        "RDS-EVENT-0047", # The DB instance was patched.
        "RDS-EVENT-0170", # DB cluster created.
        "RDS-EVENT-0268", # The engine version upgrade finished.
      ]
    }
  })
}

resource "aws_cloudwatch_event_target" "this" {
  arn  = module.lambda.arn
  rule = aws_cloudwatch_event_rule.this.id
}
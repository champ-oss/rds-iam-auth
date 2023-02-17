data "aws_iam_policy_document" "this" {
  statement {
    actions = [
      "sqs:*Message",
      "sqs:*MessageBatch",
      "sqs:GetQueueAttributes",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "rds:Describe*",
      "tag:Get*",
      "ssm:Get*",
      "ssm:Describe*"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "this" {
  name_prefix = "${var.git}-rds-iam-auth"
  policy      = data.aws_iam_policy_document.this.json
}

resource "aws_iam_role_policy_attachment" "this" {
  policy_arn = aws_iam_policy.this.arn
  role       = module.lambda.role_name
}
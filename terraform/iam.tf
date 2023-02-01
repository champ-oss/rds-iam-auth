data "aws_iam_policy_document" "this" {
  statement {
    actions = [
      "sqs:ReceiveMessage",
      "sqs:GetQueueAttributes",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "rds:Describe*"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "this" {
  name_prefix = var.git
  policy      = data.aws_iam_policy_document.this.json
}

resource "aws_iam_role_policy_attachment" "this" {
  policy_arn = aws_iam_policy.this.arn
  role       = module.lambda.role_name
}
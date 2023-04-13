variable "alarms_email" {
  description = "Email address to send alarms"
  type        = string
  default     = null
}

variable "db_iam_read_username" {
  description = "IAM read only username"
  type        = string
  default     = "db_iam_read"
}

variable "db_iam_admin_username" {
  description = "IAM admin username"
  type        = string
  default     = "db_iam_admin"
}

variable "enable_alarms" {
  description = "Enable CloudWatch metric alarms"
  type        = bool
  default     = true
}

variable "event_triggers_enabled" {
  description = "Trigger based on RDS events"
  type        = bool
  default     = true
}

variable "git" {
  description = "Name of the Git repo"
  type        = string
  default     = "rds-iam-auth"
}

variable "max_receive_count" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue"
  type        = number
  default     = 120
}

variable "private_subnet_ids" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_cluster#subnet_ids"
  type        = list(string)
}

variable "tags" {
  description = "Map of tags to assign to resources"
  type        = map(string)
  default     = {}
}

variable "retry_delay_seconds" {
  description = "How many seconds to wait before retrying the IAM user setup"
  type        = number
  default     = 30
}

variable "retry_timeout_minutes" {
  description = "How many minutes to retry the IAM user setup before giving up"
  type        = number
  default     = 4 * 60
}

variable "schedule_expression" {
  description = "schedule expression using cron"
  type        = string
  default     = "cron(0 4 * * ? *)"
}

variable "ssm_search_patterns" {
  description = "Search strings used to find the SSM parameter containing the database password"
  type        = list(string)
  default = [
    "%s-mysql",
    "/mysql/%s/password",
  ]
}

variable "treat_missing_data" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_metric_alarm#treat_missing_data"
  type        = string
  default     = "ignore"
}

variable "vpc_id" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group#vpc_id"
  type        = string
}


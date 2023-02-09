output "region" {
  value = module.this.region
}

output "test_mysql_endpoint" {
  value = module.mysql.address
}

output "test_aurora_endpoint" {
  value = module.aurora.endpoint
}

output "function_name" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function#function_name"
  value       = module.this.function_name
}

output "db_iam_read_username" {
  value = module.this.db_iam_read_username
}

output "db_iam_admin_username" {
  value = module.this.db_iam_admin_username
}
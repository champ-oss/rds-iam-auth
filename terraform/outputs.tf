output "region" {
  value = data.aws_region.this.name
}

output "function_name" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function#function_name"
  value       = module.lambda.function_name
}
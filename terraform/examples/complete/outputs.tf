output "region" {
  value = module.this.region
}

output "test_mysql_endpoint" {
  value = module.mysql.address
}

output "test_aurora_endpoint" {
  value = module.aurora.endpoint
}
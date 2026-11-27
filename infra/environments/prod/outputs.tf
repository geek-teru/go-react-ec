output "vpc_id" {
  value = module.network.vpc_id
}

output "public_subnet_ids" {
  value = module.network.public_subnet_ids
}

output "private_app_subnet_ids" {
  value = module.network.private_app_subnet_ids
}

output "private_db_subnet_ids" {
  value = module.network.private_db_subnet_ids
}

output "api_ecr_repository_url" {
  value = module.ecs.api_ecr_repository_url
}

output "web_ecr_repository_url" {
  value = module.ecs.web_ecr_repository_url
}

output "ecs_cluster_name" {
  value = module.ecs.cluster_name
}

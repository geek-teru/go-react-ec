# { api = "<account>.dkr.ecr.ap-northeast-1.amazonaws.com/go-react-ec-prod/api", web = ... }
output "api_ecr_repository_url" {
  value = aws_ecr_repository.api.repository_url
}

output "web_ecr_repository_url" {
  value = aws_ecr_repository.web.repository_url
}

output "cluster_name" {
  value = aws_ecs_cluster.app.name
}

output "cluster_arn" {
  value = aws_ecs_cluster.app.arn
}

output "task_execution_role_arn" {
  value = aws_iam_role.ecs_task_execution.arn
}

output "api_task_role_arn" {
  value = aws_iam_role.api_task.arn
}

output "web_task_role_arn" {
  value = aws_iam_role.web_task.arn
}

output "api_log_group_name" {
  value = aws_cloudwatch_log_group.api.name
}

output "web_log_group_name" {
  value = aws_cloudwatch_log_group.web.name
}

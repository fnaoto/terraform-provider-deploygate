output "app_collaborator_id" {
  value = data.deploygate_app_collaborator.current.id
}

output "app_collaborator_users" {
  value = data.deploygate_app_collaborator.current.users
}

output "app_collaborator_teams" {
  value = data.deploygate_app_collaborator.current.teams
}

output "app_collaborator_usage" {
  value = data.deploygate_app_collaborator.current.usage
}

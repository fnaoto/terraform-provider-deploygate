# For data app_collaborator

output "data_app_collaborator_id" {
  value = data.deploygate_app_collaborator.current.id
}

output "data_app_collaborator_users" {
  value = data.deploygate_app_collaborator.current.users
}

output "data_app_collaborator_teams" {
  value = data.deploygate_app_collaborator.current.teams
}

output "data_app_collaborator_usage" {
  value = data.deploygate_app_collaborator.current.usage
}

# For resource app_collaborator

output "resource_app_collaborator_id" {
  value = deploygate_app_collaborator.current.id
}

output "resource_app_collaborator_users" {
  value = deploygate_app_collaborator.current.users
}

output "resource_app_collaborator_teams" {
  value = deploygate_app_collaborator.current.teams
}

output "resource_app_collaborator_usage" {
  value = deploygate_app_collaborator.current.usage
}

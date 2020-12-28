# For data app_collaborator

output "data_app_collaborator_id" {
  value = data.deploygate_app_collaborator.current.id
}

output "data_app_collaborator_users" {
  value = data.deploygate_app_collaborator.current.users
}

# For data deploygate_organization_member

output "data_organization_member_id" {
  value = data.deploygate_organization_member.current.id
}

output "data_organization_member_members" {
  value = data.deploygate_organization_member.current.members
}

# For resource app_collaborator

output "resource_app_collaborator_id" {
  value = deploygate_app_collaborator.current.id
}

output "resource_app_collaborator_users" {
  value = deploygate_app_collaborator.current.users
}

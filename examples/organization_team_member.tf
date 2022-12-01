# Data

data "deploygate_organization_team_member" "current" {
  provider     = deploygate.organization
  organization = var.organization
  team         = var.team
}

# Resource

resource "deploygate_organization_team_member" "current" {
  provider     = deploygate.organization
  organization = var.organization
  team         = var.team

  members {
    name = var.add_member_name
  }

  depends_on = [
    deploygate_organization_member.current
  ]
}

# Output

output "data_organization_team_member_id" {
  value = data.deploygate_organization_team_member.current.id
}

output "data_organization_team_member_members" {
  value = data.deploygate_organization_team_member.current.members
}

output "resource_organization_team_member_id" {
  value = deploygate_organization_team_member.current.id
}

output "resource_organization_team_member_members" {
  value = deploygate_organization_team_member.current.members
}

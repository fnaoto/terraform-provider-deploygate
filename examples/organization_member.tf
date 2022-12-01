# Data

data "deploygate_organization_member" "current" {
  provider     = deploygate.organization
  organization = var.organization
}

# Resource

resource "deploygate_organization_member" "current" {
  provider     = deploygate.organization
  organization = var.organization
  members {
    name = var.add_member_name
  }

  depends_on = [
    deploygate_enterprise_member.current
  ]
}

# Output

output "data_organization_member_id" {
  value = data.deploygate_organization_member.current.id
}

output "data_organization_member_members" {
  value = data.deploygate_organization_member.current.members
}

output "resource_organization_member_id" {
  value = deploygate_organization_member.current.id
}

output "resource_organization_member_members" {
  value = deploygate_organization_member.current.members
}

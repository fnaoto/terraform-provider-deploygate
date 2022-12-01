# Data

data "deploygate_enterprise_member" "current" {
  provider   = deploygate.enterprise
  enterprise = var.enterprise
}

# Resource

resource "deploygate_enterprise_member" "current" {
  provider   = deploygate.enterprise
  enterprise = var.enterprise
  members {
    name = var.add_member_name
  }
}

# Output

output "data_enterprise_member_id" {
  value = data.deploygate_enterprise_member.current.id
}

output "data_enterprise_member_members" {
  value = data.deploygate_enterprise_member.current.members
}

output "resource_enterprise_member_id" {
  value = deploygate_enterprise_member.current.id
}

output "resource_enterprise_member_members" {
  value = deploygate_enterprise_member.current.members
}

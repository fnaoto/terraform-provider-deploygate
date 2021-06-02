# Variables

variable "organization" {
  type = string
}

variable "add_member_name" {
  type = string
}

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

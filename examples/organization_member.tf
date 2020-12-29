# Variables

variable "organization" {
  type = string
}

# Data

data "deploygate_organization_member" "current" {
  organization = var.organization
}

# Resource

resource "deploygate_organization_member" "current" {
  organization = var.organization
}

# Output

output "data_organization_member_id" {
  value = data.deploygate_organization_member.current.id
}

output "data_organization_member_members" {
  value = data.deploygate_organization_member.current.members
}

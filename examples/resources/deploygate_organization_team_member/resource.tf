provider "deploygate" {
  alias   = "organization"
  api_key = var.organization_api_key
}

variable "organization_api_key" {
  type = string
}

variable "organization" {
  type = string
}

variable "team" {
  type = string
}

variable "add_member_name" {
  type = string
}

resource "deploygate_organization_member" "current" {
  provider     = deploygate.organization
  organization = var.organization

  members {
    name = var.add_member_name
  }
}

resource "deploygate_organization_team_member" "current" {
  provider     = deploygate.organization
  organization = var.organization
  team         = var.team

  users {
    name = var.add_member_name
  }

  depends_on = [
    deploygate_organization_member.current
  ]
}

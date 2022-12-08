# Variables

provider "deploygate" {
  alias   = "enterprise"
  api_key = var.enterprise_api_key
}

variable "enterprise_api_key" {
  type = string
}

provider "deploygate" {
  alias   = "organization"
  api_key = var.organization_api_key
}

variable "organization_api_key" {
  type = string
}

variable "organization_in_enterprise" {
  type = string
}

variable "organization_only" {
  type = string
}

variable "add_member_name" {
  type = string
}

# For enterprise organization.

resource "deploygate_enterprise_member" "enterprise" {
  provider   = deploygate.enterprise
  enterprise = var.enterprise

  users {
    name = var.add_member_name
  }
}

resource "deploygate_organization_member" "enterprise" {
  provider     = deploygate.enterprise
  organization = var.organization_in_enterprise

  members {
    name = var.add_member_name
  }

  depends_on = [
    deploygate_enterprise_member.enterprise
  ]
}

# For Organization only

resource "deploygate_organization_member" "organization" {
  provider     = deploygate.organization
  organization = var.organization_only

  members {
    name = var.add_member_name
  }
}

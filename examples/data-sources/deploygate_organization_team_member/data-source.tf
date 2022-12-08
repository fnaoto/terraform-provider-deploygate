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

data "deploygate_organization_team_member" "current" {
  provider     = deploygate.organization
  organization = var.organization
  team         = var.team
}

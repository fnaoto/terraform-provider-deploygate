provider "deploygate" {
  alias   = "enterprise"
  api_key = var.enterprise_api_key
}

variable "enterprise_api_key" {
  type = string
}

variable "enterprise" {
  type = string
}

variable "enterprise_organization" {
  type = string
}

data "deploygate_enterprise_organization_member" "current" {
  provider     = deploygate.enterprise
  enterprise   = var.enterprise
  organization = var.enterprise_organization
}

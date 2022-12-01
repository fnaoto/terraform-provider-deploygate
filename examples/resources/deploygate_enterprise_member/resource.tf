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

variable "add_member_name" {
  type = string
}

resource "deploygate_enterprise_member" "current" {
  provider   = deploygate.enterprise
  enterprise = var.enterprise

  users {
    name = var.add_member_name
  }
}

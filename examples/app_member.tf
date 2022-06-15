# Variables

variable "platform" {
  type = string
}

variable "app_id" {
  type = string
}

variable "owner" {
  type = string
}

variable "add_user_name" {
  type = string
}

# Data

data "deploygate_app_member" "current" {
  provider = deploygate.user
  platform = var.platform
  app_id   = var.app_id
  owner    = var.owner
}

# Resource

resource "deploygate_app_member" "current" {
  provider = deploygate.user
  platform = var.platform
  app_id   = var.app_id
  owner    = var.owner
  users {
    name = var.add_user_name
  }
}

# Output

output "data_app_member_id" {
  value = data.deploygate_app_member.current.id
}

output "data_app_member_users" {
  value = data.deploygate_app_member.current.users
}

output "resource_app_member_id" {
  value = deploygate_app_member.current.id
}

output "resource_app_member_users" {
  value = deploygate_app_member.current.users
}

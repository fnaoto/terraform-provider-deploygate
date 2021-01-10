# terraform-provider-deploygate
Terraform DeployGate provider

## Installation

Make binary of terraform plugin.

```
$ git clone git@github.com:fnaoto/terraform-provider-deploygate.git
$ make
```

## Plugin Usage

Terraform resource for app_collaborator.

```terraform
# Provider

provider "deploygate" {
  alias   = "user"
  api_key = var.user_api_key # or export DG_API_KEY
}

variable "user_api_key" {
  type = string
}

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

data "deploygate_app_collaborator" "current" {
  provider = deploygate.user
  platform = var.platform
  app_id   = var.app_id
  owner    = var.owner
}

# Resource

resource "deploygate_app_collaborator" "current" {
  provider = deploygate.user
  platform = var.platform
  app_id   = var.app_id
  owner    = var.owner
  users {
    name = var.add_user_name
  }
}

# Output

output "data_app_collaborator_id" {
  value = data.deploygate_app_collaborator.current.id
}

output "data_app_collaborator_users" {
  value = data.deploygate_app_collaborator.current.users
}

output "resource_app_collaborator_id" {
  value = deploygate_app_collaborator.current.id
}

output "resource_app_collaborator_users" {
  value = deploygate_app_collaborator.current.users
}
```

Terraform resource for organization_member

```terraform
# Provider

provider "deploygate" {
  alias   = "organization"
  api_key = var.organization_api_key # or export DG_API_KEY
}

variable "organization_api_key" {
  type = string
}

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
```

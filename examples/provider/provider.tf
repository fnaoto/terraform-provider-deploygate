terraform {
  required_providers {
    deploygate = {
      source = "fnaoto/deploygate"
    }
  }
}

# Provider for organization api key

provider "deploygate" {
  alias   = "organization"
  api_key = var.organization_api_key
}

variable "organization_api_key" {
  type = string
}

# Provider for enterprise api key

provider "deploygate" {
  alias   = "enterprise"
  api_key = var.enterprise_api_key
}

variable "enterprise_api_key" {
  type = string
}

# Provider via environment variable which DG_API_KEY
#
## export DG_API_KEY="< api_key >"

provider "deploygate" {}

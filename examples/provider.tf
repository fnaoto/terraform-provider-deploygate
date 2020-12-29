provider "deploygate" {
  alias   = "user"
  api_key = var.user_api_key
}

variable "user_api_key" {
  type = string
}

provider "deploygate" {
  alias   = "organization"
  api_key = var.organization_api_key
}

variable "organization_api_key" {
  type = string
}

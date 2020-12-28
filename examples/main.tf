provider "deploygate" {}

# For data app_collaborator

data "deploygate_app_collaborator" "current" {
  platform = var.platform
  app_id   = var.app_id
  owner    = var.owner
}

# For data organization_member

data "deploygate_organization_member" "current" {
  organization = var.organization
}

# For resource app_collaborator

resource "deploygate_app_collaborator" "current" {
  platform = var.platform
  app_id   = var.app_id
  owner    = var.owner
  users {
    name = var.add_user_name
  }
}

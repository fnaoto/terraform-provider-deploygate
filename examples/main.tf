provider "deploygate" {}

data "deploygate_app_collaborator" "current" {
  platform = var.platform
  app_id   = var.app_id
  owner    = var.owner
}

resource "deploygate_app_collaborator" "current" {
  platform = var.platform
  app_id   = var.app_id
  owner    = var.owner
  users {
    name = var.add_user_name
  }
}

terraform {
  required_providers {
    interfere = {
      source = "interfere-inc/interfere"
    }
  }
}

variable "workspace_slug" {
  type = string
}

variable "surface_slug" {
  type = string
}

provider "interfere" {}

resource "interfere_private_key" "api" {
  workspace_slug           = var.workspace_slug
  name                     = "Automation"
  scopes                   = ["org:surfaces:read", "org:surfaces:write", "org:surfaces:delete"]
  seconds_until_expiration = 2592000
}

resource "interfere_private_key" "release" {
  workspace_slug           = var.workspace_slug
  surface_slug             = var.surface_slug
  name                     = "Releases"
  scopes                   = ["release:write"]
  seconds_until_expiration = null
}

output "api_key" {
  value     = interfere_private_key.api.secret
  sensitive = true
}

output "release_key" {
  value     = interfere_private_key.release.secret
  sensitive = true
}

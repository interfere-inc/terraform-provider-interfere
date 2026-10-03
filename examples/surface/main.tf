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

provider "interfere" {}

resource "interfere_surface" "app" {
  workspace_slug = var.workspace_slug
  name           = "Example application"
  type           = "react"
}

output "surface_slug" {
  value = interfere_surface.app.slug
}

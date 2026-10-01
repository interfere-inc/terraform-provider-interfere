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

variable "creation_id" {
  type        = string
  description = "A nonzero UUID for this creation attempt. Keep it unchanged when retrying; supply a new UUID for a replacement."
}

provider "interfere" {}

resource "interfere_surface" "app" {
  workspace_slug  = var.workspace_slug
  idempotency_key = var.creation_id
  name            = "Example application"
  type            = "react"
}

output "surface_slug" {
  value = interfere_surface.app.slug
}

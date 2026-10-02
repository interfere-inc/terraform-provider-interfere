terraform {
  required_providers {
    interfere = {
      source = "interfere-inc/interfere"
    }
  }
}

provider "interfere" {}

variable "workspace_slug" {
  type = string
}

variable "surface_slug" {
  type = string
}

variable "creation_id" {
  type = string
}

resource "interfere_public_key" "browser" {
  workspace_slug  = var.workspace_slug
  surface_slug    = var.surface_slug
  idempotency_key = var.creation_id
  name            = "Browser"
}

output "public_key" {
  value = interfere_public_key.browser.content
}

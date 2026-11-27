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

resource "interfere_public_key" "browser" {
  workspace_slug = var.workspace_slug
  name           = "Browser"
}

data "interfere_public_keys" "workspace" {
  workspace_slug = var.workspace_slug
}

output "public_key" {
  value = interfere_public_key.browser.content
}

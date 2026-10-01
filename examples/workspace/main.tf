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

resource "interfere_workspace" "team" {
  workspace_slug = var.workspace_slug
  name           = "Example team"
}

output "data_residency_location" {
  value = interfere_workspace.team.data_residency_location
}

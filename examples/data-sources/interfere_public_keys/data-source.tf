data "interfere_public_keys" "workspace" {
  workspace_slug = "example"
}

output "browser_public_key" {
  value = one([for key in data.interfere_public_keys.workspace.public_keys : key.content if key.name == "Browser"])
}

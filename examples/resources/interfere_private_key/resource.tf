resource "interfere_private_key" "automation" {
  workspace_slug           = "example"
  name                     = "Storefront"
  scopes                   = ["release:write", "org:telemetry:write"]
  seconds_until_expiration = 2592000
}

resource "interfere_private_key" "automation" {
  workspace_slug           = "example"
  name                     = "Automation"
  scopes                   = ["org:problems:read"]
  seconds_until_expiration = 2592000
}

resource "interfere_private_key" "automation" {
  workspace_slug           = "example"
  idempotency_key          = "11111111-1111-4111-8111-111111111111"
  name                     = "Automation"
  scopes                   = ["org:surfaces:read"]
  seconds_until_expiration = 2592000
}

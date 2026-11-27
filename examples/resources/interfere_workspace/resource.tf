resource "interfere_workspace" "team" {
  workspace_slug = "example"
  name           = "Example team"

  allow_member_workspace_creation = false
  enabled_auth_methods            = ["enterprise_sso", "passkey", "email_code"]
  impersonation_granted           = false
}

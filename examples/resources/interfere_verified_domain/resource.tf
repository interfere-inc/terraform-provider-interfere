resource "interfere_verified_domain" "company" {
  workspace_slug  = "example"
  name            = "example.com"
  enrollment_mode = "manual_invitation"
}

output "domain_verification_status" {
  value = interfere_verified_domain.company.verification_status
}

resource "interfere_tracking_domain" "telemetry" {
  workspace_slug  = "example"
  name            = "telemetry.example.com"
  idempotency_key = "574e2f70-dc84-4ce0-8db7-f5ecf2acdc82"
}
output "tracking_cname_target" {
  value = interfere_tracking_domain.telemetry.cname_target
}

resource "interfere_tracking_domain" "telemetry" {
  workspace_slug = "example"
  name           = "telemetry.example.com"
}
output "tracking_cname_target" {
  value = interfere_tracking_domain.telemetry.cname_target
}

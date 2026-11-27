data "interfere_tracking_domain" "track" {
  workspace_slug = "example"
  name           = "track.example.com"
}

output "tracking_cname_target" {
  value = data.interfere_tracking_domain.track.cname_target
}

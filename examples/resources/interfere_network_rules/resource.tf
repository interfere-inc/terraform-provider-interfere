resource "interfere_network_rules" "workspace" {
  workspace_slug = "example"

  rules = [
    {
      label = "Office"
      value = "203.0.113.0/24"
    },
    {
      label = "Terraform runner egress"
      value = "198.51.100.7/32"
    },
    {
      label = "Corporate VPN"
      value = "2001:db8:1200::/40"
    },
    {
      label = "Zero Trust egress"
      value = "AS13335"
    },
  ]
}

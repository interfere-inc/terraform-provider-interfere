resource "interfere_surface" "app" {
  workspace_slug          = "example"
  idempotency_key         = "11111111-1111-4111-8111-111111111111"
  name                    = "Application"
  type                    = "react"
  anonymous_user_tracking = false
}

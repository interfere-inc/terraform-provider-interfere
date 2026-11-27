resource "interfere_service" "website" {
  workspace_slug = "example"
  id             = "75000000-0000-4000-8000-000000000001"
  type           = "website"

  website_settings = {
    anonymous_replay_policy = "on_error"
    anonymous_user_tracking = false
    environment_filtering   = true
    plugins = {
      replay = true
      logs   = false
    }
  }
}

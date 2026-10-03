resource "interfere_surface_repository" "web" {
  workspace_slug    = "example"
  surface_slug      = "web"
  integration_id    = "installation-id"
  repository_id     = "123456"
  working_directory = "apps/web"
}

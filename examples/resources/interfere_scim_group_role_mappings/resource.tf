resource "interfere_scim_group_role_mappings" "okta" {
  workspace_slug = "example"
  directory_id   = "scim_dir_2abc"
  mappings = [
    { scim_group_id = "00g1admins", role = "org:admin" },
    { scim_group_id = "00g1engineering", role = "org:member" },
  ]
}

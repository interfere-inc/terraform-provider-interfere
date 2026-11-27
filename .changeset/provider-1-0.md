---
"terraform-provider-interfere": major
---

Interfere models telemetry as services. SDK setup uses a workspace secret (`INTERFERE_API_KEY`) and a stable service name; the framework build plugin bakes a release token into the browser bundle and handles source maps automatically. This release drops the legacy App (surface) resources and key selectors that 0.x deprecated.

#### Added

- `interfere_service` manages an existing service's type and website settings without an App. `data.interfere_service` and `data.interfere_services` expose service identities. Destroying the resource preserves its service, settings, credentials, and telemetry.
- Terraform 1.8+ can move legacy `interfere_surface` state into `interfere_service`. A normal refresh resolves the original App identity to its canonical service; permission failures or missing aliases preserve the existing state.

- `interfere_scim_group_role_mappings` accepts the `org:product`, `org:support`, `org:viewer`, `org:billing` and `org:auditor` roles. The workspace refuses them until plan-gated role assignment ships.

#### Breaking changes

- Removed the `interfere_surface`, `interfere_surface_repository` and `interfere_surface_destination` resources.
- Removed the `interfere_surface` and `interfere_surfaces` data sources.
- Removed `surface_slug` from `interfere_public_key`, `interfere_private_key` and `data.interfere_public_keys`. Public keys are workspace keys, and their import ID is `workspace-slug/key-uuid`.

#### Upgrading from 0.x

Do steps 1 and 2 while you are still on 0.x, then upgrade.

1. **Legacy public keys.** Upgrade SDK builds to a workspace secret with `release:write`, keep the service name stable, and rebuild to inject a release token. Keep old keys until deployed bundles no longer need them. Stop Terraform management of an App-scoped key with a `removed` block and `lifecycle { destroy = false }` while still on 0.x to preserve the credential. Simply removing the resource revokes it. The 1.0 provider reads public keys through the workspace API and cannot manage an App-scoped key; do not carry that resource into the upgrade. New SDK setup does not need an `interfere_public_key` resource.
2. **Surfaces and their mappings.** With Terraform 1.8+, replace each `interfere_surface.website` configuration with `interfere_service.website`, keeping `workspace_slug` and setting `id` to the App's canonical service UUID. Add `moved { from = interfere_surface.website to = interfere_service.website }` as a multiline block and run a normal plan with refresh enabled after upgrading. Review the service type and website settings; App framework values such as `react` become service type `website`. For Apps without a canonical service, stop managing the old resource without destroying it. Stop managing `interfere_surface_repository` and `interfere_surface_destination` the same way. Use `removed` blocks with `lifecycle { destroy = false }` (Terraform 1.7 or later) or `terraform state rm`. Replace old surface data sources with the singular or plural service data source.
3. **Private keys.** Delete `surface_slug` from `interfere_private_key` configuration. Existing keys keep working and are not replaced, and the attribute leaves state on the next refresh. New keys are workspace keys.
4. Set the provider version constraint to `~> 1.0` and run `terraform init -upgrade`. Then run `terraform plan`: it should show no changes for the resources you kept.

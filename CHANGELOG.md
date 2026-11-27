# terraform-provider-interfere

## 1.0.0

### Major Changes

- a5dd672: Interfere models telemetry as services. SDK setup uses a workspace secret (`INTERFERE_API_KEY`) and a stable service name; the framework build plugin bakes a release token into the browser bundle and handles source maps automatically. This release drops the legacy App (surface) resources and key selectors that 0.x deprecated.
  
  #### Added
  
  - `interfere_service` manages an existing service's type and website settings without an App. `data.interfere_service` and `data.interfere_services` expose service identities. Destroying the resource preserves its service, settings, credentials, and telemetry.
  
  - `interfere_scim_group_role_mappings` accepts the `org:product`, `org:support`, `org:viewer`, `org:billing` and `org:auditor` roles. The workspace refuses them until plan-gated role assignment ships.
  
  #### Breaking changes
  
  - Removed the `interfere_surface`, `interfere_surface_repository` and `interfere_surface_destination` resources.
  - Removed the `interfere_surface` and `interfere_surfaces` data sources.
  - Removed `surface_slug` from `interfere_public_key`, `interfere_private_key` and `data.interfere_public_keys`. Public keys are workspace keys, and their import ID is `workspace-slug/key-uuid`.
  
  #### Upgrading from 0.x
  
  Do steps 1 and 2 while you are still on 0.x, then upgrade.
  
  1. **Legacy public keys.** Upgrade SDK builds to a workspace secret with `release:write`, keep the service name stable, and rebuild to inject a release token. Keep old keys until deployed bundles no longer need them. Stop Terraform management of an App-scoped key with a `removed` block and `lifecycle { destroy = false }` while still on 0.x to preserve the credential. Simply removing the resource revokes it. The 1.0 provider reads public keys through the workspace API and cannot manage an App-scoped key; do not carry that resource into the upgrade. New SDK setup does not need an `interfere_public_key` resource.
  2. **Surfaces and their mappings.** On 1.0, with Terraform 1.7+, replace each `interfere_surface.website` configuration with `interfere_service.website`, keeping `workspace_slug` and setting `id` to the UUID of the service that replaced the App (list them with `data.interfere_services`). Add a `removed { from = interfere_surface.website  lifecycle { destroy = false } }` block and an `import { to = interfere_service.website  id = "workspace-slug/service-uuid" }` block, plan and apply, then delete both blocks. The provider cannot move `interfere_surface` state. Review the service type and website settings; App framework values such as `react` become service type `website`. Stop managing `interfere_surface_repository` and `interfere_surface_destination` the same way with `removed` blocks and `lifecycle { destroy = false }`. Replace old surface data sources with the singular or plural service data source.
  3. **Private keys.** Delete `surface_slug` from `interfere_private_key` configuration. Existing keys keep working and are not replaced, and the attribute leaves state on the next refresh. New keys are workspace keys.
  4. Set the provider version constraint to `~> 1.0` and run `terraform init -upgrade`. Then run `terraform plan`: it should show no changes for the resources you kept.

## 0.9.0

### Minor Changes

- Regenerate the provider against the deployed Terraform API contract (efdab7054a048a89f67d865395d68eb89f5a67c901c2795a119a66ea9455cd99 → d11476e1f0e43a0bcf3cdb050b0be5d7ea19a9add3a25598af012e9b11165c10). Review generated schema changes and lifecycle compatibility before publishing.

### Patch Changes

- `interfere_surface_repository` treats a detached repository source (connection removed, repository remembered after a GitHub uninstall) as unlinked, so the plan relinks it instead of failing with "Invalid mapping".

## 0.8.0

### Minor Changes

- Add `anonymous_replay_policy` to the `interfere_surface` resource and the `interfere_surface` data source. Set it to `always`, `on_error` or `never` to choose when sessions without an identified user record session replay; omit it to preserve the API setting. Changes update the surface in place, and imports read the current policy. Regenerated against the deployed Terraform API contract (f27531ff9b2de0bd97a8ba2cc302ca8a26bf201c0e7c079593ee712d7e7b2fc4 → efdab7054a048a89f67d865395d68eb89f5a67c901c2795a119a66ea9455cd99), which adds the `surfaces.setAnonymousReplayPolicy` operation. Indirect `golang.org/x` dependencies are updated.

## 0.7.0

### Minor Changes

- Regenerate the provider against the deployed Terraform API contract (65525892b0d6270545f6df1c9651c55e25022c96bd10f6ca6a6ccd46f57c03bd → f27531ff9b2de0bd97a8ba2cc302ca8a26bf201c0e7c079593ee712d7e7b2fc4). Review generated schema changes and lifecycle compatibility before publishing.

## 0.6.0

### Minor Changes

- Regenerate the provider against the deployed Terraform API contract (797a4cdda55b659b5eb1783a60e14a2b050243b9ea032db4d13ea1437214ec23 → 65525892b0d6270545f6df1c9651c55e25022c96bd10f6ca6a6ccd46f57c03bd). Review generated schema changes and lifecycle compatibility before publishing.
- 3aa315e: Upgrade Fern generation and adopt typed API resource contracts. Confirmed missing resources now leave Terraform state; denied or malformed reads preserve it. Surface creation uses the same 48-character name limit as renaming.

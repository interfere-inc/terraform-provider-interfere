# Interfere Terraform provider

Manage existing workspace settings, surfaces, integration mappings, tracking domains, and private and public API keys with Terraform Plugin Framework and a Fern-generated Go client. Install it from the [Terraform Registry](https://registry.terraform.io/providers/interfere-inc/interfere/latest).

| Resource                        | Behavior                                                                                                                            |
| ------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `interfere_workspace`           | Adopts an existing workspace and manages its name and slug. Destroying removes Terraform management without deleting the workspace. |
| `interfere_surface_repository`  | Links a surface to an installed GitHub integration, repository, and working directory. Destroy unlinks only the mapping.            |
| `interfere_surface_destination` | Links a surface to an installed Vercel, Cloudflare, or CLI destination. Destroy unlinks only the mapping.                           |
| `interfere_tracking_domain`     | Provisions a tracking hostname and exposes DNS/TLS verification status.                                                             |
| `interfere_surface`             | Creates, reads, renames, deletes, and imports surfaces.                                                                             |
| `interfere_private_key`         | Creates and revokes workspace or surface keys with explicit scopes and optional expiry. Configuration changes replace the key.      |
| `interfere_public_key`          | Creates, reads, imports, and revokes surface publishable keys. External rotation refreshes the current value.                       |

## Install

```hcl
terraform {
  required_providers {
    interfere = {
      source  = "interfere-inc/interfere"
      version = "~> 0.3.0"
    }
  }
}

provider "interfere" {}
```

Set `INTERFERE_TOKEN` through your shell or secret manager, then run `terraform init` and `terraform plan`. The credential needs the permissions described below for the resources you manage.

## Generate and test

Building and testing requires Go 1.26 or newer and Terraform. The public repository includes generated Go code, so `make build test check` works without Fern or Docker.

Regeneration additionally requires Docker and the Fern CLI pinned in `fern/fern.config.json`. The CLI prerelease is available from the [public Fern snapshot release](https://github.com/skve/fern/releases/tag/pr-17949-ffdd423). Fern Go generator 1.64.2 runs locally without an enterprise entitlement in the tested configuration.

Supply the Terraform OpenAPI 3.1 publication containing the provider’s resource operations:

```sh
make generate OPENAPI=/absolute/path/to/openapi.json FERN=/absolute/path/to/fern
make build test check
```

Resource schemas and lifecycle code are maintained together in `internal/provider`. `make generate` recreates only the Go client, which is included in public source snapshots. Go dependencies are locked in `go.mod` and `go.sum`. Tests use Terraform with a local HTTP fixture and do not create remote resources. Run `make docs` to regenerate the registry documentation.

Fern consumes the API export unchanged. Terraform schemas declare replacement rules and project the created private-key secret into a sensitive attribute. Surface creation credentials are excluded from state. Surface names are limited to 48 characters so they can also be renamed through the API.

## Try the local binary

Create a separate Terraform CLI configuration file with the absolute path to the built binary's directory:

```hcl
provider_installation {
  dev_overrides {
    "interfere-inc/interfere" = "/absolute/path/to/terraform-provider-interfere/bin"
  }
  direct {}
}
```

Set `TF_CLI_CONFIG_FILE` to that file. Set `INTERFERE_TOKEN` through your shell or secret manager to a workspace API key with `org:surfaces:read`, `org:surfaces:write`, and `org:surfaces:delete`. Session and delegated OAuth access tokens also work. Release-only surface keys cannot manage workspace resources.

Create a workspace key through `POST /v3/workspaces/{workspaceSlug}/api-keys` using an authenticated workspace user with `org:workspace_auth:write` and the permissions being granted:

```json
{
  "idempotencyKey": "11111111-1111-4111-8111-111111111111",
  "name": "Terraform",
  "scopes": ["org:surfaces:read", "org:surfaces:write", "org:surfaces:delete"],
  "secondsUntilExpiration": null
}
```

Use a fresh UUID for each key creation. Set `secondsUntilExpiration` to a positive integer for an expiring key, or `null` for no expiry. Both grants fields are required. The response includes `apiKey.secret`; store it in your secret manager. Revoke the key with `DELETE /v3/workspaces/{workspaceSlug}/api-keys/{apiKeyId}`. Key access is limited to its explicit scopes and the creator's current workspace permissions. Removing the creator's membership removes the key's management access.

## Existing workspace settings

Use `examples/workspace/main.tf` with an existing workspace slug. The first apply adopts that workspace and sets its name. Later changes to `workspace_slug` rename it. The ID stays unchanged, and `data_residency_location` is read-only. Destroying the resource leaves the workspace and its current settings intact. This resource never calls workspace creation or deletion APIs.

The provider credential needs `org:workspace_basics:read` and `org:workspace_basics:write`. Prefer a workspace API key for slug changes; a session token carrying the old slug needs to be refreshed afterward.

```sh
terraform import interfere_workspace.team existing-workspace
```

## Private keys

See `examples/private-key/main.tf` for both workspace API access and surface release keys. Omit `surface_slug` for a workspace key. Set `seconds_until_expiration` to a positive number of seconds, or omit it for no expiry. The `expires_at` output is Unix milliseconds or null.

Workspace keys require `org:workspace_auth:read` and `org:workspace_auth:write` on the provider credential. Surface keys require `org:surfaces:read` and `org:surfaces:write`. The credential must also hold every permission being granted; release and ingest grants require surface write permission.

The computed `secret` attribute is sensitive and stored in Terraform state. Use a protected state backend. Import reads metadata only and leaves `secret` null:

```sh
terraform import interfere_private_key.api existing-workspace/11111111-1111-4111-8111-111111111111
```

For an imported key, set `idempotency_key` to its UUID and match its scopes and lifetime in configuration. Give every replacement a new UUID, including replacement after expiry or external revocation. The provider rejects configuration changes that would reuse the old creation UUID before revoking the existing key. Use a different name as well when requesting `create_before_destroy`, because active key names are unique within their workspace or surface.

Refresh never fetches a secret. External rotation of a key created by Terraform reports an error instead of presenting the old secret as current; import it for metadata-only management or replace it with a new key.

## Public keys

Use `interfere_public_key` for a surface's publishable browser credential. Its `content` output can feed build parameters. Creation requires surface write permission; refresh and import require surface read permission.

Import an existing key without changing its value:

```sh
terraform import interfere_public_key.browser existing-workspace/existing-surface/11111111-1111-4111-8111-111111111111
```

Set `idempotency_key` to the imported UUID and match its name. External rotation updates `content` on refresh. Changing the name, workspace, or surface requires replacement with a fresh UUID. Use a distinct name with `create_before_destroy`. Revocation or deletion of the parent surface removes the resource from state; recreating a revoked key requires a fresh UUID. Destroy revokes the current value and its rotation grace-period values.

## Surfaces

Use `examples/surface/main.tf`, supply `workspace_slug` and a nonzero UUID for `creation_id`, then run `terraform plan`. For local development, the override loads your built provider without `terraform init`. Run `terraform apply` only against a workspace where you intend to create a surface.

Keep the creation UUID unchanged after an uncertain response. Use a new UUID for a new creation or a replacement. Changing `type` or `workspace_slug` requires replacement. Creation also issues default credentials, which the provider discards and never stores in state.

Import an existing surface with:

```sh
terraform import interfere_surface.app example-workspace/example-surface
```

The creation UUID cannot be recovered during import. Supply a new UUID in configuration; the next apply records it without recreating the imported surface.

## Read and deletion behavior

An explicit `deletedAt` value removes a surface from state. Private keys leave state on explicit revocation/expiry or a declared key-not-found response. Permission failures, outages, malformed responses, and unexpected resource identities preserve state and report an error. Workspace and surface queries can return `null` when permissions are insufficient, so an empty result cannot safely mean deletion.

The provider disables automatic HTTP retries. Rename and delete acknowledgements must contain `success: true`. API response bodies are excluded from diagnostics because they can contain credentials.

## Releases

The monorepo remains the development source. The public repository contains only this provider directory and its generated Go code, with a separate Git history. Infrastructure provisions the repository, Engineering access, and branch protection.

Sync a reviewed provider snapshot into the public repository, including `internal/sdk`. Exclude `.tools`, `bin`, `dist`, Terraform state, and local Fern metadata. Review the resulting diff and run `make build test check` before committing it.

The public repository's release workflow builds eight OS/architecture combinations when a `v*` tag is pushed. It signs SHA-256 checksums with the dedicated `GPG_PRIVATE_KEY` and `PASSPHRASE` repository secrets, and creates a draft release. Verify its artifacts before publishing. Register the matching public signing key and provider in HCP Terraform once; subsequent published releases notify the registry through its webhook.

Before publishing a release, deploy the corresponding API changes and run acceptance tests against a disposable workspace.

## Surface settings and integration mappings

Set `anonymous_user_tracking` on `interfere_surface` to manage collection of anonymous-user telemetry. Omit it to preserve the API setting. Changes update the existing surface in place. SDK plugin, tracing, and log-sourcing configuration are not managed by this version.

Use `interfere_surface_repository` for a GitHub, GitHub Enterprise Cloud, or GitHub Enterprise Server repository and optional repository-relative `working_directory`. Use `interfere_surface_destination` for a Vercel project, Cloudflare Worker, or CLI destination. Both require an existing installation and `org:surfaces:read` plus `org:integration_settings:write`. Installation and OAuth consent remain outside Terraform.

Import mappings using `workspace-slug/surface-slug`. Creating over a different existing mapping fails and requires import first. Changing integration, repository, working directory, or project updates the link without replacing the surface. Destroy only removes the link.

## Tracking domains

Use `interfere_tracking_domain` with a lowercase hostname and a creation UUID. The credential needs `org:workspace_domains:read` and `org:workspace_domains:write`. Configure a DNS-only CNAME from `name` to the computed `cname_target` through your DNS provider.

Creation returns while verification is pending so the DNS record can depend on this resource. `status` and `ssl_status` refresh as the service processes verification. Apply completion does not mean DNS or TLS is active. The resource manages the workspace hostname, not DNS records.

Import using `workspace-slug/domain-UUID` and set `idempotency_key` to that domain UUID. Use a fresh UUID when replacing or recreating a deleted domain. Name and workspace changes require replacement. Destroy removes the hostname registration and its managed TLS provisioning.

## Data sources

`data.interfere_workspace` looks up workspace identity, name, and residency with `org:workspace_basics:read`. `data.interfere_surface` reads surface identity, type, and anonymous-user tracking with `org:surfaces:read`. Neither takes ownership of the object.

`data.interfere_integration` looks up an installation using `integration_provider` and requires `org:integration_settings:read`. When multiple installations match, set `integration_id` explicitly. Missing, ambiguous, or inaccessible results produce an error. Only installation identity, provider, and status enter state; installation credentials and provider metadata do not.

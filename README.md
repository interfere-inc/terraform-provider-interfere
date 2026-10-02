# Interfere Terraform provider

Manage existing workspace settings, surfaces, and private API keys with HashiCorp's generated schemas and a Fern-generated Go client. Install it from the [Terraform Registry](https://registry.terraform.io/providers/interfere-inc/interfere/latest).

| Resource                | Behavior                                                                                                                            |
| ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `interfere_workspace`   | Adopts an existing workspace and manages its name and slug. Destroying removes Terraform management without deleting the workspace. |
| `interfere_surface`     | Creates, reads, renames, deletes, and imports surfaces.                                                                             |
| `interfere_private_key` | Creates and revokes workspace or surface keys with explicit scopes and optional expiry. Configuration changes replace the key.      |

## Install

```hcl
terraform {
  required_providers {
    interfere = {
      source  = "interfere-inc/interfere"
      version = "~> 0.1.0"
    }
  }
}

provider "interfere" {}
```

Set `INTERFERE_TOKEN` through your shell or secret manager, then run `terraform init` and `terraform plan`. The credential needs the permissions described below for the resources you manage.

## Generate and test

Building and testing requires Go 1.26 or newer and Terraform. The public repository includes generated Go code, so `make build test check` works without Fern or Docker.

Regeneration additionally requires Docker and the Fern CLI pinned in `fern/fern.config.json`. The CLI prerelease is available from the [public Fern snapshot release](https://github.com/skve/fern/releases/tag/pr-17949-ffdd423). Fern Go generator 1.64.2 runs locally without an enterprise entitlement in the tested configuration.

Supply the Terraform OpenAPI 3.1 publication containing the operations mapped in `generator_config.yml`:

```sh
make generate OPENAPI=/absolute/path/to/openapi.json FERN=/absolute/path/to/fern
make build test check
```

The two HashiCorp generators are pinned in the Makefile. Generated schemas and the Go client are recreated by `make generate` and included in public source snapshots. Go dependencies are locked in `go.mod` and `go.sum`. Tests use Terraform with a local HTTP fixture and do not create remote resources. Run `make docs` to regenerate the registry documentation.

Both generators consume the API export unchanged. `generator_config.yml` uses HashiCorp's native operation mappings. The provider adds replacement rules and projects the created private-key secret into a sensitive Terraform attribute. Surface creation credentials are excluded from state. Surface names are limited to 48 characters so they can also be renamed through the API.

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

Sync a reviewed provider snapshot into the public repository, including `internal/sdk` and `internal/resource_*`. Exclude `.tools`, `bin`, `dist`, Terraform state, and local Fern metadata. Review the resulting diff and run `make build test check` before committing it.

The public repository's release workflow builds eight OS/architecture combinations when a `v*` tag is pushed. It signs SHA-256 checksums with the dedicated `GPG_PRIVATE_KEY` and `PASSPHRASE` repository secrets, and creates a draft release. Verify its artifacts before publishing. Register the matching public signing key and provider in HCP Terraform once; subsequent published releases notify the registry through its webhook.

Before publishing a release, deploy the corresponding API changes and run acceptance tests against a disposable workspace.

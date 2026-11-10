# terraform-provider-interfere

## 0.9.0

### Minor Changes

- Regenerate the provider against the deployed Terraform API contract (efdab7054a048a89f67d865395d68eb89f5a67c901c2795a119a66ea9455cd99 → d11476e1f0e43a0bcf3cdb050b0be5d7ea19a9add3a25598af012e9b11165c10). Review generated schema changes and lifecycle compatibility before publishing.

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

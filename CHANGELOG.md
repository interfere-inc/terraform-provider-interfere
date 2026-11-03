# terraform-provider-interfere

## 0.6.0

### Minor Changes

- Regenerate the provider against the deployed Terraform API contract (797a4cdda55b659b5eb1783a60e14a2b050243b9ea032db4d13ea1437214ec23 → 65525892b0d6270545f6df1c9651c55e25022c96bd10f6ca6a6ccd46f57c03bd). Review generated schema changes and lifecycle compatibility before publishing.
- 3aa315e: Upgrade Fern generation and adopt typed API resource contracts. Confirmed missing resources now leave Terraform state; denied or malformed reads preserve it. Surface creation uses the same 48-character name limit as renaming.

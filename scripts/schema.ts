import { createHash } from "node:crypto";

export function canonicalSchema(input: string): string {
  const document = JSON.parse(input);
  if (
    typeof document.openapi !== "string" ||
    !document.openapi.startsWith("3.1.") ||
    !document.paths ||
    !Object.keys(document.paths).some((path) => path.startsWith("/v3/workspaces/"))
  ) {
    throw new Error("Expected the Terraform OpenAPI 3.1 publication");
  }
  return `${JSON.stringify(document, null, 2)}\n`;
}

export function schemaHash(input: string): string {
  return createHash("sha256").update(canonicalSchema(input)).digest("hex");
}

export function schemaChangeset(previous: string, next: string, version: string): string {
  const bump = version.startsWith("0.") ? "minor" : "major";
  return `---\n"terraform-provider-interfere": ${bump}\n---\n\nRegenerate the provider against the deployed Terraform API contract (${schemaHash(previous)} → ${schemaHash(next)}). Review generated schema changes and lifecycle compatibility before publishing.\n`;
}

export async function downloadSchema(): Promise<string> {
  const response = await fetch("https://api.interfere.com/openapi-terraform.json", {
    signal: AbortSignal.timeout(60_000),
  });
  if (!response.ok) throw new Error(`Terraform schema download failed: HTTP ${response.status}`);
  return canonicalSchema(await response.text());
}

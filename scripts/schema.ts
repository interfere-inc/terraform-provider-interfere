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

const DESCRIPTIVE = new Set(["description", "summary", "title", "example", "examples", "externalDocs", "deprecated"]);

function contractFacts(node: unknown, at: string, facts: Set<string>) {
  if (Array.isArray(node)) {
    node.forEach((item, index) => {
      const named = typeof item === "object" && item !== null && "name" in item && "in" in item;
      const key = named ? `${at}[${String(item.in)}:${String(item.name)}]` : `${at}[${index}]`;
      if (typeof item === "object" && item !== null) {
        facts.add(key);
        contractFacts(item, key, facts);
      } else {
        facts.add(`${at}=${JSON.stringify(item)}`);
      }
    });
    return;
  }
  if (typeof node !== "object" || node === null) {
    facts.add(`${at}=${JSON.stringify(node)}`);
    return;
  }
  for (const [key, value] of Object.entries(node)) {
    if (DESCRIPTIVE.has(key)) continue;
    const path = `${at}/${key.replaceAll("~", "~0").replaceAll("/", "~1")}`;
    facts.add(path);
    contractFacts(value, path, facts);
  }
}

function factsOf(input: string) {
  const document = JSON.parse(canonicalSchema(input));
  const facts = new Set<string>();
  contractFacts(document.paths ?? {}, "/paths", facts);
  contractFacts(document.components ?? {}, "/components", facts);
  return facts;
}

const operation = /^\/paths\/[^/]+\/[^/[=]+/;

const tightensRequest = (fact: string) =>
  /\/requestBody\/(?:.*\/)?required=(?!false$)/.test(fact) ||
  /\/parameters\[[^\]]+\]\/required=true$/.test(fact);

export function removesContract(previous: string, next: string): boolean {
  const before = factsOf(previous);
  const after = factsOf(next);
  for (const fact of before) if (!after.has(fact)) return true;
  for (const fact of after) {
    if (before.has(fact) || !tightensRequest(fact)) continue;
    const existing = fact.match(operation)?.[0];
    if (existing && before.has(existing)) return true;
  }
  return false;
}

export function schemaChangeset(previous: string, next: string, version: string): string {
  const bump = version.startsWith("0.") || !removesContract(previous, next) ? "minor" : "major";
  return `---\n"terraform-provider-interfere": ${bump}\n---\n\nRegenerate the provider against the deployed Terraform API contract (${schemaHash(previous)} → ${schemaHash(next)}). Review generated schema changes and lifecycle compatibility before publishing.\n`;
}

export async function downloadSchema(): Promise<string> {
  const response = await fetch("https://api.interfere.com/openapi-terraform.json", {
    signal: AbortSignal.timeout(60_000),
  });
  if (!response.ok) throw new Error(`Terraform schema download failed: HTTP ${response.status}`);
  return canonicalSchema(await response.text());
}

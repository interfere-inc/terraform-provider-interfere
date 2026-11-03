import { appendFile, readdir, readFile, writeFile } from "node:fs/promises";
import { canonicalSchema, schemaHash } from "./schema";

const destination = process.env.SCHEMA_FILE;
if (!destination) throw new Error("SCHEMA_FILE is required");
const response = await fetch("https://api.interfere.com/openapi-terraform.json", {
  signal: AbortSignal.timeout(60_000),
});
if (!response.ok) throw new Error(`Terraform schema download failed: HTTP ${response.status}`);
const next = canonicalSchema(await response.text());
const previous = await readFile("openapi.json", "utf8");
await writeFile(destination, next);
const changed = schemaHash(previous) !== schemaHash(next);
const pending = (await readdir(".changeset")).some(
  (file) => file.endsWith(".md") && file !== "README.md",
);
if (process.env.GITHUB_OUTPUT) {
  await appendFile(process.env.GITHUB_OUTPUT, `needed=${changed || pending}\n`);
}
console.log(`API changed: ${changed}; pending changesets: ${pending}`);

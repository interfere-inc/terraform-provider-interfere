import { readFile, writeFile } from "node:fs/promises";
import { canonicalSchema, schemaChangeset, schemaHash } from "./schema";

async function run(command: string[]) {
  const child = Bun.spawn(command, { stdout: "inherit", stderr: "inherit" });
  if ((await child.exited) !== 0) throw new Error(`Command failed: ${command.join(" ")}`);
}

const source = process.env.SCHEMA_FILE;
if (!source) throw new Error("SCHEMA_FILE is required");
const previous = await readFile("openapi.json", "utf8");
const next = canonicalSchema(await readFile(source, "utf8"));
if (schemaHash(previous) !== schemaHash(next)) {
  const { version } = JSON.parse(await readFile("package.json", "utf8"));
  await writeFile("openapi.json", next);
  await run(["make", "generate", "FERN=./node_modules/.bin/fern"]);
  await writeFile(".changeset/api-refresh.md", schemaChangeset(previous, next, version));
  await run(["git", "add", "--force", "openapi.json", "internal/sdk"]);
}
await run(["bun", "run", "changeset", "version"]);
await run(["bun", "install", "--lockfile-only", "--ignore-scripts"]);
await run(["make", "build", "test", "check"]);
await run(["git", "add", "openapi.sha256", "internal", "docs"]);
await run(["make", "check-generated"]);

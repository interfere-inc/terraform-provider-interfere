import { readdir, rm, writeFile } from "node:fs/promises";
import { downloadSchema } from "./schema";

async function run(command: string[]) {
  const child = Bun.spawn(command, { stdout: "inherit", stderr: "inherit" });
  if ((await child.exited) !== 0) throw new Error(`Command failed: ${command.join(" ")}`);
}

const generated = (await readdir("internal")).filter(
  (name) => name.startsWith("resource_") || name.startsWith("datasource_"),
);
await Promise.all(
  [...generated.map((name) => `internal/${name}`), "internal/sdk", "docs"].map((path) =>
    rm(path, { recursive: true, force: true }),
  ),
);
await writeFile("openapi.json", await downloadSchema());
await run(["make", "generate", "FERN=./node_modules/.bin/fern"]);
await run(["git", "add", "--all", "--force", "openapi.json", "openapi.sha256", "internal", "docs"]);

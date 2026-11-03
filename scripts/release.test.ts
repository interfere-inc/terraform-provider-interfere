import { expect, test } from "bun:test";
import { mkdtemp, mkdir, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import { schemaChangeset } from "./schema";

const root = resolve(import.meta.dir, "..");
const original = await readFile(join(root, "openapi.json"), "utf8");
const next = JSON.parse(original);
next.info.description = "Updated API contract";
const updated = JSON.stringify(next);

test("release detection handles unchanged, changed, pending and invalid API responses", async () => {
  const directory = await mkdtemp(join(tmpdir(), "provider-detection-"));
  try {
    await mkdir(join(directory, ".changeset"));
    await writeFile(join(directory, "openapi.json"), original);
    const preload = join(directory, "fetch.ts");
    const output = join(directory, "output");
    const schema = join(directory, "candidate.json");
    for (const [body, pending, expected] of [
      [original, false, "needed=false\n"],
      [updated, false, "needed=true\n"],
      [original, true, "needed=true\n"],
      ["<html>unavailable</html>", false, null],
    ] as const) {
      await writeFile(preload, `globalThis.fetch = async () => new Response(${JSON.stringify(body)});`);
      await writeFile(output, "");
      if (pending) await writeFile(join(directory, ".changeset/manual.md"), "pending");
      else await rm(join(directory, ".changeset/manual.md"), { force: true });
      const child = Bun.spawn([process.execPath, "--preload", preload, join(root, "scripts/check-release.ts")], {
        cwd: directory,
        env: { ...process.env, SCHEMA_FILE: schema, GITHUB_OUTPUT: output },
        stdout: "pipe",
        stderr: "pipe",
      });
      expect(await child.exited).toBe(expected === null ? 1 : 0);
      expect(await readFile(output, "utf8")).toBe(expected ?? "");
      expect(await readFile(join(directory, "openapi.json"), "utf8")).toBe(original);
    }
  } finally {
    await rm(directory, { recursive: true, force: true });
  }
});

test("real Changesets CLI versions the private Go provider and consumes the API changeset", async () => {
  const directory = await mkdtemp(join(tmpdir(), "provider-version-"));
  try {
    await mkdir(join(directory, ".changeset"));
    await writeFile(join(directory, "package.json"), JSON.stringify({
      name: "terraform-provider-interfere", version: "0.5.0", private: true,
    }));
    await writeFile(join(directory, ".changeset/config.json"), await readFile(join(root, ".changeset/config.json")));
    await writeFile(join(directory, ".changeset/api-refresh.md"), schemaChangeset(original, updated, "0.5.0"));
    const child = Bun.spawn([process.execPath, join(root, "node_modules/@changesets/cli/bin.js"), "version"], {
      cwd: directory, stdout: "pipe", stderr: "pipe",
    });
    const stderr = await new Response(child.stderr).text();
    expect({ code: await child.exited, stderr }).toEqual({ code: 0, stderr: "" });
    expect(JSON.parse(await readFile(join(directory, "package.json"), "utf8")).version).toBe("0.6.0");
    expect(await readFile(join(directory, "CHANGELOG.md"), "utf8")).toContain("Regenerate the provider");
    expect(await Bun.file(join(directory, ".changeset/api-refresh.md")).exists()).toBe(false);
  } finally {
    await rm(directory, { recursive: true, force: true });
  }
});

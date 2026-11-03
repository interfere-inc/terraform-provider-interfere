import { expect, test } from "bun:test";
import { canonicalSchema, schemaChangeset, schemaHash } from "./schema";

const original = JSON.stringify({
  openapi: "3.1.0",
  paths: { "/v3/workspaces/{workspaceSlug}/surfaces": { post: {} } },
});

test("formatting alone does not trigger a release", () => {
  expect(schemaHash(original)).toBe(schemaHash(canonicalSchema(original)));
});

test("rejects an error response or the general HTTP publication", () => {
  expect(() => canonicalSchema("<html>unavailable</html>")).toThrow();
  expect(() => canonicalSchema('{"openapi":"3.0.0","paths":{}}')).toThrow();
  expect(() => canonicalSchema('{"openapi":"3.1.0","paths":{}}')).toThrow();
});

test("API changes produce a conservative version bump with provenance", () => {
  const next = original.replace('"post":{}', '"post":{"description":"Updated"}');
  expect(schemaHash(next)).not.toBe(schemaHash(original));
  expect(schemaChangeset(original, next, "0.5.0")).toContain('"terraform-provider-interfere": minor');
  expect(schemaChangeset(original, next, "1.0.0")).toContain('"terraform-provider-interfere": major');
  expect(schemaChangeset(original, next, "0.5.0")).toContain(schemaHash(next));
});

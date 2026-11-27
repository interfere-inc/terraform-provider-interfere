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

test("API changes produce a version bump with provenance", () => {
  const next = original.replace('"post":{}', '"post":{"description":"Updated"}');
  expect(schemaHash(next)).not.toBe(schemaHash(original));
  expect(schemaChangeset(original, next, "0.5.0")).toContain('"terraform-provider-interfere": minor');
  expect(schemaChangeset(original, next, "0.5.0")).toContain(schemaHash(next));
});

const contract = (services: object | null, schemas: object = {}, extra: object = {}) =>
  JSON.stringify({
    openapi: "3.1.0",
    info: { title: "Interfere", version: crypto.randomUUID() },
    paths: {
      "/v3/workspaces/{workspaceSlug}/domains": { get: {} },
      ...(services ? { "/v3/workspaces/{workspaceSlug}/services": services } : {}),
      ...extra,
    },
    components: { schemas },
  });

const create = (body: object = {}) => ({
  requestBody: {
    content: {
      "application/json": {
        schema: { type: "object", properties: { name: { type: "string" } }, ...body },
      },
    },
  },
});

const service = {
  type: "object",
  required: ["id"],
  properties: { id: { type: "string" }, type: { type: "string", enum: ["website", "worker"] } },
};

const withProperties = (properties: object) => ({ Service: { ...service, properties } });

const base = contract({ get: {}, post: create() }, { Service: service });

const bump = (next: string, version = "1.0.0") =>
  schemaChangeset(base, next, version).match(/"terraform-provider-interfere": (\w+)/)?.[1];

test.each([
  ["a new path", contract({ get: {}, post: create() }, { Service: service }, { "/v3/workspaces/{workspaceSlug}/keys": { get: {} } })],
  ["a new operation", contract({ get: {}, post: create(), delete: {} }, { Service: service })],
  ["a new operation with required fields", contract({ get: {}, post: create() }, { Service: service }, { "/v3/workspaces/{workspaceSlug}/keys": { post: create({ required: ["name"] }) } })],
  ["a new optional property", contract({ get: {}, post: create() }, withProperties({ ...service.properties, name: { type: "string" } }))],
  ["a new enum value", contract({ get: {}, post: create() }, withProperties({ ...service.properties, type: { type: "string", enum: ["website", "worker", "mobile"] } }))],
  ["a description edit and a new API version", contract({ get: { description: "Lists services", summary: "List" }, post: create() }, { Service: service })],
])("on 1.x an additive refresh (%s) is a minor bump", (_change, next) => {
  expect(bump(next)).toBe("minor");
});

test.each([
  ["a removed path", contract(null, { Service: service })],
  ["a removed operation", contract({ post: create() }, { Service: service })],
  ["a removed schema", contract({ get: {}, post: create() })],
  ["a removed property", contract({ get: {}, post: create() }, withProperties({ id: { type: "string" } }))],
  ["a removed enum value", contract({ get: {}, post: create() }, withProperties({ ...service.properties, type: { type: "string", enum: ["website"] } }))],
  ["a changed type", contract({ get: {}, post: create() }, withProperties({ ...service.properties, id: { type: "integer" } }))],
  ["a newly required field on an existing request", contract({ get: {}, post: create({ required: ["name"] }) }, { Service: service })],
  ["a newly required parameter on an existing operation", contract({ get: { parameters: [{ name: "cursor", in: "query", required: true }] }, post: create() }, { Service: service })],
])("on 1.x a removal (%s) is a major bump", (_change, next) => {
  expect(bump(next)).toBe("major");
});

test("the 0.x line keeps every refresh minor, removals included", () => {
  expect(bump(contract(null, { Service: service }), "0.9.0")).toBe("minor");
});

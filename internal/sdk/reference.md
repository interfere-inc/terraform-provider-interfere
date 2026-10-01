# Reference
## Workspaces
<details><summary><code>client.Workspaces.QueryOrganizationsCurrent(WorkspaceSlug, request) -> *sdk.QueryOrganizationsCurrentResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Read the authenticated workspace's current details and settings. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &sdk.QueryOrganizationsCurrentRequest{
    WorkspaceSlug: "workspaceSlug",
}
client.Workspaces.QueryOrganizationsCurrent(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workspaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**args:** `map[string]any` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.QuerySurfacesGetBySlugIncludeDeleted(WorkspaceSlug, request) -> *sdk.QuerySurfacesGetBySlugIncludeDeletedResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Read a surface by slug even if deleted. Use only when investigating historical references; prefer surfaces.getBySlug for active resources. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &sdk.QuerySurfacesGetBySlugIncludeDeletedRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.QuerySurfacesGetBySlugIncludeDeletedRequestArgs{
        SurfaceSlug: "surfaceSlug",
    },
}
client.Workspaces.QuerySurfacesGetBySlugIncludeDeleted(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workspaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.QuerySurfacesGetBySlugIncludeDeletedRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.MutationOrganizationsUpdateBasics(WorkspaceSlug, request) -> *sdk.MutationOrganizationsUpdateBasicsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Update workspace name or other basic settings described by the input schema. Read organizations.current first. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &sdk.MutationOrganizationsUpdateBasicsRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.MutationOrganizationsUpdateBasicsRequestArgs{},
}
client.Workspaces.MutationOrganizationsUpdateBasics(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workspaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.MutationOrganizationsUpdateBasicsRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.MutationSurfacesDelete(WorkspaceSlug, request) -> *sdk.MutationSurfacesDeleteResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Delete a surface. This is destructive; resolve its exact ID and slug with surfaces.list and verify the requested target. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &sdk.MutationSurfacesDeleteRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.MutationSurfacesDeleteRequestArgs{
        SurfaceSlug: "surfaceSlug",
    },
}
client.Workspaces.MutationSurfacesDelete(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workspaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.MutationSurfacesDeleteRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.MutationSurfacesUpdateName(WorkspaceSlug, request) -> *sdk.MutationSurfacesUpdateNameResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Rename a surface. Resolve its ID with surfaces.list and preserve the rest of its settings. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &sdk.MutationSurfacesUpdateNameRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.MutationSurfacesUpdateNameRequestArgs{
        Name: "name",
        SurfaceSlug: "surfaceSlug",
    },
}
client.Workspaces.MutationSurfacesUpdateName(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workspaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.MutationSurfacesUpdateNameRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.CreateWorkspaceAPIKey(WorkspaceSlug, request) -> *sdk.CreateWorkspaceAPIKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a private API key with explicit scopes and expiration. Supply surfaceSlug for a surface key; omit it for a workspace key. Grants cannot exceed the caller's permissions. Reuse the same idempotency key and request after an uncertain response.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &sdk.CreateWorkspaceAPIKeyRequest{
    WorkspaceSlug: "workspaceSlug",
    IdempotencyKey: "idempotencyKey",
    Name: "name",
    Scopes: []sdk.CreateWorkspaceAPIKeyRequestScopesItem{
        sdk.CreateWorkspaceAPIKeyRequestScopesItemOrgSurfacesRead,
    },
}
client.Workspaces.CreateWorkspaceAPIKey(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workspaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**surfaceSlug:** `*string` 
    
</dd>
</dl>

<dl>
<dd>

**idempotencyKey:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**scopes:** `[]sdk.CreateWorkspaceAPIKeyRequestScopesItem` 
    
</dd>
</dl>

<dl>
<dd>

**secondsUntilExpiration:** `*int` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.ReadPrivateKey(WorkspaceSlug, APIKeyID) -> *sdk.ReadPrivateKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Reads a workspace or surface key's grants and expiration without retrieving its secret. Requires workspace authentication read permission for workspace keys or surface read permission for surface keys.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &sdk.ReadPrivateKeyRequest{
    WorkspaceSlug: "workspaceSlug",
    APIKeyID: "apiKeyId",
}
client.Workspaces.ReadPrivateKey(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workspaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**apiKeyID:** `string` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.RevokeWorkspaceAPIKey(WorkspaceSlug, APIKeyID) -> *sdk.RevokeWorkspaceAPIKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Revokes a workspace or surface API key so it can no longer authorize API requests. Requires workspace authentication write permission for workspace keys or surface write permission for surface keys.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &sdk.RevokeWorkspaceAPIKeyRequest{
    WorkspaceSlug: "workspaceSlug",
    APIKeyID: "apiKeyId",
}
client.Workspaces.RevokeWorkspaceAPIKey(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workspaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**apiKeyID:** `string` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Surfaces
<details><summary><code>client.Surfaces.CreateSurface(WorkspaceSlug, request) -> *sdk.CreateSurfaceResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a surface in the workspace. Reuse the original request identity when retrying the same creation attempt.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &sdk.CreateSurfaceRequest{
    WorkspaceSlug: "workspaceSlug",
    APIKey: &sdk.CreateSurfaceRequestAPIKey{
        Scopes: []sdk.CreateSurfaceRequestAPIKeyScopesItem{
            sdk.CreateSurfaceRequestAPIKeyScopesItemOrgSurfacesRead,
        },
    },
    IdempotencyKey: "idempotencyKey",
    Name: "name",
    Type: sdk.CreateSurfaceRequestTypeElysia,
}
client.Surfaces.CreateSurface(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workspaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**apiKey:** `*sdk.CreateSurfaceRequestAPIKey` 
    
</dd>
</dl>

<dl>
<dd>

**idempotencyKey:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**type_:** `sdk.CreateSurfaceRequestType` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>


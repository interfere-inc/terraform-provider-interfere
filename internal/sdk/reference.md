# Reference
## Workspaces
<details><summary><code>client.Workspaces.Current(WorkspaceSlug, request) -> *sdk.QueryOrganizationsCurrentResponse</code></summary>
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
request := &sdk.CurrentWorkspacesRequest{
    WorkspaceSlug: "workspaceSlug",
}
client.Workspaces.Current(
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

<details><summary><code>client.Workspaces.UpdateBasics(WorkspaceSlug, request) -> *sdk.MutationOrganizationsUpdateBasicsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Update workspace name or other basic settings described by the input schema. Read workspaces.current first. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.UpdateBasicsWorkspacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.UpdateBasicsWorkspacesRequestArgs{},
}
client.Workspaces.UpdateBasics(
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

**args:** `*sdk.UpdateBasicsWorkspacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Surfaces
<details><summary><code>client.Surfaces.GetBySlugIncludeDeleted(WorkspaceSlug, request) -> *sdk.QuerySurfacesGetBySlugIncludeDeletedResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Read a surface by slug even if deleted. Use only when investigating historical references; prefer surfaces.get for active resources. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
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
request := &sdk.GetBySlugIncludeDeletedSurfacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.GetBySlugIncludeDeletedSurfacesRequestArgs{
        SurfaceSlug: "surfaceSlug",
    },
}
client.Surfaces.GetBySlugIncludeDeleted(
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

**args:** `*sdk.GetBySlugIncludeDeletedSurfacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Surfaces.Delete(WorkspaceSlug, request) -> *sdk.MutationSurfacesDeleteResponse</code></summary>
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
request := &sdk.DeleteSurfacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.DeleteSurfacesRequestArgs{
        SurfaceSlug: "surfaceSlug",
    },
}
client.Surfaces.Delete(
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

**args:** `*sdk.DeleteSurfacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Surfaces.UpdateName(WorkspaceSlug, request) -> *sdk.MutationSurfacesUpdateNameResponse</code></summary>
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
request := &sdk.UpdateNameSurfacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.UpdateNameSurfacesRequestArgs{
        Name: "name",
        SurfaceSlug: "surfaceSlug",
    },
}
client.Surfaces.UpdateName(
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

**args:** `*sdk.UpdateNameSurfacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

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

## Keys Private
<details><summary><code>client.Keys.Private.Create(WorkspaceSlug, request) -> *sdk.CreateWorkspaceAPIKeyResponse</code></summary>
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
request := &keys.CreatePrivateRequest{
    WorkspaceSlug: "workspaceSlug",
    IdempotencyKey: "idempotencyKey",
    Name: "name",
    Scopes: []keys.CreatePrivateRequestScopesItem{
        keys.CreatePrivateRequestScopesItemOrgSurfacesRead,
    },
}
client.Keys.Private.Create(
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

**scopes:** `[]keys.CreatePrivateRequestScopesItem` 
    
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

<details><summary><code>client.Keys.Private.Get(WorkspaceSlug, APIKeyID) -> *sdk.ReadPrivateKeyResponse</code></summary>
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
request := &keys.GetPrivateRequest{
    WorkspaceSlug: "workspaceSlug",
    APIKeyID: "apiKeyId",
}
client.Keys.Private.Get(
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

<details><summary><code>client.Keys.Private.Delete(WorkspaceSlug, APIKeyID) -> *sdk.RevokeWorkspaceAPIKeyResponse</code></summary>
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
request := &keys.DeletePrivateRequest{
    WorkspaceSlug: "workspaceSlug",
    APIKeyID: "apiKeyId",
}
client.Keys.Private.Delete(
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

## Keys Public
<details><summary><code>client.Keys.Public.Create(WorkspaceSlug, SurfaceSlug, request) -> *sdk.CreateSurfacePublicKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a publishable key for the surface.
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
request := &keys.CreatePublicRequest{
    WorkspaceSlug: "workspaceSlug",
    SurfaceSlug: "surfaceSlug",
    IdempotencyKey: "idempotencyKey",
    Name: "name",
}
client.Keys.Public.Create(
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

**surfaceSlug:** `string` 
    
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
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Keys.Public.Get(WorkspaceSlug, SurfaceSlug, PublicKeyID) -> *sdk.ReadPublicKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Reads a publishable key's identity, current value and revocation status. Requires surface read permission. Revoked keys and keys belonging to deleted surfaces report revoked: true.
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
request := &keys.GetPublicRequest{
    WorkspaceSlug: "workspaceSlug",
    SurfaceSlug: "surfaceSlug",
    PublicKeyID: "publicKeyId",
}
client.Keys.Public.Get(
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

**surfaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**publicKeyID:** `string` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Keys.Public.Delete(WorkspaceSlug, SurfaceSlug, PublicKeyID) -> *sdk.RevokeSurfacePublicKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Revokes a publishable key and all its rotation grace-period values. Requires surface write permission. Repeating revocation of an existing key succeeds without changing its revocation time.
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
request := &keys.DeletePublicRequest{
    WorkspaceSlug: "workspaceSlug",
    SurfaceSlug: "surfaceSlug",
    PublicKeyID: "publicKeyId",
}
client.Keys.Public.Delete(
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

**surfaceSlug:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**publicKeyID:** `string` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>


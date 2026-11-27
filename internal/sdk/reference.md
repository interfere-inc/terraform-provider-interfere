# Reference
## Domains
<details><summary><code>client.Domains.ByIDIncludeDeleted(WorkspaceSlug, request) -> *sdk.QueryDomainsByIDIncludeDeletedResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Read a tracking domain by ID, including its deletion marker. Requires workspace domain read permission; null can also indicate insufficient access. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Requires org:workspace_domains:read; insufficient permission returns an access-denied error. An authorized read with no matching resource returns WORKSPACE_RESOURCE_NOT_FOUND (404). Only use pagination arguments declared in the schema.
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
request := &sdk.ByIDIncludeDeletedDomainsRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.ByIDIncludeDeletedDomainsRequestArgs{
        DomainID: "domainId",
    },
}
_, err := client.Domains.ByIDIncludeDeleted(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.ByIDIncludeDeletedDomainsRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Domains.List(WorkspaceSlug, request) -> sdk.QueryDomainsForOrgResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List domains associated with the workspace, including their verification and enrollment configuration. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
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
request := &sdk.ListDomainsRequest{
    WorkspaceSlug: "workspaceSlug",
}
_, err := client.Domains.List(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
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

## Integrations
<details><summary><code>client.Integrations.Installations(WorkspaceSlug, request) -> sdk.QueryIntegrationsInstallationsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Legacy alias for connections.installations. List integration installations to discover provider, connection state and installation identifiers. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
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
request := &sdk.InstallationsIntegrationsRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.InstallationsIntegrationsRequestArgs{
        Providers: []string{
            "providers",
        },
    },
}
_, err := client.Integrations.Installations(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.InstallationsIntegrationsRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Integrations.LinkSurfaceToDestination(WorkspaceSlug, request) -> *sdk.MutationIntegrationsLinkSurfaceToDestinationResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Link an app (surface) to an integration destination. Inspect the existing installation and mapping first. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.LinkSurfaceToDestinationIntegrationsRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.LinkSurfaceToDestinationIntegrationsRequestArgs{
        ConnectionID: "connectionId",
        ProjectID: "projectId",
        SurfaceSlug: "surfaceSlug",
    },
}
_, err := client.Integrations.LinkSurfaceToDestination(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.LinkSurfaceToDestinationIntegrationsRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Integrations.LinkSurfaceToRepository(WorkspaceSlug, request) -> *sdk.MutationIntegrationsLinkSurfaceToRepositoryResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Link an app (surface) to a repository. Resolve repository and app IDs from their list queries first. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.LinkSurfaceToRepositoryIntegrationsRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.LinkSurfaceToRepositoryIntegrationsRequestArgs{
        ConnectionID: "connectionId",
        RepositoryID: "repositoryId",
        SurfaceSlug: "surfaceSlug",
    },
}
_, err := client.Integrations.LinkSurfaceToRepository(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.LinkSurfaceToRepositoryIntegrationsRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Integrations.UnlinkSurfaceDestinationMapping(WorkspaceSlug, request) -> *sdk.MutationIntegrationsUnlinkSurfaceDestinationMappingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Remove an app's destination mapping. This changes integration routing. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.UnlinkSurfaceDestinationMappingIntegrationsRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.UnlinkSurfaceDestinationMappingIntegrationsRequestArgs{
        SurfaceSlug: "surfaceSlug",
    },
}
_, err := client.Integrations.UnlinkSurfaceDestinationMapping(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.UnlinkSurfaceDestinationMappingIntegrationsRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Integrations.UnlinkSurfaceMapping(WorkspaceSlug, request) -> *sdk.MutationIntegrationsUnlinkSurfaceMappingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Remove an app's integration mapping. Inspect the selected mapping before unlinking it. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.UnlinkSurfaceMappingIntegrationsRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.UnlinkSurfaceMappingIntegrationsRequestArgs{
        SurfaceSlug: "surfaceSlug",
    },
}
_, err := client.Integrations.UnlinkSurfaceMapping(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.UnlinkSurfaceMappingIntegrationsRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Workspaces
<details><summary><code>client.Workspaces.Current(WorkspaceSlug, request) -> *sdk.QueryWorkspacesCurrentResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Read the authenticated workspace's current details and settings. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Requires org:workspace_basics:read; insufficient permission returns an access-denied error. An authorized read with no matching resource returns WORKSPACE_RESOURCE_NOT_FOUND (404). Only use pagination arguments declared in the schema.
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
_, err := client.Workspaces.Current(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
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

<details><summary><code>client.Workspaces.AddDomain(WorkspaceSlug, request) -> *sdk.MutationWorkspacesAddDomainResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Add a workspace domain. This changes domain configuration; inspect domains.list first. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.AddDomainWorkspacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.AddDomainWorkspacesRequestArgs{
        ID: "id",
        Name: "name",
    },
}
_, err := client.Workspaces.AddDomain(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.AddDomainWorkspacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.AddProxyDomain(WorkspaceSlug, request) -> *sdk.MutationWorkspacesAddProxyDomainResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Add a tracking proxy domain to the workspace. Check domains.list and domains.proxyInUse first. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.AddProxyDomainWorkspacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.AddProxyDomainWorkspacesRequestArgs{
        ID: "id",
        Name: "name",
    },
}
_, err := client.Workspaces.AddProxyDomain(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.AddProxyDomainWorkspacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.RemoveDomain(WorkspaceSlug, request) -> *sdk.MutationWorkspacesRemoveDomainResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Remove a workspace domain identified from domains.list. Review authentication and enrollment implications first. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.RemoveDomainWorkspacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.RemoveDomainWorkspacesRequestArgs{
        ClerkDomainID: "clerkDomainId",
    },
}
_, err := client.Workspaces.RemoveDomain(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.RemoveDomainWorkspacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.RemoveProxyDomain(WorkspaceSlug, request) -> *sdk.MutationWorkspacesRemoveProxyDomainResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Remove a tracking proxy domain. This can affect telemetry sent through that hostname. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.RemoveProxyDomainWorkspacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.RemoveProxyDomainWorkspacesRequestArgs{
        Name: "name",
    },
}
_, err := client.Workspaces.RemoveProxyDomain(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.RemoveProxyDomainWorkspacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.UpdateAuthMethods(WorkspaceSlug, request) -> *sdk.MutationWorkspacesUpdateAuthMethodsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Change allowed workspace authentication methods. Inspect existing configuration to avoid unintended sign-in restrictions. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.UpdateAuthMethodsWorkspacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.UpdateAuthMethodsWorkspacesRequestArgs{
        Methods: []sdk.UpdateAuthMethodsWorkspacesRequestArgsMethodsItem{
            sdk.UpdateAuthMethodsWorkspacesRequestArgsMethodsItemEnterpriseSSO,
        },
    },
}
_, err := client.Workspaces.UpdateAuthMethods(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.UpdateAuthMethodsWorkspacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.UpdateBasics(WorkspaceSlug, request) -> *sdk.MutationWorkspacesUpdateBasicsResponse</code></summary>
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
_, err := client.Workspaces.UpdateBasics(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
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

<details><summary><code>client.Workspaces.UpdateDomainEnrollmentMode(WorkspaceSlug, request) -> *sdk.MutationWorkspacesUpdateDomainEnrollmentModeResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Change how users with a workspace domain can enroll. Inspect the domain and current enrollment policy first. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.UpdateDomainEnrollmentModeWorkspacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.UpdateDomainEnrollmentModeWorkspacesRequestArgs{
        ClerkDomainID: "clerkDomainId",
        EnrollmentMode: sdk.UpdateDomainEnrollmentModeWorkspacesRequestArgsEnrollmentModeManualInvitation,
    },
}
_, err := client.Workspaces.UpdateDomainEnrollmentMode(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.UpdateDomainEnrollmentModeWorkspacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.UpdateMemberWorkspaceCreation(WorkspaceSlug, request) -> *sdk.MutationWorkspacesUpdateMemberWorkspaceCreationResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Change whether workspace members can create additional workspaces. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.UpdateMemberWorkspaceCreationWorkspacesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.UpdateMemberWorkspaceCreationWorkspacesRequestArgs{
        AllowMemberWorkspaceCreation: true,
    },
}
_, err := client.Workspaces.UpdateMemberWorkspaceCreation(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.UpdateMemberWorkspaceCreationWorkspacesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.CreateWorkspacePublicKey(WorkspaceSlug, request) -> *sdk.CreateWorkspacePublicKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a publishable key for the whole workspace. Every app in the workspace can ship it; the gateway resolves the app from the telemetry. Requires workspace authentication write permission. Reuse the same idempotency key and request after an uncertain response.
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
request := &sdk.CreateWorkspacePublicKeyRequest{
    WorkspaceSlug: "workspaceSlug",
    IdempotencyKey: "idempotencyKey",
    Name: "name",
}
_, err := client.Workspaces.CreateWorkspacePublicKey(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**idempotencyKey:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — Display name of the key.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.ReadWorkspacePublicKey(WorkspaceSlug, PublicKeyID) -> *sdk.ReadWorkspacePublicKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Reads a workspace publishable key's identity, current value and revocation status. Requires workspace authentication read or app read permission.
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
request := &sdk.ReadWorkspacePublicKeyRequest{
    WorkspaceSlug: "workspaceSlug",
    PublicKeyID: "publicKeyId",
}
_, err := client.Workspaces.ReadWorkspacePublicKey(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
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

<details><summary><code>client.Workspaces.RevokeWorkspacePublicKey(WorkspaceSlug, PublicKeyID) -> *sdk.RevokeWorkspacePublicKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Revokes a workspace publishable key and all its rotation grace-period values. Every app shipping this key stops authenticating, so confirm with the user before revoking the last live workspace key. Requires workspace authentication write permission. Repeating revocation of an existing key succeeds without changing its revocation time.
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
request := &sdk.RevokeWorkspacePublicKeyRequest{
    WorkspaceSlug: "workspaceSlug",
    PublicKeyID: "publicKeyId",
}
_, err := client.Workspaces.RevokeWorkspacePublicKey(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
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

<details><summary><code>client.Workspaces.ReadNetworkRules(WorkspaceSlug) -> *sdk.ReadNetworkRulesResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the networks allowed to reach the workspace dashboard and API. An empty list means no restriction. Requires workspace authentication read permission.
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
request := &sdk.ReadNetworkRulesRequest{
    WorkspaceSlug: "workspaceSlug",
}
_, err := client.Workspaces.ReadNetworkRules(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workspaces.SaveNetworkRules(WorkspaceSlug, request) -> *sdk.SaveNetworkRulesResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Replaces the networks allowed to reach the workspace dashboard and API, and records an audit event. An empty list removes all restrictions. The saving administrator's connection must remain allowed.
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
request := &sdk.SaveNetworkRulesRequest{
    WorkspaceSlug: "workspaceSlug",
    Rules: []*sdk.SaveNetworkRulesRequestRulesItem{
        &sdk.SaveNetworkRulesRequestRulesItem{
            ID: "id",
            Label: "label",
            Value: "value",
        },
    },
}
_, err := client.Workspaces.SaveNetworkRules(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**rules:** `[]*sdk.SaveNetworkRulesRequestRulesItem` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Services
<details><summary><code>client.Services.Get(WorkspaceSlug, request) -> *sdk.QueryServicesByIDResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Read one OpenTelemetry service in this workspace by its ID. Requires integration settings read permission. Use services.forWorkspace to discover service IDs. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
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
request := &sdk.GetServicesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.GetServicesRequestArgs{
        ServiceID: "serviceId",
    },
}
_, err := client.Services.Get(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.GetServicesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Services.ForWorkspace(WorkspaceSlug, request) -> sdk.QueryServicesForWorkspaceResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List the OpenTelemetry services in this workspace, one row for each service.namespace and service.name pair. Requires integration settings read permission. A null namespace means the service sent none; unknown_service means it sent no service.name. Read last-seen times from the service activity endpoint. Returns up to limit rows (default and maximum 500), ordered by namespace and name. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
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
request := &sdk.ForWorkspaceServicesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.ForWorkspaceServicesRequestArgs{},
}
_, err := client.Services.ForWorkspace(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.ForWorkspaceServicesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Services.SetType(WorkspaceSlug, request) -> *sdk.MutationServicesSetTypeResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Set a service's type to website, backend, other or unclassified. The explicit choice takes precedence over automatic inference. Resolve the service ID from services.forWorkspace first. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.SetTypeServicesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.SetTypeServicesRequestArgs{
        ServiceID: "serviceId",
        Type: sdk.SetTypeServicesRequestArgsTypeUnclassified,
    },
}
_, err := client.Services.SetType(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.SetTypeServicesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Services.SetWebsiteSettings(WorkspaceSlug, request) -> *sdk.MutationServicesSetWebsiteSettingsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Replace a service's website settings, including anonymous replay policy, anonymous-user tracking, environment filtering and browser SDK plugins. Read services.byId first and preserve settings the user did not ask to change. This works without an App and leaves workspace defaults unchanged. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.SetWebsiteSettingsServicesRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.SetWebsiteSettingsServicesRequestArgs{
        ServiceID: "serviceId",
        WebsiteSettings: &sdk.SetWebsiteSettingsServicesRequestArgsWebsiteSettings{
            AnonymousReplayPolicy: sdk.SetWebsiteSettingsServicesRequestArgsWebsiteSettingsAnonymousReplayPolicyAlways,
            AnonymousUserTracking: true,
            EnvironmentFiltering: true,
        },
    },
}
_, err := client.Services.SetWebsiteSettings(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.SetWebsiteSettingsServicesRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Support
<details><summary><code>client.Support.SetAccess(WorkspaceSlug, request) -> *sdk.MutationWorkspacesUpdateSupportAccessResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Change support access to the workspace. Verify the requested access setting before applying it. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &sdk.SetAccessSupportRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &sdk.SetAccessSupportRequestArgs{
        ImpersonationGranted: true,
    },
}
_, err := client.Support.SetAccess(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*sdk.SetAccessSupportRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Keys Public
<details><summary><code>client.Keys.Public.ListForSurface(WorkspaceSlug, request) -> sdk.QueryPublicKeysForSurfaceResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List publishable key records for an app (surface). These are separate from private build API keys. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
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
request := &keys.ListForSurfacePublicRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &keys.ListForSurfacePublicRequestArgs{
        SurfaceSlug: "surfaceSlug",
    },
}
_, err := client.Keys.Public.ListForSurface(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*keys.ListForSurfacePublicRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Keys.Public.ForWorkspace(WorkspaceSlug, request) -> sdk.QueryPublicKeysForWorkspaceResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List workspace-scoped publishable key records (surfaceId null). Every app in the workspace can ship a workspace key; these are separate from app keys and from private build API keys. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
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
request := &keys.ForWorkspacePublicRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: map[string]any{
        "key": "value",
    },
}
_, err := client.Keys.Public.ForWorkspace(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
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

<details><summary><code>client.Keys.Public.Get(WorkspaceSlug, SurfaceSlug, PublicKeyID) -> *sdk.ReadPublicKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Reads a publishable key's identity, current value and revocation status. Requires app read permission. Revoked keys report revoked: true, and so do keys of a deleted app unless the key is still mapped to a service and keeps ingesting.
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
_, err := client.Keys.Public.Get(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**surfaceSlug:** `string` — Surface URL slug, using lowercase letters, numbers, and hyphens.
    
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

Revokes a publishable key and all its rotation grace-period values. Requires app write permission. Repeating revocation of an existing key succeeds without changing its revocation time.
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
_, err := client.Keys.Public.Delete(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**surfaceSlug:** `string` — Surface URL slug, using lowercase letters, numbers, and hyphens.
    
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

## Keys Private
<details><summary><code>client.Keys.Private.Create(WorkspaceSlug, request) -> *sdk.CreateWorkspaceAPIKeyResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a workspace secret API key with explicit scopes and expiration. One key with the release:write scope publishes releases for every app in the workspace. Grants cannot exceed the caller's permissions. Reuse the same idempotency key and request after an uncertain response.
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
        keys.CreatePrivateRequestScopesItemOrgTelemetryWrite,
    },
}
_, err := client.Keys.Private.Create(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**idempotencyKey:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — Display name of the key.
    
</dd>
</dl>

<dl>
<dd>

**scopes:** `[]keys.CreatePrivateRequestScopesItem` — Permission scopes granted to the key.
    
</dd>
</dl>

<dl>
<dd>

**secondsUntilExpiration:** `*int` — Positive lifetime in seconds, or null for no expiry.
    
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

Reads a workspace or app key's grants and expiration without retrieving its secret. Requires workspace authentication read permission for workspace keys or app read permission for app keys.
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
_, err := client.Keys.Private.Get(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
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

Revokes a workspace or app API key so it can no longer authorize API requests. Requires workspace authentication write permission for workspace keys or app write permission for app keys.
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
_, err := client.Keys.Private.Delete(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
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

## SCIM Groups
<details><summary><code>client.SCIM.Groups.List(WorkspaceSlug, request) -> sdk.QuerySCIMGroupsForOrgResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List provisioned SCIM groups and their workspace role mappings. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation is read-only despite using POST. Permission-filtered queries may return an empty result; do not assume that proves the resource does not exist. Only use pagination arguments declared in the schema.
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
request := &scim.ListGroupsRequest{
    WorkspaceSlug: "workspaceSlug",
}
_, err := client.SCIM.Groups.List(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
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

<details><summary><code>client.SCIM.Groups.Reorder(WorkspaceSlug, request) -> *sdk.MutationWorkspacesReorderSCIMGroupRolesResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Reorder a SCIM directory's group role mappings; earlier mappings take precedence when a member is in several groups. Send every mapped group of the directory. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &scim.ReorderGroupsRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &scim.ReorderGroupsRequestArgs{
        ClerkSCIMDirectoryID: "clerkScimDirectoryId",
        IDs: []string{
            "ids",
        },
    },
}
_, err := client.SCIM.Groups.Reorder(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*scim.ReorderGroupsRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.SCIM.Groups.SetRole(WorkspaceSlug, request) -> *sdk.MutationWorkspacesSetSCIMGroupRoleResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Change the role assigned to a SCIM group. This can change access for multiple provisioned members. Pass the workspace slug in the URL and operation arguments inside the JSON body as { args: ... }. Inspect the request schema for required fields and exact identifier formats. This operation writes data under the caller's existing permissions and records its audit event. Do not automatically retry after a timeout; read the affected resource to determine whether the write completed.
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
request := &scim.SetRoleGroupsRequest{
    WorkspaceSlug: "workspaceSlug",
    Args: &scim.SetRoleGroupsRequestArgs{
        ID: "id",
    },
}
_, err := client.SCIM.Groups.SetRole(
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

**workspaceSlug:** `string` — Workspace URL slug, using lowercase letters, numbers, and hyphens.
    
</dd>
</dl>

<dl>
<dd>

**args:** `*scim.SetRoleGroupsRequestArgs` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>


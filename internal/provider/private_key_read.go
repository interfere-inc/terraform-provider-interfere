package provider

import (
	"context"
	"errors"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/keys"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

var errPrivateKeyNotFound = errors.New("Private key no longer exists")

func (r *privateKeyResource) read(ctx context.Context, data privateKeyModel) (*sdk.ReadPrivateKeyResponse, error) {
	remote, err := r.client.Keys.Private.Get(ctx, &keys.GetPrivateRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), APIKeyID: data.Id.ValueString()})
	if err != nil {
		var missing *sdk.NotFoundError
		if errors.As(err, &missing) {
			if body, ok := missing.Body.(map[string]interface{}); ok && body["code"] == "SURFACE_KEY_NOT_FOUND" {
				return nil, errPrivateKeyNotFound
			}
		}
		return nil, errors.New(apiError(err))
	}
	if remote == nil || remote.ID != data.Id.ValueString() || remote.Version == "" || remote.Name == "" || len(remote.Scopes) == 0 {
		return nil, errors.New("The API returned incomplete or unexpected key metadata. Terraform state was preserved.")
	}
	if !remote.Revoked && !data.Secret.IsNull() && !data.Version.IsNull() && !data.Version.IsUnknown() && remote.Version != data.Version.ValueString() {
		return nil, errors.New("This key was rotated outside Terraform. Its stored secret is stale. Import the key to manage metadata only, or revoke and replace it using a new idempotency_key.")
	}
	return remote, nil
}

func (r *privateKeyResource) refresh(ctx context.Context, data *privateKeyModel, remote *sdk.ReadPrivateKeyResponse, diagnostics *diag.Diagnostics) {
	data.Name = types.StringValue(remote.Name)
	data.Version = types.StringValue(remote.Version)
	data.SurfaceSlug = types.StringPointerValue(remote.SurfaceSlug)
	data.ExpiresAt = types.Int64Null()
	data.SecondsUntilExpiration = types.Int64Null()
	if remote.SecondsUntilExpiration != nil {
		data.SecondsUntilExpiration = types.Int64Value(int64(*remote.SecondsUntilExpiration))
	}
	if remote.ExpiresAt != nil {
		data.ExpiresAt = types.Int64Value(int64(*remote.ExpiresAt))
	}
	scopes := make([]string, len(remote.Scopes))
	for i, scope := range remote.Scopes {
		scopes[i] = string(scope)
	}
	var existing []string
	if !data.Scopes.IsNull() && !data.Scopes.IsUnknown() {
		diagnostics.Append(data.Scopes.ElementsAs(ctx, &existing, false)...)
	}
	ordered := slices.Clone(existing)
	slices.Sort(ordered)
	sortedScopes := slices.Clone(scopes)
	slices.Sort(sortedScopes)
	if !slices.Equal(ordered, sortedScopes) || data.Scopes.IsNull() || data.Scopes.IsUnknown() {
		value, diags := types.ListValueFrom(ctx, types.StringType, scopes)
		diagnostics.Append(diags...)
		data.Scopes = value
	}
}

func (r *privateKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data privateKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if errors.Is(err, errPrivateKeyNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read private key", err.Error())
		return
	}
	if remote.Revoked {
		resp.State.RemoveResource(ctx)
		return
	}
	r.refresh(ctx, &data, remote, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

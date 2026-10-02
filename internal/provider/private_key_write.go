package provider

import (
	"context"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/keys"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *privateKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data privateKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var scopes []string
	resp.Diagnostics.Append(data.Scopes.ElementsAs(ctx, &scopes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := &keys.CreatePrivateRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), SurfaceSlug: data.SurfaceSlug.ValueStringPointer(), Name: data.Name.ValueString(), IdempotencyKey: data.IdempotencyKey.ValueString()}
	for _, scope := range scopes {
		request.Scopes = append(request.Scopes, keys.CreatePrivateRequestScopesItem(scope))
	}
	request.SetSecondsUntilExpiration(nil)
	if !data.SecondsUntilExpiration.IsNull() {
		seconds := int(data.SecondsUntilExpiration.ValueInt64())
		request.SetSecondsUntilExpiration(&seconds)
	}
	created, err := r.client.Keys.Private.Create(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create private key", apiError(err))
		return
	}
	if created == nil || created.APIKey == nil || created.APIKey.ID != data.IdempotencyKey.ValueString() || created.APIKey.Secret == "" {
		resp.Diagnostics.AddError("Invalid key creation response", "Retry with the same idempotency_key to recover the credential.")
		return
	}
	data.Id = types.StringValue(created.APIKey.ID)
	data.Version = types.StringNull()
	data.ExpiresAt = types.Int64Null()
	data.Secret = types.StringValue(created.APIKey.Secret)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to verify created key", err.Error())
		return
	}
	if remote.Revoked {
		resp.Diagnostics.AddError("Created key is inactive", "The API returned a revoked or expired credential.")
		return
	}
	r.refresh(ctx, &data, remote, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *privateKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data privateKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	result, err := r.client.Keys.Private.Delete(ctx, &keys.DeletePrivateRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), APIKeyID: data.Id.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Unable to revoke private key", apiError(err))
		return
	}
	if result == nil || !result.Success {
		resp.Diagnostics.AddError("Revocation was not acknowledged", "The API did not confirm revocation. Terraform state was preserved.")
	}
}

package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/keys"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
)

func (r *publicKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data publicKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.IdempotencyKey = creationIdentity(data.IdempotencyKey)
	created, err := r.client.Keys.Public.Create(ctx, &keys.CreatePublicRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), SurfaceSlug: data.SurfaceSlug.ValueString(), IdempotencyKey: data.IdempotencyKey.ValueString(), Name: data.Name.ValueString()}, option.WithMaxAttempts(3))
	if err != nil {
		resp.Diagnostics.AddError("Unable to create public key", creationError(data.IdempotencyKey, apiError(err)))
		return
	}
	if created == nil || created.Content == "" || created.Name != data.Name.ValueString() {
		resp.Diagnostics.AddError("Invalid key creation response", creationError(data.IdempotencyKey, "The API did not return the expected credential."))
		return
	}
	data.Id = data.IdempotencyKey
	data.Content = types.StringValue(created.Content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to verify created public key", err.Error())
		return
	}
	if remote.Revoked || remote.Name != data.Name.ValueString() {
		resp.Diagnostics.AddError("Created key is inactive or changed", "Refresh the key before retrying.")
		return
	}
	data.Content = types.StringValue(remote.Content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *publicKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data publicKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if errors.Is(err, errPublicKeyNotFound) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read public key before revocation", err.Error())
		return
	}
	if remote.Revoked {
		return
	}
	result, err := r.client.Keys.Public.Delete(ctx, &keys.DeletePublicRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), SurfaceSlug: data.SurfaceSlug.ValueString(), PublicKeyID: data.Id.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Unable to revoke public key", apiError(err))
		return
	}
	if result == nil || !result.Success {
		resp.Diagnostics.AddError("Revocation was not acknowledged", "The API did not confirm revocation. Terraform state was preserved.")
	}
}

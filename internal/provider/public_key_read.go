package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/keys"
)

var errPublicKeyNotFound = errors.New("Public key no longer exists")

func (r *publicKeyResource) read(ctx context.Context, data publicKeyModel) (*sdk.ReadPublicKeyResponse, error) {
	remote, err := r.client.Keys.Public.Get(ctx, &keys.GetPublicRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), SurfaceSlug: data.SurfaceSlug.ValueString(), PublicKeyID: data.Id.ValueString()})
	if err != nil {
		var missing *sdk.NotFoundError
		if errors.As(err, &missing) {
			if body, ok := missing.Body.(map[string]interface{}); ok && body["code"] == "SURFACE_KEY_NOT_FOUND" {
				return nil, errPublicKeyNotFound
			}
		}
		return nil, errors.New(apiError(err))
	}
	if remote == nil || remote.ID != data.Id.ValueString() || remote.SurfaceSlug != data.SurfaceSlug.ValueString() || remote.Name == "" || remote.Content == "" {
		return nil, errors.New("The API returned incomplete or unexpected public-key metadata. Terraform state was preserved.")
	}
	return remote, nil
}

func (r *publicKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data publicKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if errors.Is(err, errPublicKeyNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read public key", err.Error())
		return
	}
	if remote.Revoked {
		resp.State.RemoveResource(ctx)
		return
	}
	data.Name = types.StringValue(remote.Name)
	data.Content = types.StringValue(remote.Content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

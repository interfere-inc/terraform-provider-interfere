package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *publicKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 || !slugPattern.MatchString(parts[0]) || !slugPattern.MatchString(parts[1]) || !uuidPattern.MatchString(parts[2]) {
		resp.Diagnostics.AddError("Invalid public key import ID", "Use workspace-slug/surface-slug/key-uuid.")
		return
	}
	data := publicKeyModel{WorkspaceSlug: types.StringValue(parts[0]), SurfaceSlug: types.StringValue(parts[1]), Id: types.StringValue(parts[2]), IdempotencyKey: types.StringValue(parts[2])}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to import public key", err.Error())
		return
	}
	if remote.Revoked {
		resp.Diagnostics.AddError("Cannot import revoked public key", "Import an active public key.")
		return
	}
	data.Name = types.StringValue(remote.Name)
	data.Content = types.StringValue(remote.Content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

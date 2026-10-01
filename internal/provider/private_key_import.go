package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *privateKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || !slugPattern.MatchString(parts[0]) || !uuidPattern.MatchString(parts[1]) {
		resp.Diagnostics.AddError("Invalid private key import ID", "Use workspace-slug/key-uuid.")
		return
	}
	var data privateKeyModel
	data.Id = types.StringValue(parts[1])
	data.IdempotencyKey = types.StringValue(parts[1])
	data.WorkspaceSlug = types.StringValue(parts[0])
	data.Scopes = types.ListNull(types.StringType)

	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to import private key", err.Error())
		return
	}
	if remote.Revoked {
		resp.Diagnostics.AddError("Cannot import inactive private key", "Import a key that has not expired or been revoked.")
		return
	}
	r.refresh(ctx, &data, remote, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

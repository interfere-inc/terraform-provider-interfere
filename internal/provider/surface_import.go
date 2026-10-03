package provider

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *surfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	slug := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	if len(parts) != 2 || !slug.MatchString(parts[0]) || !slug.MatchString(parts[1]) {
		resp.Diagnostics.AddError("Invalid import identifier", "Use workspace-slug/surface-slug.")
		return
	}
	data := surfaceModel{WorkspaceSlug: types.StringValue(parts[0]), Slug: types.StringValue(parts[1])}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to import surface", err.Error())
		return
	}
	if remote.DeletedAt != nil {
		resp.Diagnostics.AddError("Cannot import deleted surface", "Import an active surface.")
		return
	}
	data.Id = types.StringValue(remote.ID)
	data.AnonymousUserTracking = types.BoolValue(remote.AnonymousUserTracking)
	data.Name = types.StringValue(remote.Name)
	data.Type = types.StringValue(string(remote.Type))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

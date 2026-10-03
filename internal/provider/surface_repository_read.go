package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func (r *surfaceRepositoryResource) read(ctx context.Context, data surfaceRepositoryModel) (*sdk.QuerySurfacesGetBySlugIncludeDeletedResponse, error) {
	return (&surfaceResource{client: r.client}).read(ctx, surfaceModel{Id: data.ID, WorkspaceSlug: data.WorkspaceSlug, Slug: data.SurfaceSlug})
}

func (data *surfaceRepositoryModel) refresh(remote *sdk.QuerySurfacesGetBySlugIncludeDeletedResponse) {
	data.ID = types.StringValue(remote.ID)
	data.IntegrationID = types.StringPointerValue(remote.SourceIntegrationID)
	data.MappingID = types.StringPointerValue(remote.SourceMappingID)
	data.WorkingDirectory = types.StringPointerValue(remote.SourceWorkingDirectory)
}

func (r *surfaceRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data surfaceRepositoryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read mapping", err.Error())
		return
	}
	if remote.DeletedAt != nil || (remote.SourceIntegrationID == nil && remote.SourceMappingID == nil) {
		resp.State.RemoveResource(ctx)
		return
	}
	if remote.SourceIntegrationID == nil || remote.SourceMappingID == nil {
		resp.Diagnostics.AddError("Invalid mapping", "The API returned an incomplete mapping. State was preserved.")
		return
	}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

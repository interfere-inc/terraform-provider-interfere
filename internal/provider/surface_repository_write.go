package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func (r *surfaceRepositoryResource) write(ctx context.Context, data *surfaceRepositoryModel, creating bool) error {
	remote, err := r.read(ctx, *data)
	if err != nil {
		return err
	}
	if remote.DeletedAt != nil {
		return errors.New("The surface was deleted.")
	}
	if creating && remote.SourceIntegrationID != nil && (*remote.SourceIntegrationID != data.IntegrationID.ValueString() || remote.SourceMappingID == nil || *remote.SourceMappingID != data.MappingID.ValueString() || types.StringPointerValue(remote.SourceWorkingDirectory) != data.WorkingDirectory) {
		return errors.New("The surface already has a different mapping. Import it before changing it.")
	}
	result, err := r.client.Integrations.LinkSurfaceToRepository(ctx, &sdk.LinkSurfaceToRepositoryIntegrationsRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.LinkSurfaceToRepositoryIntegrationsRequestArgs{IntegrationID: data.IntegrationID.ValueString(), SurfaceSlug: data.SurfaceSlug.ValueString(), RepositoryID: data.MappingID.ValueString(), WorkingDirectory: data.WorkingDirectory.ValueStringPointer()}})
	if err != nil {
		return errors.New(apiError(err))
	}
	if result == nil || !result.Success {
		return errors.New("The API did not acknowledge the mapping update.")
	}
	data.ID = types.StringValue(remote.ID)
	return nil
}

func (r *surfaceRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data surfaceRepositoryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.write(ctx, &data, true); err != nil {
		resp.Diagnostics.AddError("Unable to create mapping", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *surfaceRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data surfaceRepositoryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.write(ctx, &data, false); err != nil {
		resp.Diagnostics.AddError("Unable to update mapping", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *surfaceRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data surfaceRepositoryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete mapping", err.Error())
		return
	}
	if remote.DeletedAt != nil || (remote.SourceIntegrationID == nil && remote.SourceMappingID == nil) {
		return
	}
	if remote.SourceIntegrationID == nil || remote.SourceMappingID == nil || *remote.SourceIntegrationID != data.IntegrationID.ValueString() || *remote.SourceMappingID != data.MappingID.ValueString() {
		resp.Diagnostics.AddError("Mapping changed", "Refresh the plan before unlinking a changed mapping.")
		return
	}
	result, err := r.client.Integrations.UnlinkSurfaceMapping(ctx, &sdk.UnlinkSurfaceMappingIntegrationsRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.UnlinkSurfaceMappingIntegrationsRequestArgs{SurfaceSlug: data.SurfaceSlug.ValueString()}})
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete mapping", apiError(err))
		return
	}
	if result == nil || !result.Success {
		resp.Diagnostics.AddError("Deletion not acknowledged", "The API did not confirm unlinking. State was preserved.")
	}
}

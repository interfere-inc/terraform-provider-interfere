package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func (r *surfaceDestinationResource) write(ctx context.Context, data *surfaceDestinationModel, creating bool) error {
	remote, err := r.read(ctx, *data)
	if err != nil {
		return err
	}
	if remote.DeletedAt != nil {
		return errors.New("The surface was deleted.")
	}
	if creating && remote.DestinationIntegrationID != nil && (*remote.DestinationIntegrationID != data.IntegrationID.ValueString() || remote.DestinationMappingID == nil || *remote.DestinationMappingID != data.MappingID.ValueString()) {
		return errors.New("The surface already has a different mapping. Import it before changing it.")
	}
	result, err := r.client.Integrations.LinkSurfaceToDestination(ctx, &sdk.LinkSurfaceToDestinationIntegrationsRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.LinkSurfaceToDestinationIntegrationsRequestArgs{IntegrationID: data.IntegrationID.ValueString(), SurfaceSlug: data.SurfaceSlug.ValueString(), ProjectID: data.MappingID.ValueString()}})
	if err != nil {
		return errors.New(apiError(err))
	}
	if result == nil || !result.Success {
		return errors.New("The API did not acknowledge the mapping update.")
	}
	data.ID = types.StringValue(remote.ID)
	return nil
}

func (r *surfaceDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data surfaceDestinationModel
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

func (r *surfaceDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data surfaceDestinationModel
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

func (r *surfaceDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data surfaceDestinationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete mapping", err.Error())
		return
	}
	if remote.DeletedAt != nil || (remote.DestinationIntegrationID == nil && remote.DestinationMappingID == nil) {
		return
	}
	if remote.DestinationIntegrationID == nil || remote.DestinationMappingID == nil || *remote.DestinationIntegrationID != data.IntegrationID.ValueString() || *remote.DestinationMappingID != data.MappingID.ValueString() {
		resp.Diagnostics.AddError("Mapping changed", "Refresh the plan before unlinking a changed mapping.")
		return
	}
	result, err := r.client.Integrations.UnlinkSurfaceDestinationMapping(ctx, &sdk.UnlinkSurfaceDestinationMappingIntegrationsRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.UnlinkSurfaceDestinationMappingIntegrationsRequestArgs{SurfaceSlug: data.SurfaceSlug.ValueString()}})
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete mapping", apiError(err))
		return
	}
	if result == nil || !result.Success {
		resp.Diagnostics.AddError("Deletion not acknowledged", "The API did not confirm unlinking. State was preserved.")
	}
}

package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func (r *workspaceResource) apply(ctx context.Context, data *workspaceModel, previous workspaceModel) error {
	remote, err := r.read(ctx, previous)
	if err != nil {
		return err
	}
	if remote.Name != data.Name.ValueString() || remote.Slug != data.WorkspaceSlug.ValueString() {
		name, slug := data.Name.ValueString(), data.WorkspaceSlug.ValueString()
		result, err := r.client.Workspaces.UpdateBasics(ctx, &sdk.UpdateBasicsWorkspacesRequest{
			WorkspaceSlug: previous.WorkspaceSlug.ValueString(),
			Args:          &sdk.UpdateBasicsWorkspacesRequestArgs{Name: &name, Slug: &slug},
		})
		if err != nil {
			return errors.New(apiError(err))
		}
		if result == nil || !result.Success {
			return errors.New("The API did not acknowledge the settings update. Terraform state was preserved.")
		}
	}
	data.Id = types.StringValue(remote.ID)
	data.DataResidencyLocation = types.StringValue(string(remote.DataResidencyLocation))
	return nil
}

func (r *workspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data workspaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.apply(ctx, &data, data); err != nil {
		resp.Diagnostics.AddError("Unable to manage workspace", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *workspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, previous workspaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previous)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.apply(ctx, &data, previous); err != nil {
		resp.Diagnostics.AddError("Unable to update workspace", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

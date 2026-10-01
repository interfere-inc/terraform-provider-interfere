package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_workspace"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func (r *workspaceResource) read(ctx context.Context, data resource_workspace.WorkspaceModel) (*sdk.QueryOrganizationsCurrentResponse, error) {
	remote, err := r.client.Workspaces.QueryOrganizationsCurrent(ctx, &sdk.QueryOrganizationsCurrentRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString()})
	if err != nil {
		return nil, errors.New(apiError(err))
	}
	if remote == nil || remote.ID == "" || remote.Name == "" || remote.Slug != data.WorkspaceSlug.ValueString() || remote.DataResidencyLocation == "" {
		return nil, errors.New("The API returned no workspace or an invalid workspace identity. Verify access; Terraform state was preserved.")
	}
	if !data.Id.IsNull() && !data.Id.IsUnknown() && remote.ID != data.Id.ValueString() {
		return nil, errors.New("This slug identifies a different workspace. Terraform state was preserved.")
	}
	return remote, nil
}

func (r *workspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_workspace.WorkspaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read workspace", err.Error())
		return
	}
	data.Id = types.StringValue(remote.ID)
	data.Name = types.StringValue(remote.Name)
	data.DataResidencyLocation = types.StringValue(string(remote.DataResidencyLocation))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

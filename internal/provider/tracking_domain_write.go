package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
)

func (r *trackingDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data trackingDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.IdempotencyKey = creationIdentity(data.IdempotencyKey)
	result, err := r.client.Workspaces.AddProxyDomain(ctx, &sdk.AddProxyDomainWorkspacesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.AddProxyDomainWorkspacesRequestArgs{ID: data.IdempotencyKey.ValueString(), Name: data.Name.ValueString()}}, option.WithMaxAttempts(3))
	if err != nil {
		resp.Diagnostics.AddError("Unable to create domain", creationError(data.IdempotencyKey, apiError(err)))
		return
	}
	if result == nil || !result.Success {
		resp.Diagnostics.AddError("Creation not acknowledged", creationError(data.IdempotencyKey, "The API did not confirm domain creation."))
		return
	}
	data.ID = data.IdempotencyKey
	data.Status = types.StringValue("pending")
	data.SSLStatus = types.StringNull()
	data.CNAMETarget = types.StringValue("proxy.interfere.domains")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read created domain", err.Error())
		return
	}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *trackingDomainResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Domain replacement required", "Domain name and workspace changes require replacement.")
}

func (r *trackingDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data trackingDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete domain", err.Error())
		return
	}
	if remote.DeletedAt != nil {
		return
	}
	result, err := r.client.Workspaces.RemoveProxyDomain(ctx, &sdk.RemoveProxyDomainWorkspacesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.RemoveProxyDomainWorkspacesRequestArgs{ID: data.ID.ValueStringPointer(), Name: data.Name.ValueString()}})
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete domain", apiError(err))
		return
	}
	if result == nil || !result.Success {
		resp.Diagnostics.AddError("Deletion not acknowledged", "The API did not confirm domain removal. Terraform state was preserved.")
	}
}

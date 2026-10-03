package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func (r *trackingDomainResource) read(ctx context.Context, data trackingDomainModel) (*sdk.QueryDomainsByIDIncludeDeletedResponse, error) {
	remote, err := r.client.Domains.ByIDIncludeDeleted(ctx, &sdk.ByIDIncludeDeletedDomainsRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.ByIDIncludeDeletedDomainsRequestArgs{DomainID: data.ID.ValueString()}})
	if err != nil {
		return nil, errors.New(apiError(err))
	}
	if remote == nil || remote.ID != data.ID.ValueString() || remote.Name == "" || remote.Type == nil || *remote.Type != "proxy" || (remote.DeletedAt == nil && remote.Status == nil) {
		return nil, errors.New("The API returned no domain or an invalid identity. Check workspace domain permissions; Terraform state was preserved.")
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() && remote.Name != data.Name.ValueString() {
		return nil, errors.New("The domain identity now has a different hostname. Terraform state was preserved.")
	}
	return remote, nil
}

func (data *trackingDomainModel) refresh(remote *sdk.QueryDomainsByIDIncludeDeletedResponse) {
	data.Name = types.StringValue(remote.Name)
	data.Status = types.StringValue(string(*remote.Status))
	data.CNAMETarget = types.StringValue("proxy.interfere.domains")
	data.SSLStatus = types.StringNull()
	if remote.DomainMetadata != nil {
		data.SSLStatus = types.StringPointerValue(remote.DomainMetadata.CfSslStatus)
	}
}

func (r *trackingDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data trackingDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read domain", err.Error())
		return
	}
	if remote.DeletedAt != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

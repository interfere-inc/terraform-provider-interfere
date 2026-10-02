package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/core"
)

func apiError(err error) string {
	var response *core.APIError
	if errors.As(err, &response) {
		return fmt.Sprintf("The Interfere API returned HTTP %d. Check workspace access and the resource configuration. The response body is omitted because it may contain credentials.", response.StatusCode)
	}
	return "The API request failed or returned an invalid response. Check connectivity and credentials. For an uncertain create result, retry with the same idempotency_key."
}

func (r *surfaceResource) read(ctx context.Context, data surfaceModel) (*sdk.QuerySurfacesGetBySlugIncludeDeletedResponse, error) {
	remote, err := r.client.Surfaces.GetBySlugIncludeDeleted(ctx, &sdk.GetBySlugIncludeDeletedSurfacesRequest{
		WorkspaceSlug: data.WorkspaceSlug.ValueString(),
		Args:          &sdk.GetBySlugIncludeDeletedSurfacesRequestArgs{SurfaceSlug: data.Slug.ValueString()},
	})
	if err != nil {
		return nil, errors.New(apiError(err))
	}
	if remote == nil {
		return nil, errors.New("The API returned no surface. This can mean insufficient permissions, so Terraform state was preserved. Verify workspace access before retrying.")
	}
	if remote.ID == "" || remote.Name == "" || remote.Type == "" || remote.Slug != data.Slug.ValueString() {
		return nil, errors.New("The API returned an incomplete or unexpected surface. Terraform state was preserved.")
	}
	if !data.Id.IsNull() && remote.ID != data.Id.ValueString() {
		return nil, errors.New("This slug now identifies a different surface. Terraform state was preserved to avoid managing the wrong resource.")
	}
	return remote, nil
}

func (r *surfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data surfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read surface", err.Error())
		return
	}
	if remote.DeletedAt != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	data.Id = types.StringValue(remote.ID)
	data.Name = types.StringValue(remote.Name)
	data.Type = types.StringValue(string(remote.Type))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

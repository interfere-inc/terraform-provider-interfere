package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func (r *surfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data surfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	grants := &sdk.CreateSurfaceRequestAPIKey{Scopes: []sdk.CreateSurfaceRequestAPIKeyScopesItem{sdk.CreateSurfaceRequestAPIKeyScopesItemReleaseWrite}}
	grants.SetSecondsUntilExpiration(nil)
	created, err := r.client.Surfaces.CreateSurface(ctx, &sdk.CreateSurfaceRequest{
		APIKey:         grants,
		WorkspaceSlug:  data.WorkspaceSlug.ValueString(),
		IdempotencyKey: data.IdempotencyKey.ValueString(),
		Name:           data.Name.ValueString(),
		Type:           sdk.CreateSurfaceRequestType(data.Type.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create surface", apiError(err))
		return
	}
	if created == nil || created.Surface == nil || created.Surface.ID == "" || created.Surface.Slug == "" {
		resp.Diagnostics.AddError("Invalid creation response", "The API did not return a surface identity. Retry with the same idempotency_key to recover the creation result.")
		return
	}
	data.Id = types.StringValue(created.Surface.ID)
	data.Slug = types.StringValue(created.Surface.Slug)
	data.Name = types.StringValue(created.Surface.Name)
	data.Type = types.StringValue(string(created.Surface.Type))
	partial := data
	if partial.AnonymousUserTracking.IsUnknown() {
		partial.AnonymousUserTracking = types.BoolNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &partial)...)
	if err := r.applySettings(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Unable to configure surface", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *surfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data surfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update surface", err.Error())
		return
	}
	if remote.DeletedAt != nil {
		resp.Diagnostics.AddError("Surface was deleted", "Refresh the plan before recreating this surface.")
		return
	}
	if remote.Name != data.Name.ValueString() {
		updated, err := r.client.Surfaces.UpdateName(ctx, &sdk.UpdateNameSurfacesRequest{
			WorkspaceSlug: data.WorkspaceSlug.ValueString(),
			Args:          &sdk.UpdateNameSurfacesRequestArgs{SurfaceSlug: data.Slug.ValueString(), Name: data.Name.ValueString()},
		})
		if err != nil {
			resp.Diagnostics.AddError("Unable to update surface", apiError(err))
			return
		}
		if updated == nil || !updated.Success {
			resp.Diagnostics.AddError("Update was not acknowledged", "The API did not confirm the rename. Refresh the plan before retrying.")
			return
		}
	}
	if err := r.applySettings(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Unable to configure surface", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *surfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data surfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete surface", err.Error())
		return
	}
	if remote.DeletedAt != nil {
		return
	}
	deleted, err := r.client.Surfaces.Delete(ctx, &sdk.DeleteSurfacesRequest{
		WorkspaceSlug: data.WorkspaceSlug.ValueString(),
		Args:          &sdk.DeleteSurfacesRequestArgs{SurfaceSlug: data.Slug.ValueString()},
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete surface", apiError(err))
		return
	}
	if deleted == nil || !deleted.Success {
		resp.Diagnostics.AddError("Deletion was not acknowledged", "The API did not confirm deletion. Terraform state was preserved.")
	}
}

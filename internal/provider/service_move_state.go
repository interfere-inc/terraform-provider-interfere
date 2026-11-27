package provider

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

const legacyServiceStateKey = "legacy_surface_identity"

type legacyServiceIdentity struct {
	ID            string `json:"id"`
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	WorkspaceSlug string `json:"workspace_slug"`
}

func (r *serviceResource) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{{StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
		if req.SourceTypeName != "interfere_surface" || req.SourceSchemaVersion != 0 || !strings.HasSuffix(req.SourceProviderAddress, "/interfere-inc/interfere") {
			return
		}
		var previous legacyServiceIdentity
		if req.SourceRawState == nil || json.Unmarshal(req.SourceRawState.JSON, &previous) != nil || !uuidPattern.MatchString(previous.ID) || !slugPattern.MatchString(previous.Slug) || !slugPattern.MatchString(previous.WorkspaceSlug) {
			resp.Diagnostics.AddError("Invalid legacy surface state", "The stored App identity is incomplete. Source state was preserved.")
			return
		}
		payload, err := json.Marshal(previous)
		if err != nil {
			resp.Diagnostics.AddError("Invalid legacy surface state", "The stored App identity could not be preserved.")
			return
		}
		resp.Diagnostics.Append(resp.TargetPrivate.SetKey(ctx, legacyServiceStateKey, payload)...)
		if resp.Diagnostics.HasError() {
			return
		}
		data := serviceModel{ID: types.StringValue(previous.ID), WorkspaceSlug: types.StringValue(previous.WorkspaceSlug), Name: types.StringValue(previous.Name), Namespace: types.StringNull(), Type: types.StringNull(), WebsiteSettings: types.ObjectNull(serviceSettingTypes(ctx))}
		resp.Diagnostics.Append(resp.TargetState.Set(ctx, &data)...)
	}}}
}

func (r *serviceResource) resolveLegacyIdentity(ctx context.Context, data *serviceModel, payload []byte) error {
	var previous legacyServiceIdentity
	if json.Unmarshal(payload, &previous) != nil || previous.ID != data.ID.ValueString() || previous.WorkspaceSlug != data.WorkspaceSlug.ValueString() || !slugPattern.MatchString(previous.Slug) {
		return errors.New("The stored App identity does not match this resource. Terraform state was preserved.")
	}
	app, err := r.client.Surfaces.GetBySlugIncludeDeleted(ctx, &sdk.GetBySlugIncludeDeletedSurfacesRequest{WorkspaceSlug: previous.WorkspaceSlug, Args: &sdk.GetBySlugIncludeDeletedSurfacesRequestArgs{SurfaceSlug: previous.Slug}})
	if err != nil {
		return errors.New(apiError(err))
	}
	if app == nil || app.ID != previous.ID || app.Slug != previous.Slug || app.ServiceID == nil || !uuidPattern.MatchString(*app.ServiceID) {
		return errors.New("The stored App has no canonical service alias or its identity changed. Terraform state was preserved.")
	}
	data.ID = types.StringValue(*app.ServiceID)
	return nil
}

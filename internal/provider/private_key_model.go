package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type privateKeyModel struct {
	ExpiresAt              types.Int64  `tfsdk:"expires_at"`
	Id                     types.String `tfsdk:"id"`
	IdempotencyKey         types.String `tfsdk:"idempotency_key"`
	Name                   types.String `tfsdk:"name"`
	Scopes                 types.List   `tfsdk:"scopes"`
	SecondsUntilExpiration types.Int64  `tfsdk:"seconds_until_expiration"`
	SurfaceSlug            types.String `tfsdk:"surface_slug"`
	Version                types.String `tfsdk:"version"`
	WorkspaceSlug          types.String `tfsdk:"workspace_slug"`
	Secret                 types.String `tfsdk:"secret"`
}

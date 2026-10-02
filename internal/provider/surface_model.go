package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type surfaceModel struct {
	Id             types.String `tfsdk:"id"`
	IdempotencyKey types.String `tfsdk:"idempotency_key"`
	Name           types.String `tfsdk:"name"`
	Slug           types.String `tfsdk:"slug"`
	Type           types.String `tfsdk:"type"`
	WorkspaceSlug  types.String `tfsdk:"workspace_slug"`
}

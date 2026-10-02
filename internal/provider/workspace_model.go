package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type workspaceModel struct {
	DataResidencyLocation types.String `tfsdk:"data_residency_location"`
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	WorkspaceSlug         types.String `tfsdk:"workspace_slug"`
}

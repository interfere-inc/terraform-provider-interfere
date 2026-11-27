package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type workspaceModel struct {
	AllowMemberWorkspaceCreation types.Bool   `tfsdk:"allow_member_workspace_creation"`
	DataResidencyLocation        types.String `tfsdk:"data_residency_location"`
	EnabledAuthMethods           types.Set    `tfsdk:"enabled_auth_methods"`
	Id                           types.String `tfsdk:"id"`
	ImpersonationExpiresAt       types.Int64  `tfsdk:"impersonation_expires_at"`
	ImpersonationGranted         types.Bool   `tfsdk:"impersonation_granted"`
	Name                         types.String `tfsdk:"name"`
	WorkspaceSlug                types.String `tfsdk:"workspace_slug"`
}

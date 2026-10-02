package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type workspaceResource struct{ client *client.Client }

func NewWorkspaceResource() resource.Resource { return &workspaceResource{} }

func (r *workspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (r *workspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	configured, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Invalid provider configuration", "Expected an Interfere API client.")
		return
	}
	r.client = configured
}

func (r *workspaceResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
}

func (r *workspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !slugPattern.MatchString(req.ID) {
		resp.Diagnostics.AddError("Invalid workspace import ID", "Use the existing workspace slug.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_slug"), req.ID)...)
}

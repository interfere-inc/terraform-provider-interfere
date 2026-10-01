package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_workspace"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type workspaceResource struct{ client *client.Client }

func NewWorkspaceResource() resource.Resource { return &workspaceResource{} }

func (r *workspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (r *workspaceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_workspace.WorkspaceResourceSchema(ctx)
	resp.Schema.Description = "Manage an existing workspace's name and slug. Applying adopts the workspace; destroying only removes it from Terraform state and never deletes the workspace. Import using its slug."
	for _, key := range []string{"name", "workspace_slug"} {
		field := resp.Schema.Attributes[key].(schema.StringAttribute)
		field.Required, field.Optional, field.Computed = true, false, false
		resp.Schema.Attributes[key] = field
	}
	name := resp.Schema.Attributes["name"].(schema.StringAttribute)
	name.Validators = []validator.String{stringvalidator.LengthBetween(1, 255)}
	resp.Schema.Attributes["name"] = name
	for _, key := range []string{"id", "data_residency_location"} {
		field := resp.Schema.Attributes[key].(schema.StringAttribute)
		field.PlanModifiers = []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
		resp.Schema.Attributes[key] = field
	}
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

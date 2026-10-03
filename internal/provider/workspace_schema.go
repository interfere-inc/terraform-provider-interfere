package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_workspace"
)

func (r *workspaceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_workspace.WorkspaceResourceSchema(ctx)
	resp.Schema.Description = "Manage an existing workspace's name and slug. Applying adopts the workspace; destroying only removes it from Terraform state and never deletes the workspace. Import using its slug."
	for _, field := range []string{"id", "data_residency_location"} {
		attribute := resp.Schema.Attributes[field].(schema.StringAttribute)
		attribute.PlanModifiers = []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
		resp.Schema.Attributes[field] = attribute
	}
	name := resp.Schema.Attributes["name"].(schema.StringAttribute)
	name.Required, name.Computed = true, false
	name.Validators = []validator.String{stringvalidator.LengthBetween(1, 255)}
	resp.Schema.Attributes["name"] = name
	slug := resp.Schema.Attributes["workspace_slug"].(schema.StringAttribute)
	slug.Required, slug.Optional, slug.Computed = true, false, false
	resp.Schema.Attributes["workspace_slug"] = slug
}

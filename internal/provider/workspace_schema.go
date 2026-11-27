package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_workspace"
)

func (r *workspaceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_workspace.WorkspaceResourceSchema(ctx)
	resp.Schema.Description = "Manage an existing workspace's name, slug, and security settings. Applying adopts the workspace; destroying only removes it from Terraform state and never deletes the workspace or resets its settings. Omitted security settings keep their current values. Import using its slug."
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
	for _, field := range []string{"allow_member_workspace_creation", "impersonation_granted"} {
		attribute := resp.Schema.Attributes[field].(schema.BoolAttribute)
		attribute.Optional = true
		attribute.PlanModifiers = []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}
		resp.Schema.Attributes[field] = attribute
	}
	methods := resp.Schema.Attributes["enabled_auth_methods"].(schema.ListAttribute)
	resp.Schema.Attributes["enabled_auth_methods"] = schema.SetAttribute{
		ElementType:         types.StringType,
		Optional:            true,
		Computed:            true,
		Description:         methods.Description,
		MarkdownDescription: methods.MarkdownDescription,
		PlanModifiers:       []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		Validators: []validator.Set{
			setvalidator.ValueStringsAre(stringvalidator.OneOf(authMethods...)),
			alwaysEnabledAuthMethodsValidator{},
		},
	}
	expiry := resp.Schema.Attributes["impersonation_expires_at"].(schema.NumberAttribute)
	resp.Schema.Attributes["impersonation_expires_at"] = schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		Description:         expiry.Description,
		MarkdownDescription: expiry.MarkdownDescription,
		PlanModifiers:       []planmodifier.Int64{supportAccessExpiryModifier{}},
		Validators:          []validator.Int64{int64validator.AtLeast(1)},
	}
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *surfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a surface's name and framework. Import using workspace-slug/surface-slug. Creation also issues default credentials, which this resource does not store.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"slug":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":           schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.LengthBetween(1, 48)}},
			"type":           schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.OneOf("elysia", "nest", "nextjs", "python", "react")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"workspace_slug": schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(slugPattern, "Must be a workspace slug.")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"idempotency_key": schema.StringAttribute{
				Required:    true,
				Description: "Nonzero UUID for this creation attempt. Reuse it after an uncertain create result; choose a new UUID when replacing or recreating the surface.",
				Validators:  []validator.String{stringvalidator.RegexMatches(uuidPattern, "Must be a UUID."), stringvalidator.NoneOf("00000000-0000-0000-0000-000000000000")},
			},
		},
	}
}

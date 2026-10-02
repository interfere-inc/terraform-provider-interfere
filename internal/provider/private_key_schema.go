package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *privateKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a workspace or surface private API key. Optional surface_slug selects a surface. Grants and lifetime changes replace the key; choose a fresh idempotency_key for each replacement. secret is sensitive and is stored in Terraform state. Import retrieves metadata only, never the secret.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"version":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":            schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.LengthBetween(1, 64)}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"workspace_slug":  schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(slugPattern, "Must be a workspace slug.")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"surface_slug":    schema.StringAttribute{Optional: true, Validators: []validator.String{stringvalidator.RegexMatches(slugPattern, "Must be a surface slug.")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"idempotency_key": schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(uuidPattern, "Must be a UUID."), stringvalidator.NoneOf("00000000-0000-0000-0000-000000000000")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"scopes":          schema.ListAttribute{Required: true, ElementType: types.StringType, Validators: []validator.List{listvalidator.SizeAtLeast(1), listvalidator.UniqueValues()}, PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()}},
			"seconds_until_expiration": schema.Int64Attribute{
				Optional:      true,
				Description:   "Positive lifetime in seconds. Omit or set null for no expiry. Changing this value requires a replacement with a fresh idempotency_key.",
				Validators:    []validator.Int64{int64validator.Between(1, 9007199254740991)},
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"expires_at": schema.Int64Attribute{Computed: true, Description: "Expiration as Unix milliseconds, or null for no expiry.", PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"secret":     schema.StringAttribute{Computed: true, Sensitive: true, Description: "Created key secret. Stored in Terraform state; unavailable after import.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		},
	}
}

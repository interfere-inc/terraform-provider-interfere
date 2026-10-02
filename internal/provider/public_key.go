package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type publicKeyModel struct {
	Id             types.String `tfsdk:"id"`
	IdempotencyKey types.String `tfsdk:"idempotency_key"`
	WorkspaceSlug  types.String `tfsdk:"workspace_slug"`
	SurfaceSlug    types.String `tfsdk:"surface_slug"`
	Name           types.String `tfsdk:"name"`
	Content        types.String `tfsdk:"content"`
}

type publicKeyResource struct{ client *client.Client }

func NewPublicKeyResource() resource.Resource { return &publicKeyResource{} }

func (r *publicKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_key"
}

func (r *publicKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a surface's publishable key. Configuration changes replace the key and require a fresh idempotency_key. Import using workspace-slug/surface-slug/key-uuid. External rotation refreshes content; destruction revokes the key and its rotation grace-period values.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"content":         schema.StringAttribute{Computed: true, Description: "Publishable credential for browser SDK configuration."},
			"name":            schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.LengthBetween(1, 64)}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"workspace_slug":  schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(slugPattern, "Must be a workspace slug.")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"surface_slug":    schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(slugPattern, "Must be a surface slug.")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"idempotency_key": schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(uuidPattern, "Must be a UUID."), stringvalidator.NoneOf("00000000-0000-0000-0000-000000000000")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *publicKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *publicKeyResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Public keys require replacement", "Choose a fresh idempotency_key when changing a public key.")
}

func (r *publicKeyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var prior, planned publicKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planned)...)
	if resp.Diagnostics.HasError() || planned.IdempotencyKey.IsUnknown() || !prior.IdempotencyKey.Equal(planned.IdempotencyKey) {
		return
	}
	if !prior.Name.Equal(planned.Name) || !prior.SurfaceSlug.Equal(planned.SurfaceSlug) || !prior.WorkspaceSlug.Equal(planned.WorkspaceSlug) {
		resp.Diagnostics.AddError("New key creation identity required", "Changing public-key configuration replaces the key. Set idempotency_key to a fresh UUID before applying.")
	}
}

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type trackingDomainResource struct{ client *client.Client }
type trackingDomainModel struct {
	ID             types.String `tfsdk:"id"`
	WorkspaceSlug  types.String `tfsdk:"workspace_slug"`
	Name           types.String `tfsdk:"name"`
	IdempotencyKey types.String `tfsdk:"idempotency_key"`
	Status         types.String `tfsdk:"status"`
	SSLStatus      types.String `tfsdk:"ssl_status"`
	CNAMETarget    types.String `tfsdk:"cname_target"`
}

func NewTrackingDomainResource() resource.Resource { return &trackingDomainResource{} }
func (r *trackingDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tracking_domain"
}

func (r *trackingDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *trackingDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Manage a workspace tracking hostname. Creation returns while DNS and TLS verification is pending. Configure a DNS-only CNAME to cname_target using your DNS provider. Import with workspace-slug/domain-UUID.", Attributes: map[string]schema.Attribute{
		"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"workspace_slug":  schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(slugPattern, "Must be a workspace slug.")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"name":            schema.StringAttribute{Required: true, Description: "Lowercase fully qualified tracking hostname, without a trailing dot.", Validators: []validator.String{trackingHostnameValidator{}}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"idempotency_key": schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(uuidPattern, "Must be a UUID."), stringvalidator.NoneOf("00000000-0000-0000-0000-000000000000")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "UUID for this domain creation. Choose a fresh UUID for replacement or recreation."},
		"status":          schema.StringAttribute{Computed: true}, "ssl_status": schema.StringAttribute{Computed: true}, "cname_target": schema.StringAttribute{Computed: true},
	}}
}

func (r *trackingDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || !slugPattern.MatchString(parts[0]) || !uuidPattern.MatchString(parts[1]) {
		resp.Diagnostics.AddError("Invalid import identifier", "Use workspace-slug/domain-UUID.")
		return
	}
	data := trackingDomainModel{WorkspaceSlug: types.StringValue(parts[0]), ID: types.StringValue(parts[1]), IdempotencyKey: types.StringValue(parts[1])}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to import domain", err.Error())
		return
	}
	if remote.DeletedAt != nil {
		resp.Diagnostics.AddError("Domain was deleted", "Only an active domain can be imported.")
		return
	}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *trackingDomainResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var state, plan trackingDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.IdempotencyKey.IsUnknown() || plan.Name.IsUnknown() || plan.WorkspaceSlug.IsUnknown() {
		return
	}
	if plan.IdempotencyKey.Equal(state.IdempotencyKey) && (!plan.Name.Equal(state.Name) || !plan.WorkspaceSlug.Equal(state.WorkspaceSlug)) {
		resp.Diagnostics.AddError("Replacement needs a new creation UUID", "Choose a new idempotency_key when changing the domain name or workspace.")
	}
}

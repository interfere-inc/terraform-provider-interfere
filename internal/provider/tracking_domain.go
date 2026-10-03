package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_tracking_domain"
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

func (r *trackingDomainResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_tracking_domain.TrackingDomainResourceSchema(ctx)
	resp.Schema.Description = "Manage a workspace tracking hostname. Creation returns while DNS and TLS verification is pending. Configure a DNS-only CNAME to cname_target using your DNS provider. Import with workspace-slug/domain-UUID."
	attributes := resp.Schema.Attributes
	attributes["id"] = stableString(attributes["id"])
	attributes["workspace_slug"] = requiredString(attributes["workspace_slug"], stringplanmodifier.RequiresReplace())
	name := requiredString(attributes["name"], stringplanmodifier.RequiresReplace())
	name.Validators = []validator.String{trackingHostnameValidator{}}
	attributes["name"] = name
	attributes["ssl_status"] = attributes["domain_metadata"].(schema.SingleNestedAttribute).Attributes["cf_ssl_status"]
	delete(attributes, "domain_metadata")
	attributes["cname_target"] = schema.StringAttribute{Computed: true, Description: "DNS-only CNAME target for the tracking hostname."}
	attributes["idempotency_key"] = creationIdentityAttribute(true)
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
	planCreationIdentity(ctx, req, resp, "name", "workspace_slug")
}

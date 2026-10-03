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

type surfaceDestinationResource struct{ client *client.Client }
type surfaceDestinationModel struct {
	ID            types.String `tfsdk:"id"`
	WorkspaceSlug types.String `tfsdk:"workspace_slug"`
	SurfaceSlug   types.String `tfsdk:"surface_slug"`
	IntegrationID types.String `tfsdk:"integration_id"`
	MappingID     types.String `tfsdk:"project_id"`
}

func NewSurfaceDestinationResource() resource.Resource { return &surfaceDestinationResource{} }
func (r *surfaceDestinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_surface_destination"
}

func (r *surfaceDestinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *surfaceDestinationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Manage a surface's destination mapping through an existing integration. Import with workspace-slug/surface-slug. Destroy unlinks the mapping without deleting the surface or integration.", Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"workspace_slug": schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(slugPattern, "Must be a workspace slug.")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"surface_slug":   schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.RegexMatches(slugPattern, "Must be a surface slug.")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"integration_id": schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.LengthAtLeast(1)}},
		"project_id":     schema.StringAttribute{Required: true, Validators: []validator.String{stringvalidator.LengthAtLeast(1)}},
	}}
}

func (r *surfaceDestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || !slugPattern.MatchString(parts[0]) || !slugPattern.MatchString(parts[1]) {
		resp.Diagnostics.AddError("Invalid import identifier", "Use workspace-slug/surface-slug.")
		return
	}
	data := surfaceDestinationModel{WorkspaceSlug: types.StringValue(parts[0]), SurfaceSlug: types.StringValue(parts[1])}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to import mapping", err.Error())
		return
	}
	if remote.DeletedAt != nil || remote.DestinationIntegrationID == nil || remote.DestinationMappingID == nil {
		resp.Diagnostics.AddError("Mapping not found", "Only an active mapping can be imported.")
		return
	}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

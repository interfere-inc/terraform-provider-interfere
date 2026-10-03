package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_surface_repository"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type surfaceRepositoryResource struct{ client *client.Client }
type surfaceRepositoryModel struct {
	ID               types.String `tfsdk:"id"`
	WorkspaceSlug    types.String `tfsdk:"workspace_slug"`
	SurfaceSlug      types.String `tfsdk:"surface_slug"`
	IntegrationID    types.String `tfsdk:"integration_id"`
	MappingID        types.String `tfsdk:"repository_id"`
	WorkingDirectory types.String `tfsdk:"working_directory"`
}

func NewSurfaceRepositoryResource() resource.Resource { return &surfaceRepositoryResource{} }
func (r *surfaceRepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_surface_repository"
}

func (r *surfaceRepositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *surfaceRepositoryResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_surface_repository.SurfaceRepositoryResourceSchema(ctx)
	resp.Schema.Description = "Manage a surface's repository mapping through an existing integration. Import with workspace-slug/surface-slug. Destroy unlinks the mapping without deleting the surface or integration."
	attributes := resp.Schema.Attributes
	flattenAttributes(attributes, "args")
	attributes["id"] = stableString(attributes["id"])
	for _, field := range []string{"workspace_slug", "surface_slug"} {
		attributes[field] = requiredString(attributes[field], stringplanmodifier.RequiresReplace())
	}
	for _, field := range []string{"integration_id", "repository_id"} {
		attribute := requiredString(attributes[field])
		attribute.Validators = []validator.String{stringvalidator.LengthAtLeast(1)}
		attributes[field] = attribute
	}
	directory := attributes["working_directory"].(schema.StringAttribute)
	directory.Computed = false
	attributes["working_directory"] = directory
}

func (r *surfaceRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || !slugPattern.MatchString(parts[0]) || !slugPattern.MatchString(parts[1]) {
		resp.Diagnostics.AddError("Invalid import identifier", "Use workspace-slug/surface-slug.")
		return
	}
	data := surfaceRepositoryModel{WorkspaceSlug: types.StringValue(parts[0]), SurfaceSlug: types.StringValue(parts[1])}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to import mapping", err.Error())
		return
	}
	if remote.DeletedAt != nil || remote.SourceIntegrationID == nil || remote.SourceMappingID == nil {
		resp.Diagnostics.AddError("Mapping not found", "Only an active mapping can be imported.")
		return
	}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

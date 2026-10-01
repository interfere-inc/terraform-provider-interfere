package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_surface"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type surfaceResource struct{ client *client.Client }

func NewSurfaceResource() resource.Resource { return &surfaceResource{} }

func (r *surfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_surface"
}

func (r *surfaceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_surface.SurfaceResourceSchema(ctx)
	resp.Schema.Description = "Manage a surface's name and framework. Import using workspace-slug/surface-slug. Creation also issues default credentials, which this resource does not store."
	for _, name := range []string{"id", "slug"} {
		attribute := resp.Schema.Attributes[name].(schema.StringAttribute)
		attribute.PlanModifiers = []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
		resp.Schema.Attributes[name] = attribute
	}
	for _, name := range []string{"type", "workspace_slug"} {
		attribute := resp.Schema.Attributes[name].(schema.StringAttribute)
		attribute.Optional, attribute.Computed, attribute.Required = false, false, true
		attribute.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
		resp.Schema.Attributes[name] = attribute
	}
	key := resp.Schema.Attributes["idempotency_key"].(schema.StringAttribute)
	key.Description = "Nonzero UUID for this creation attempt. Reuse it after an uncertain create result; choose a new UUID when replacing or recreating the surface."
	key.Validators = []validator.String{
		stringvalidator.RegexMatches(regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`), "Must be a UUID."),
		stringvalidator.NoneOf("00000000-0000-0000-0000-000000000000"),
	}
	resp.Schema.Attributes["idempotency_key"] = key
	name := resp.Schema.Attributes["name"].(schema.StringAttribute)
	name.Validators = append(name.Validators, stringvalidator.LengthBetween(1, 48))
	resp.Schema.Attributes["name"] = name
}

func (r *surfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

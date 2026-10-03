package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_public_key"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type publicKeyModel = resource_public_key.PublicKeyModel

type publicKeyResource struct{ client *client.Client }

func NewPublicKeyResource() resource.Resource { return &publicKeyResource{} }

func (r *publicKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_key"
}

func (r *publicKeyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_public_key.PublicKeyResourceSchema(ctx)
	resp.Schema.Description = "Manage a surface's publishable key. Configuration changes replace the key. Creation UUIDs are generated automatically unless explicitly configured. Import using workspace-slug/surface-slug/key-uuid. External rotation refreshes content; destruction revokes the key and its rotation grace-period values."
	attributes := resp.Schema.Attributes
	attributes["id"] = stableString(attributes["id"])
	for _, field := range []string{"name", "workspace_slug", "surface_slug"} {
		attributes[field] = requiredString(attributes[field], stringplanmodifier.RequiresReplace())
	}
	slug := attributes["surface_slug"].(schema.StringAttribute)
	slug.Validators = []validator.String{stringvalidator.RegexMatches(slugPattern, "Must be a surface slug.")}
	attributes["surface_slug"] = slug
	attributes["idempotency_key"] = creationIdentityAttribute(true)
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
	planCreationIdentity(ctx, req, resp, "name", "surface_slug", "workspace_slug")
}

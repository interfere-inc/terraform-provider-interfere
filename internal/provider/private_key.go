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
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_private_key"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type privateKeyModel struct {
	resource_private_key.PrivateKeyModel
	Secret types.String `tfsdk:"secret"`
}

type privateKeyResource struct{ client *client.Client }

func NewPrivateKeyResource() resource.Resource { return &privateKeyResource{} }

func (r *privateKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_key"
}

func (r *privateKeyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_private_key.PrivateKeyResourceSchema(ctx)
	resp.Schema.Description = "Manage a workspace or surface private API key. Optional surface_slug selects a surface. Grants and lifetime changes replace the key; choose a fresh idempotency_key for each replacement. secret is sensitive and is stored in Terraform state. Import retrieves metadata only, never the secret."
	for _, key := range []string{"name", "workspace_slug", "surface_slug", "idempotency_key"} {
		field := resp.Schema.Attributes[key].(schema.StringAttribute)
		field.Optional, field.Required, field.Computed = key == "surface_slug", key != "surface_slug", false
		field.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
		resp.Schema.Attributes[key] = field
	}
	identity := resp.Schema.Attributes["idempotency_key"].(schema.StringAttribute)
	identity.Validators = []validator.String{stringvalidator.RegexMatches(uuidPattern, "Must be a UUID."), stringvalidator.NoneOf("00000000-0000-0000-0000-000000000000")}
	resp.Schema.Attributes["idempotency_key"] = identity
	for _, key := range []string{"id", "version"} {
		field := resp.Schema.Attributes[key].(schema.StringAttribute)
		field.PlanModifiers = []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
		resp.Schema.Attributes[key] = field
	}
	scopes := resp.Schema.Attributes["scopes"].(schema.ListAttribute)
	scopes.Validators = append(scopes.Validators, listvalidator.UniqueValues())
	scopes.PlanModifiers = []planmodifier.List{listplanmodifier.RequiresReplace()}
	resp.Schema.Attributes["scopes"] = scopes
	lifetime := resp.Schema.Attributes["seconds_until_expiration"].(schema.Int64Attribute)
	lifetime.Required, lifetime.Optional = false, true
	lifetime.Description = "Positive lifetime in seconds. Omit or set null for no expiry. Changing this value requires a replacement with a fresh idempotency_key."
	lifetime.Validators = []validator.Int64{int64validator.Between(1, 9007199254740991)}
	lifetime.PlanModifiers = []planmodifier.Int64{int64planmodifier.RequiresReplace()}
	resp.Schema.Attributes["seconds_until_expiration"] = lifetime
	expires := resp.Schema.Attributes["expires_at"].(schema.Int64Attribute)
	expires.Description = "Expiration as Unix milliseconds, or null for no expiry."
	expires.PlanModifiers = []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}
	resp.Schema.Attributes["expires_at"] = expires
	resp.Schema.Attributes["secret"] = schema.StringAttribute{Computed: true, Sensitive: true, Description: "Created key secret. Stored in Terraform state; unavailable after import.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}}
}

func (r *privateKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *privateKeyResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Private keys require replacement", "Choose a fresh idempotency_key when changing a private key.")
}

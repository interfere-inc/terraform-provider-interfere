package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_private_key"
)

func (r *privateKeyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_private_key.PrivateKeyResourceSchema(ctx)
	resp.Schema.Description = "Manage a workspace or surface private API key. Optional surface_slug selects a surface. Grants and lifetime changes replace the key. Creation UUIDs are generated automatically unless explicitly configured. secret is sensitive and is stored in Terraform state. Import retrieves metadata only, never the secret."
	attributes := resp.Schema.Attributes
	flattenAttributes(attributes, "api_key")
	for _, field := range []string{"id", "version", "secret"} {
		attributes[field] = stableString(attributes[field])
	}
	for _, field := range []string{"name", "workspace_slug"} {
		attributes[field] = requiredString(attributes[field], stringplanmodifier.RequiresReplace())
	}
	slug := attributes["surface_slug"].(schema.StringAttribute)
	slug.Required, slug.Optional, slug.Computed = false, true, false
	slug.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
	attributes["surface_slug"] = slug
	scopes := attributes["scopes"].(schema.ListAttribute)
	scopes.Validators = append(scopes.Validators, listvalidator.UniqueValues())
	scopes.PlanModifiers = []planmodifier.List{listplanmodifier.RequiresReplace()}
	attributes["scopes"] = scopes
	lifetime := attributes["seconds_until_expiration"].(schema.Int64Attribute)
	lifetime.Required, lifetime.Optional, lifetime.Computed = false, true, false
	lifetime.Validators = []validator.Int64{int64validator.Between(1, 9007199254740991)}
	lifetime.PlanModifiers = []planmodifier.Int64{int64planmodifier.RequiresReplace()}
	lifetime.Description += " Omit for no expiry. Changing this value replaces the key."
	lifetime.MarkdownDescription = lifetime.Description
	attributes["seconds_until_expiration"] = lifetime
	expires := attributes["expires_at"].(schema.Int64Attribute)
	expires.PlanModifiers = []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}
	attributes["expires_at"] = expires
	secret := attributes["secret"].(schema.StringAttribute)
	secret.Sensitive = true
	secret.Description += " Stored in Terraform state; unavailable after import."
	secret.MarkdownDescription = secret.Description
	attributes["secret"] = secret
	attributes["idempotency_key"] = creationIdentityAttribute(true)
}

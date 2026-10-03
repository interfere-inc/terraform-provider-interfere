package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_surface"
)

func (r *surfaceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_surface.SurfaceResourceSchema(ctx)
	resp.Schema.Description = "Manage a surface's name, framework, and anonymous-user tracking. Import using workspace-slug/surface-slug. Creation also issues default credentials, which this resource does not store."
	attributes := resp.Schema.Attributes
	for _, field := range []string{"id", "slug"} {
		attributes[field] = stableString(attributes[field])
	}
	for _, field := range []string{"workspace_slug", "type"} {
		attributes[field] = requiredString(attributes[field], stringplanmodifier.RequiresReplace())
	}
	name := attributes["name"].(schema.StringAttribute)
	name.Validators = []validator.String{stringvalidator.LengthBetween(1, 48)}
	attributes["name"] = name
	tracking := attributes["anonymous_user_tracking"].(schema.BoolAttribute)
	tracking.Optional = true
	tracking.Description += " Omit to preserve the current API setting."
	tracking.MarkdownDescription = tracking.Description
	attributes["anonymous_user_tracking"] = tracking
	attributes["idempotency_key"] = creationIdentityAttribute(false)
}

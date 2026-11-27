package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_service"
)

func (r *serviceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_service.ServiceResourceSchema(ctx)
	resp.Schema.Description = "Manage an existing service's type and website settings. Applying adopts the service by ID. Destroying only removes Terraform management and preserves the service, settings, credentials, and telemetry. Import using workspace-slug/service-uuid."
	for _, field := range []string{"id", "workspace_slug"} {
		resp.Schema.Attributes[field] = requiredString(resp.Schema.Attributes[field], stringplanmodifier.RequiresReplace())
	}
	id := resp.Schema.Attributes["id"].(schema.StringAttribute)
	id.Validators = []validator.String{stringvalidator.RegexMatches(uuidPattern, "Use a service UUID.")}
	resp.Schema.Attributes["id"] = id
	typeAttribute := resp.Schema.Attributes["type"].(schema.StringAttribute)
	typeAttribute.Optional = true
	typeAttribute.Validators = []validator.String{stringvalidator.OneOf("unclassified", "other", "backend", "website")}
	resp.Schema.Attributes["type"] = typeAttribute
	settings := resp.Schema.Attributes["website_settings"].(schema.SingleNestedAttribute)
	settings.CustomType = nil
	settings.Optional = true
	policy := settings.Attributes["anonymous_replay_policy"].(schema.StringAttribute)
	policy.Required, policy.Computed = true, false
	policy.Validators = []validator.String{stringvalidator.OneOf("always", "on_error", "never")}
	settings.Attributes["anonymous_replay_policy"] = policy
	for _, field := range []string{"anonymous_user_tracking", "environment_filtering"} {
		attribute := settings.Attributes[field].(schema.BoolAttribute)
		attribute.Required, attribute.Computed = true, false
		settings.Attributes[field] = attribute
	}
	plugins := settings.Attributes["plugins"].(schema.SingleNestedAttribute)
	plugins.CustomType = nil
	plugins.Optional, plugins.Computed = true, false
	for field, value := range plugins.Attributes {
		attribute := value.(schema.BoolAttribute)
		attribute.Optional, attribute.Computed = true, false
		plugins.Attributes[field] = attribute
	}
	settings.Attributes["plugins"] = plugins
	resp.Schema.Attributes["website_settings"] = settings
}

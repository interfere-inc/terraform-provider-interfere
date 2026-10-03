package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func requiredString(attribute schema.Attribute, modifiers ...planmodifier.String) schema.StringAttribute {
	value := attribute.(schema.StringAttribute)
	value.Required, value.Optional, value.Computed = true, false, false
	value.PlanModifiers = modifiers
	return value
}

func stableString(attribute schema.Attribute) schema.StringAttribute {
	value := attribute.(schema.StringAttribute)
	value.PlanModifiers = []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
	return value
}

func flattenAttributes(attributes map[string]schema.Attribute, field string) {
	for name, attribute := range attributes[field].(schema.SingleNestedAttribute).Attributes {
		attributes[name] = attribute
	}
	delete(attributes, field)
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type apiDataSource struct{ client *client.Client }

func (r *apiDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func plainSetAttribute(attributes map[string]schema.Attribute, field, description string) {
	set := attributes[field].(schema.SetNestedAttribute)
	set.Description, set.MarkdownDescription = description, description
	set.NestedObject.CustomType = nil
	attributes[field] = set
}

func objectSet(ctx context.Context, attributeTypes map[string]attr.Type, items any) (types.Set, bool) {
	value, diagnostics := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attributeTypes}, items)
	return value, !diagnostics.HasError()
}

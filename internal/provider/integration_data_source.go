package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/datasource_integration"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type integrationDataSource struct{ client *client.Client }
type integrationLookupModel struct {
	ID            types.String `tfsdk:"id"`
	WorkspaceSlug types.String `tfsdk:"workspace_slug"`
	Provider      types.String `tfsdk:"integration_provider"`
	IntegrationID types.String `tfsdk:"integration_id"`
	Status        types.String `tfsdk:"status"`
}

func NewIntegrationDataSource() datasource.DataSource { return &integrationDataSource{} }
func (r *integrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *integrationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *integrationDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_integration.IntegrationDataSourceSchema(ctx)
	resp.Schema.Description = "Look up an existing integration without managing its lifecycle."
	attributes := resp.Schema.Attributes
	installation := attributes["integration"].(schema.SetNestedAttribute).NestedObject.Attributes
	attributes["id"], attributes["status"] = installation["id"], installation["status"]
	provider := installation["provider"].(schema.StringAttribute)
	provider.Required, provider.Computed = true, false
	attributes["integration_provider"] = provider
	selector := installation["id"].(schema.StringAttribute)
	selector.Optional, selector.Computed = true, false
	selector.Description += " Select a specific installation when the provider has multiple installations."
	selector.MarkdownDescription = selector.Description
	attributes["integration_id"] = selector
	delete(attributes, "integration")
}

func (r *integrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data integrationLookupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	items, err := r.client.Integrations.Installations(ctx, &sdk.InstallationsIntegrationsRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.InstallationsIntegrationsRequestArgs{Providers: []sdk.InstallationsIntegrationsRequestArgsProvidersItem{sdk.InstallationsIntegrationsRequestArgsProvidersItem(data.Provider.ValueString())}}})
	if err != nil {
		resp.Diagnostics.AddError("Unable to find integration", apiError(err))
		return
	}
	var selected *sdk.QueryIntegrationsInstallationsResponseItem
	for _, item := range items {
		if item == nil || item.ID == "" || string(item.Provider) != data.Provider.ValueString() || item.Status == "" {
			resp.Diagnostics.AddError("Invalid integration response", "The API returned an incomplete or unexpected installation.")
			return
		}
		if !data.IntegrationID.IsNull() && item.ID != data.IntegrationID.ValueString() {
			continue
		}
		if selected != nil {
			resp.Diagnostics.AddError("Ambiguous integration", "More than one installation matches. Set integration_id to select one.")
			return
		}
		selected = item
	}
	if selected == nil {
		resp.Diagnostics.AddError("Integration not found", "No matching installation was returned. Check the provider, integration_id, and workspace permissions.")
		return
	}
	data.ID = types.StringValue(selected.ID)
	data.Status = types.StringValue(selected.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

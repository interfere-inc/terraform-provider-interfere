package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/datasource_service"
	"github.com/interfere-inc/terraform-provider-interfere/internal/datasource_services"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

type serviceDataSource struct{ apiDataSource }
type servicesDataSource struct{ apiDataSource }
type servicesModel struct {
	WorkspaceSlug types.String `tfsdk:"workspace_slug"`
	Services      types.Set    `tfsdk:"services"`
}
type serviceListItem struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Namespace types.String `tfsdk:"namespace"`
	Type      types.String `tfsdk:"type"`
}

func NewServiceDataSource() datasource.DataSource  { return &serviceDataSource{} }
func NewServicesDataSource() datasource.DataSource { return &servicesDataSource{} }

func (r *serviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}
func (r *servicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

func (r *serviceDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_service.ServiceDataSourceSchema(ctx)
	resp.Schema.Description = "Read an existing service and its stored website settings by ID. Requires org:integration_settings:read."
	id := resp.Schema.Attributes["id"].(schema.StringAttribute)
	id.Required, id.Computed = true, false
	resp.Schema.Attributes["id"] = id
	settings := resp.Schema.Attributes["website_settings"].(schema.SingleNestedAttribute)
	settings.CustomType = nil
	plugins := settings.Attributes["plugins"].(schema.SingleNestedAttribute)
	plugins.CustomType = nil
	settings.Attributes["plugins"] = plugins
	resp.Schema.Attributes["website_settings"] = settings
}

func (r *serviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data serviceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := (&serviceResource{client: r.client}).read(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Unable to find service", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *servicesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_services.ServicesDataSourceSchema(ctx)
	resp.Schema.Description = "List up to 500 services in a workspace, ordered by namespace, name, and ID. Requires org:integration_settings:read. Use interfere_service to look up a specific service outside this list."
	plainSetAttribute(resp.Schema.Attributes, "services", "The first 500 services in the workspace.")
}

func (r *servicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data servicesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.client.Services.ForWorkspace(ctx, &sdk.ForWorkspaceServicesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.ForWorkspaceServicesRequestArgs{}})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list services", apiError(err))
		return
	}
	items := []serviceListItem{}
	for _, item := range remote {
		if item == nil || !uuidPattern.MatchString(item.ID) || item.Name == "" || (item.Type != nil && !validServiceType(*item.Type)) {
			resp.Diagnostics.AddError("Invalid service response", "The API returned an incomplete service.")
			return
		}
		items = append(items, serviceListItem{ID: types.StringValue(item.ID), Name: types.StringValue(item.Name), Namespace: types.StringPointerValue(item.Namespace), Type: types.StringPointerValue(item.Type)})
	}
	value, ok := objectSet(ctx, datasource_services.ServicesValue{}.AttributeTypes(ctx), items)
	if !ok {
		resp.Diagnostics.AddError("Invalid service response", "The services could not be converted to Terraform values.")
		return
	}
	data.Services = value
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/datasource_tracking_domain"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

type trackingDomainDataSource struct{ apiDataSource }
type trackingDomainLookupModel struct {
	ID            types.String `tfsdk:"id"`
	WorkspaceSlug types.String `tfsdk:"workspace_slug"`
	Name          types.String `tfsdk:"name"`
	Status        types.String `tfsdk:"status"`
	SSLStatus     types.String `tfsdk:"ssl_status"`
	CNAMETarget   types.String `tfsdk:"cname_target"`
}

func NewTrackingDomainDataSource() datasource.DataSource { return &trackingDomainDataSource{} }
func (r *trackingDomainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tracking_domain"
}

func (r *trackingDomainDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_tracking_domain.TrackingDomainDataSourceSchema(ctx)
	resp.Schema.Description = "Look up an existing workspace tracking hostname by name without managing it. Requires org:workspace_domains:read."
	attributes := resp.Schema.Attributes
	domain := attributes["tracking_domain"].(schema.SetNestedAttribute).NestedObject.Attributes
	attributes["id"], attributes["status"] = domain["id"], domain["status"]
	attributes["ssl_status"] = domain["domain_metadata"].(schema.SingleNestedAttribute).Attributes["cf_ssl_status"]
	name := domain["name"].(schema.StringAttribute)
	name.Required, name.Computed = true, false
	attributes["name"] = name
	attributes["cname_target"] = schema.StringAttribute{Computed: true, Description: "DNS-only CNAME target for the tracking hostname."}
	delete(attributes, "tracking_domain")
}

func (r *trackingDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data trackingDomainLookupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.client.Domains.List(ctx, &sdk.ListDomainsRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Unable to find tracking domain", apiError(err))
		return
	}
	var selected *sdk.QueryDomainsForOrgResponseItem
	for _, item := range remote {
		if item == nil || item.ID == "" || item.Name == "" {
			resp.Diagnostics.AddError("Invalid domain response", "The API returned an incomplete domain.")
			return
		}
		if item.DeletedAt != nil || item.Type == nil || *item.Type != "proxy" || item.Name != data.Name.ValueString() {
			continue
		}
		if item.Status == nil || selected != nil {
			resp.Diagnostics.AddError("Invalid domain response", "The API returned an incomplete or duplicate tracking domain.")
			return
		}
		selected = item
	}
	if selected == nil {
		resp.Diagnostics.AddError("Tracking domain not found", "No active tracking domain with this name was returned. Check the hostname and workspace domain permissions.")
		return
	}
	data.ID = types.StringValue(selected.ID)
	data.Status = types.StringValue(string(*selected.Status))
	data.CNAMETarget = types.StringValue(trackingCNAMETarget)
	data.SSLStatus = types.StringNull()
	if selected.DomainMetadata != nil {
		data.SSLStatus = types.StringPointerValue(selected.DomainMetadata.CfSslStatus)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

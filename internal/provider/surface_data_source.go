package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type surfaceDataSource struct{ client *client.Client }
type surfaceLookupModel struct {
	ID                    types.String `tfsdk:"id"`
	WorkspaceSlug         types.String `tfsdk:"workspace_slug"`
	Slug                  types.String `tfsdk:"slug"`
	Name                  types.String `tfsdk:"name"`
	Type                  types.String `tfsdk:"type"`
	AnonymousUserTracking types.Bool   `tfsdk:"anonymous_user_tracking"`
}

func NewSurfaceDataSource() datasource.DataSource { return &surfaceDataSource{} }
func (r *surfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_surface"
}

func (r *surfaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *surfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Look up an existing surface without managing its lifecycle.", Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{Computed: true}, "workspace_slug": schema.StringAttribute{Required: true}, "slug": schema.StringAttribute{Required: true}, "name": schema.StringAttribute{Computed: true}, "type": schema.StringAttribute{Computed: true}, "anonymous_user_tracking": schema.BoolAttribute{Computed: true},
	}}
}

func (r *surfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data surfaceLookupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := (&surfaceResource{client: r.client}).read(ctx, surfaceModel{WorkspaceSlug: data.WorkspaceSlug, Slug: data.Slug})
	if err != nil {
		resp.Diagnostics.AddError("Unable to find surface", err.Error())
		return
	}
	if remote.DeletedAt != nil {
		resp.Diagnostics.AddError("Surface was deleted", "Select an active surface.")
		return
	}
	data.ID = types.StringValue(remote.ID)
	data.Name = types.StringValue(remote.Name)
	data.Type = types.StringValue(string(remote.Type))
	data.AnonymousUserTracking = types.BoolValue(remote.AnonymousUserTracking)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

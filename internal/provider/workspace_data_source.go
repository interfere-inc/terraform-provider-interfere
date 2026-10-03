package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/datasource_workspace"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type workspaceDataSource struct{ client *client.Client }

func NewWorkspaceDataSource() datasource.DataSource { return &workspaceDataSource{} }
func (r *workspaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (r *workspaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *workspaceDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_workspace.WorkspaceDataSourceSchema(ctx)
	resp.Schema.Description = "Look up an existing workspace without managing its lifecycle."
}

func (r *workspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workspaceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := (&workspaceResource{client: r.client}).read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to find workspace", err.Error())
		return
	}
	data.Id = types.StringValue(remote.ID)
	data.Name = types.StringValue(remote.Name)
	data.DataResidencyLocation = types.StringValue(string(remote.DataResidencyLocation))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

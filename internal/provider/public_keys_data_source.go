package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/datasource_public_keys"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/keys"
)

type publicKeysDataSource struct{ apiDataSource }
type publicKeysModel struct {
	WorkspaceSlug types.String `tfsdk:"workspace_slug"`
	PublicKeys    types.Set    `tfsdk:"public_keys"`
}
type publicKeyListItem struct {
	ID      string `tfsdk:"id"`
	Name    string `tfsdk:"name"`
	Content string `tfsdk:"content"`
}

func NewPublicKeysDataSource() datasource.DataSource { return &publicKeysDataSource{} }
func (r *publicKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_keys"
}

func (r *publicKeysDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_public_keys.PublicKeysDataSourceSchema(ctx)
	resp.Schema.Description = "List active workspace publishable keys without managing them. Reads require org:workspace_auth:read or org:surfaces:read. Public key content is a browser credential, not a secret."
	plainSetAttribute(resp.Schema.Attributes, "public_keys", "Active workspace publishable keys.")
}

func (r *publicKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data publicKeysModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	items, err := r.list(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list public keys", err.Error())
		return
	}
	value, ok := objectSet(ctx, datasource_public_keys.PublicKeysValue{}.AttributeTypes(ctx), items)
	if !ok {
		resp.Diagnostics.AddError("Invalid public key response", "The public keys could not be converted to Terraform values.")
		return
	}
	data.PublicKeys = value
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *publicKeysDataSource) list(ctx context.Context, data publicKeysModel) ([]publicKeyListItem, error) {
	remote, err := r.client.Keys.Public.ForWorkspace(ctx, &keys.ForWorkspacePublicRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: map[string]any{}})
	if err != nil {
		return nil, errors.New(apiError(err))
	}
	items := []publicKeyListItem{}
	for _, item := range remote {
		if item == nil || item.ID == "" || item.Name == "" || item.Content == "" || item.SurfaceID != nil {
			return nil, errors.New("The API returned an incomplete or unexpected workspace public key.")
		}
		if item.DeletedAt == nil {
			items = append(items, publicKeyListItem{ID: item.ID, Name: item.Name, Content: item.Content})
		}
	}
	return items, nil
}

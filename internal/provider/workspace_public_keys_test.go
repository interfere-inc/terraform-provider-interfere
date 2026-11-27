package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/interfere-inc/terraform-provider-interfere/internal/datasource_public_keys"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
	"github.com/stretchr/testify/require"
)

func TestWorkspacePublicKeysDataSource(t *testing.T) {
	f := newPublicKeyFixture(t)
	f.keys[attemptID] = map[string]any{"id": attemptID, "name": "Workspace", "content": "pk_workspace", "surfaceId": nil}
	config := fmt.Sprintf(`
provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
data "interfere_public_keys" "workspace" {
 workspace_slug = "example"
}
`, f.server.URL)
	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
		Steps: []tfresource.TestStep{{Config: config, Check: tfresource.ComposeAggregateTestCheckFunc(
			tfresource.TestCheckTypeSetElemNestedAttrs("data.interfere_public_keys.workspace", "public_keys.*", map[string]string{"id": attemptID, "name": "Workspace", "content": "pk_workspace"}),
			tfresource.TestCheckResourceAttr("data.interfere_public_keys.workspace", "public_keys.#", "1"),
		)}},
	})
}

func TestWorkspacePublicKeysRejectAppKey(t *testing.T) {
	ctx := context.Background()
	f := newPublicKeyFixture(t)
	f.keys[attemptID] = map[string]any{"id": attemptID, "name": "App", "content": "pk_app", "surfaceId": workspaceID}
	r := &publicKeysDataSource{apiDataSource{client: client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())}}
	var schema datasource.SchemaResponse
	r.Schema(ctx, datasource.SchemaRequest{}, &schema)
	config := tfsdk.Config{Schema: schema.Schema}
	model := publicKeysModel{WorkspaceSlug: types.StringValue("example"), PublicKeys: types.SetNull(types.ObjectType{AttrTypes: datasource_public_keys.PublicKeysValue{}.AttributeTypes(ctx)})}
	state := tfsdk.State{Schema: schema.Schema}
	require.False(t, state.Set(ctx, model).HasError())
	config.Raw = state.Raw
	response := datasource.ReadResponse{State: state}
	r.Read(ctx, datasource.ReadRequest{Config: config}, &response)
	require.True(t, response.Diagnostics.HasError())
	require.True(t, response.State.Raw.Equal(state.Raw))
}

func TestWorkspacePublicKeyReadFailures(t *testing.T) {
	for _, scenario := range []struct {
		name    string
		status  int
		body    any
		removed bool
	}{
		{"missing", 404, map[string]string{"code": "SURFACE_KEY_NOT_FOUND"}, true},
		{"revoked", 200, map[string]any{"id": attemptID, "name": "Browser", "content": "pk_workspace", "revoked": true}, true},
		{"permission", 403, map[string]string{"message": "credential-from-error"}, false},
		{"outage", 503, map[string]string{"message": "credential-from-error"}, false},
		{"unclassified missing", 404, map[string]string{"message": "credential-from-error"}, false},
		{"null", 200, nil, false},
		{"wrong identity", 200, map[string]any{"id": "other", "name": "Browser", "content": "credential-from-error"}, false},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			ctx := context.Background()
			f := newPublicKeyFixture(t)
			f.readStatus, f.readBody = scenario.status, scenario.body
			r, state := fixturePublicKey(t, f)
			response := resource.ReadResponse{State: state}
			r.Read(ctx, resource.ReadRequest{State: state}, &response)
			if scenario.removed {
				require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
				require.True(t, response.State.Raw.IsNull())
			} else {
				require.True(t, response.Diagnostics.HasError())
				require.True(t, response.State.Raw.Equal(state.Raw))
				require.NotContains(t, fmt.Sprint(response.Diagnostics), "credential-from-error")
			}
		})
	}
}

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
	"github.com/stretchr/testify/require"
)

const serviceFixtureID = "75000000-0000-4000-8000-000000000001"

type serviceFixture struct {
	mu          sync.Mutex
	server      *httptest.Server
	serviceType string
	settings    *serviceWebsiteSettings
	readStatus  int
}

func newServiceFixture(t *testing.T) *serviceFixture {
	t.Helper()
	f := &serviceFixture{serviceType: "backend"}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer fixture-token" {
			t.Error("Unexpected service method or credentials")
			w.WriteHeader(400)
			return
		}
		var input struct {
			Args struct {
				ServiceID       string                  `json:"serviceId"`
				Type            string                  `json:"type"`
				WebsiteSettings *serviceWebsiteSettings `json:"websiteSettings"`
			} `json:"args"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Error(err)
			w.WriteHeader(400)
			return
		}
		remote := map[string]any{"id": serviceFixtureID, "organizationId": workspaceID, "name": "checkout", "namespace": nil, "type": f.serviceType, "websiteSettings": f.settings, "createdAt": 1, "updatedAt": 1, "lastSeenAt": 1}
		var result any
		switch r.URL.Path {
		case "/v3/workspaces/example/queries/services.byId":
			if input.Args.ServiceID != serviceFixtureID {
				t.Error("Wrong service identity")
				w.WriteHeader(404)
				return
			}
			if f.readStatus != 0 {
				w.WriteHeader(f.readStatus)
				_, _ = w.Write([]byte(`{"message":"fixture denial"}`))
				return
			}
			result = remote
		case "/v3/workspaces/example/queries/services.forWorkspace":
			result = []any{remote}
		case "/v3/workspaces/example/actions/services.setType":
			if input.Args.ServiceID != serviceFixtureID {
				t.Error("Wrong service mutation identity")
				w.WriteHeader(400)
				return
			}
			f.serviceType = input.Args.Type
			result = map[string]bool{"success": true}
		case "/v3/workspaces/example/actions/services.setWebsiteSettings":
			if input.Args.ServiceID != serviceFixtureID || input.Args.WebsiteSettings == nil {
				t.Error("Invalid settings mutation")
				w.WriteHeader(400)
				return
			}
			f.settings = input.Args.WebsiteSettings
			result = map[string]bool{"success": true}
		default:
			t.Errorf("Unexpected service/App creation or deletion: %s", r.URL.Path)
			w.WriteHeader(404)
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(f.server.Close)
	return f
}

func TestServiceSettingsLifecycle(t *testing.T) {
	f := newServiceFixture(t)
	config := func(settings string) string {
		return fmt.Sprintf(`
provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
resource "interfere_service" "test" {
 workspace_slug = "example"
 id = %q
 type = "website"
 %s
}
data "interfere_service" "one" {
 workspace_slug = "example"
 id = interfere_service.test.id
}
data "interfere_services" "all" { workspace_slug = "example" }
`, f.server.URL, serviceFixtureID, settings)
	}
	settings := `website_settings = {
 anonymous_replay_policy = "on_error"
 anonymous_user_tracking = false
 environment_filtering = false
 plugins = { replay = false, logs = false }
}`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
		Steps: []resource.TestStep{
			{Config: config(""), Check: resource.TestCheckResourceAttr("interfere_service.test", "type", "website")},
			{Config: config(settings), Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr("interfere_service.test", "website_settings.plugins.replay", "false"), resource.TestCheckResourceAttr("data.interfere_service.one", "website_settings.anonymous_replay_policy", "on_error"), resource.TestCheckResourceAttr("data.interfere_services.all", "services.#", "1"))},
			{ResourceName: "interfere_service.test", ImportState: true, ImportStateId: "example/" + serviceFixtureID, ImportStateVerify: true},
			{PreConfig: func() { f.mu.Lock(); f.settings.AnonymousReplayPolicy = "always"; f.mu.Unlock() }, Config: config(settings), Check: resource.TestCheckResourceAttr("interfere_service.test", "website_settings.anonymous_replay_policy", "on_error")},
			{Config: config(""), Check: resource.TestCheckResourceAttr("interfere_service.test", "website_settings.plugins.logs", "false")},
		},
		CheckDestroy: func(*terraform.State) error {
			f.mu.Lock()
			defer f.mu.Unlock()
			if f.serviceType != "website" || f.settings == nil || f.settings.AnonymousReplayPolicy != "on_error" || f.settings.Plugins["replay"] {
				return fmt.Errorf("Service settings were not retained")
			}
			return nil
		},
	})
}

func TestServiceReadFailurePreservesState(t *testing.T) {
	for _, status := range []int{401, 403, 404, 503} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			f := newServiceFixture(t)
			f.readStatus = status
			r := &serviceResource{client: client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())}
			ctx := context.Background()
			var definition frameworkresource.SchemaResponse
			r.Schema(ctx, frameworkresource.SchemaRequest{}, &definition)
			state := tfsdk.State{Schema: definition.Schema}
			require.False(t, state.Set(ctx, &serviceModel{ID: types.StringValue(serviceFixtureID), WorkspaceSlug: types.StringValue("example"), Name: types.StringValue("checkout"), Namespace: types.StringNull(), Type: types.StringValue("backend"), WebsiteSettings: types.ObjectNull(serviceSettingTypes(ctx))}).HasError())
			response := frameworkresource.ReadResponse{State: state}
			r.Read(ctx, frameworkresource.ReadRequest{State: state}, &response)
			require.True(t, response.Diagnostics.HasError())
			require.True(t, response.State.Raw.Equal(state.Raw))
		})
	}
}

func TestServiceRefusesLegacyAppStateMove(t *testing.T) {
	server, err := providerserver.NewProtocol6WithError(New("test")())()
	require.NoError(t, err)
	payload, err := json.Marshal(map[string]string{"id": "75000000-0000-4000-8000-000000000002", "workspace_slug": "example", "slug": "legacy-web", "name": "Legacy web"})
	require.NoError(t, err)
	response, err := server.MoveResourceState(context.Background(), &tfprotov6.MoveResourceStateRequest{SourceProviderAddress: "registry.terraform.io/interfere-inc/interfere", SourceTypeName: "interfere_surface", SourceSchemaVersion: 0, TargetTypeName: "interfere_service", SourceState: &tfprotov6.RawState{JSON: payload}})
	require.NoError(t, err)
	require.NotEmpty(t, response.Diagnostics)
	require.Nil(t, response.TargetState)
}

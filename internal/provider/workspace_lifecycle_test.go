package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func TestWorkspaceSettingsLifecycle(t *testing.T) {
	var mu sync.Mutex
	name, slug, updates := "Existing", "existing", 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer fixture-token" || r.Header.Get("CF-Access-Client-Id") != "proxy-client" || r.Header.Get("CF-Access-Client-Secret") != "proxy-secret" {
			t.Error("Unexpected authentication or method")
			w.WriteHeader(400)
			return
		}
		root := "/v3/workspaces/" + slug
		if r.URL.Path == root+"/queries/organizations.current" {
			if err := json.NewEncoder(w).Encode(map[string]any{"id": surfaceID, "name": name, "slug": slug, "dataResidencyLocation": "eu"}); err != nil {
				t.Error(err)
			}
			return
		}
		if r.URL.Path == root+"/actions/organizations.updateBasics" {
			var input sdk.UpdateBasicsWorkspacesRequest
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				t.Error(err)
				return
			}
			if input.Args == nil || input.Args.Name == nil || input.Args.Slug == nil {
				t.Error("Missing workspace settings")
				return
			}
			name, slug = *input.Args.Name, *input.Args.Slug
			updates++
			if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
				t.Error(err)
			}
			return
		}
		t.Errorf("Workspace creation/deletion must never be called: %s", r.URL.Path)
		w.WriteHeader(404)
	}))
	defer server.Close()
	config := func(workspaceSlug, workspaceName string) string {
		return fmt.Sprintf(`
provider "interfere" {
 headers = { "CF-Access-Client-Id" = "proxy-client", "CF-Access-Client-Secret" = "proxy-secret" }
 token = "fixture-token"
 base_url = %q
}
resource "interfere_workspace" "test" {
 workspace_slug = %q
 name = %q
}
`, server.URL, workspaceSlug, workspaceName)
	}
	check := func(expected string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("interfere_workspace.test", "id", surfaceID),
			resource.TestCheckResourceAttr("interfere_workspace.test", "name", expected),
			resource.TestCheckResourceAttr("interfere_workspace.test", "data_residency_location", "eu"),
		)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
		Steps: []resource.TestStep{
			{Config: config("existing", "Managed"), Check: check("Managed")},
			{ResourceName: "interfere_workspace.test", ImportState: true, ImportStateId: "existing", ImportStateVerify: true},
			{Config: config("existing", "Renamed"), Check: check("Renamed")},
			{PreConfig: func() { mu.Lock(); name = "Drift"; mu.Unlock() }, Config: config("existing", "Renamed"), Check: check("Renamed")},
			{Config: config("new-slug", "Renamed"), Check: resource.TestCheckResourceAttr("interfere_workspace.test", "workspace_slug", "new-slug")},
		},
		CheckDestroy: func(*terraform.State) error {
			mu.Lock()
			defer mu.Unlock()
			if name != "Renamed" || slug != "new-slug" || updates != 4 {
				return fmt.Errorf("Unexpected retained workspace: %s/%s, updates=%d", slug, name, updates)
			}
			return nil
		},
	})
	if strings.TrimSpace(name) == "" {
		t.Fatal("Workspace was lost")
	}
}

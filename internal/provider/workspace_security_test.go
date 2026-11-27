package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type workspaceSecurityFixture struct {
	mu        sync.Mutex
	creation  bool
	granted   bool
	expiresAt *float64
	methods   []string
	writes    map[string]int
}

func (f *workspaceSecurityFixture) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		root := "/v3/workspaces/existing"
		var input struct {
			Args map[string]json.RawMessage `json:"args"`
		}
		if r.URL.Path != root+"/queries/workspaces.current" {
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				t.Error(err)
				return
			}
		}
		switch r.URL.Path {
		case root + "/queries/workspaces.current":
			if err := json.NewEncoder(w).Encode(map[string]any{
				"id": workspaceID, "name": "Existing", "slug": "existing", "dataResidencyLocation": "eu",
				"allowMemberWorkspaceCreation": f.creation, "impersonationGranted": f.granted,
				"impersonationExpiresAt": f.expiresAt, "enabledAuthMethods": f.methods,
			}); err != nil {
				t.Error(err)
			}
			return
		case root + "/actions/workspaces.updateMemberWorkspaceCreation":
			if err := json.Unmarshal(input.Args["allowMemberWorkspaceCreation"], &f.creation); err != nil {
				t.Error(err)
			}
		case root + "/actions/workspaces.updateSupportAccess":
			expiry, present := input.Args["impersonationExpiresAt"]
			if !present {
				t.Error("Support access updates must send an explicit expiry")
			}
			f.expiresAt = nil
			if err := json.Unmarshal(expiry, &f.expiresAt); err != nil {
				t.Error(err)
			}
			if err := json.Unmarshal(input.Args["impersonationGranted"], &f.granted); err != nil {
				t.Error(err)
			}
		case root + "/actions/workspaces.updateAuthMethods":
			if err := json.Unmarshal(input.Args["methods"], &f.methods); err != nil {
				t.Error(err)
			}
		default:
			t.Errorf("Unexpected workspace call: %s", r.URL.Path)
			w.WriteHeader(404)
			return
		}
		f.writes[r.URL.Path[len(root)+1:]]++
		if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
			t.Error(err)
		}
	}
}

func TestWorkspaceSecuritySettingsLifecycle(t *testing.T) {
	expiry := float64(1893456000000)
	fixture := &workspaceSecurityFixture{granted: true, expiresAt: &expiry, writes: map[string]int{}}
	server := httptest.NewServer(fixture.handler(t))
	defer server.Close()
	config := func(settings string) string {
		return fmt.Sprintf(`
provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
resource "interfere_workspace" "test" {
 workspace_slug = "existing"
 name = "Existing"
%s
}
`, server.URL, settings)
	}
	const name = "interfere_workspace.test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
		Steps: []resource.TestStep{
			{Config: config(`
 enabled_auth_methods = ["oauth_google"]
`), ExpectError: regexp.MustCompile("Missing sign-in method")},
			{Config: config(""), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "allow_member_workspace_creation", "false"),
				resource.TestCheckResourceAttr(name, "impersonation_granted", "true"),
				resource.TestCheckResourceAttr(name, "impersonation_expires_at", "1893456000000"),
				resource.TestCheckNoResourceAttr(name, "enabled_auth_methods"),
			)},
			{Config: config(`
 allow_member_workspace_creation = true
 impersonation_granted = false
 enabled_auth_methods = ["passkey", "email_code", "oauth_github"]
`), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "allow_member_workspace_creation", "true"),
				resource.TestCheckResourceAttr(name, "impersonation_granted", "false"),
				resource.TestCheckNoResourceAttr(name, "impersonation_expires_at"),
				resource.TestCheckResourceAttr(name, "enabled_auth_methods.#", "3"),
				resource.TestCheckTypeSetElemAttr(name, "enabled_auth_methods.*", "oauth_github"),
			)},
			{ResourceName: name, ImportState: true, ImportStateId: "existing", ImportStateVerify: true},
			{Config: config(`
 impersonation_granted = true
 impersonation_expires_at = 1893456000000
`), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "allow_member_workspace_creation", "true"),
				resource.TestCheckResourceAttr(name, "impersonation_expires_at", "1893456000000"),
				resource.TestCheckResourceAttr(name, "enabled_auth_methods.#", "3"),
			)},
		},
		CheckDestroy: func(*terraform.State) error {
			fixture.mu.Lock()
			defer fixture.mu.Unlock()
			if !fixture.creation || !fixture.granted || fixture.expiresAt == nil || len(fixture.methods) != 3 {
				return fmt.Errorf("Destroy must keep workspace settings: %+v", fixture)
			}
			if fixture.writes["actions/workspaces.updateMemberWorkspaceCreation"] != 1 || fixture.writes["actions/workspaces.updateSupportAccess"] != 2 || fixture.writes["actions/workspaces.updateAuthMethods"] != 1 {
				return fmt.Errorf("Unexpected settings writes: %v", fixture.writes)
			}
			if fixture.methods[0] != "oauth_github" {
				return fmt.Errorf("Sign-in methods must be sent in priority order: %v", fixture.methods)
			}
			return nil
		},
	})
}

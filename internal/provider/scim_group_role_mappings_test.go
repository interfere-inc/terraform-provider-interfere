package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type scimGroupFixture struct {
	mu     sync.Mutex
	groups []map[string]any
	writes map[string]int
}

func (f *scimGroupFixture) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		root := "/v3/workspaces/existing"
		var input struct {
			Args map[string]json.RawMessage `json:"args"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Error(err)
			return
		}
		switch r.URL.Path {
		case root + "/queries/scimGroups.forOrg":
			if err := json.NewEncoder(w).Encode(f.groups); err != nil {
				t.Error(err)
			}
			return
		case root + "/actions/workspaces.setScimGroupRole":
			var id string
			var role *string
			if err := json.Unmarshal(input.Args["id"], &id); err != nil {
				t.Error(err)
			}
			raw, present := input.Args["roleKey"]
			if !present {
				t.Error("Role updates must send an explicit roleKey")
			}
			if err := json.Unmarshal(raw, &role); err != nil {
				t.Error(err)
			}
			highest := 0
			for _, group := range f.groups {
				if precedence, ok := group["precedence"].(int); ok && precedence > highest {
					highest = precedence
				}
			}
			for _, group := range f.groups {
				if group["id"] != id {
					continue
				}
				if role == nil {
					group["roleKey"], group["precedence"] = nil, nil
					continue
				}
				group["roleKey"] = *role
				if group["precedence"] == nil {
					group["precedence"] = highest + 1
				}
			}
		case root + "/actions/workspaces.reorderScimGroupRoles":
			var ids []string
			if err := json.Unmarshal(input.Args["ids"], &ids); err != nil {
				t.Error(err)
			}
			for index, id := range ids {
				for _, group := range f.groups {
					if group["id"] == id {
						group["precedence"] = index + 1
					}
				}
			}
		default:
			t.Errorf("Unexpected SCIM call: %s", r.URL.Path)
			w.WriteHeader(404)
			return
		}
		f.writes[r.URL.Path[len(root)+1:]]++
		if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
			t.Error(err)
		}
	}
}

func scimGroup(index int, name string) map[string]any {
	return map[string]any{
		"id": fmt.Sprintf("00000000-0000-4000-8000-00000000000%d", index), "clerkScimDirectoryId": "scim_dir_1", "scimGroupId": name,
		"scimGroupDisplayName": strings.ToUpper(name), "roleKey": nil, "precedence": nil, "memberCount": 0,
		"organizationId": workspaceID, "createdAt": 1, "updatedAt": 1, "deletedAt": nil,
	}
}

func TestSCIMGroupRoleMappingsLifecycle(t *testing.T) {
	fixture := &scimGroupFixture{groups: []map[string]any{scimGroup(1, "admins"), scimGroup(2, "engineers"), scimGroup(3, "contractors")}, writes: map[string]int{}}
	server := httptest.NewServer(fixture.handler(t))
	defer server.Close()
	config := func(mappings string) string {
		return fmt.Sprintf(`
provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
resource "interfere_scim_group_role_mappings" "test" {
 workspace_slug = "existing"
 directory_id = "scim_dir_1"
 mappings = [%s]
}
`, server.URL, mappings)
	}
	const name = "interfere_scim_group_role_mappings.test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
		Steps: []resource.TestStep{
			{Config: config(`{ scim_group_id = "missing", role = "org:member" }`), ExpectError: regexp.MustCompile("has not been provisioned")},
			{Config: config(`
 { scim_group_id = "admins", role = "org:admin" },
 { scim_group_id = "engineers", role = "org:member" },
`), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "id", "scim_dir_1"),
				resource.TestCheckResourceAttr(name, "mappings.#", "2"),
				resource.TestCheckResourceAttr(name, "mappings.0.scim_group_id", "admins"),
				resource.TestCheckResourceAttr(name, "mappings.1.role", "org:member"),
			)},
			{Config: config(`
 { scim_group_id = "contractors", role = "org:member" },
 { scim_group_id = "admins", role = "org:admin" },
`), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "mappings.#", "2"),
				resource.TestCheckResourceAttr(name, "mappings.0.scim_group_id", "contractors"),
				resource.TestCheckResourceAttr(name, "mappings.1.scim_group_id", "admins"),
			)},
			{ResourceName: name, ImportState: true, ImportStateId: "existing/scim_dir_1", ImportStateVerify: true},
		},
		CheckDestroy: func(*terraform.State) error {
			fixture.mu.Lock()
			defer fixture.mu.Unlock()
			for _, group := range fixture.groups {
				if group["roleKey"] != nil {
					return fmt.Errorf("Destroy must unmap every group: %v", group)
				}
			}
			if fixture.writes["actions/workspaces.reorderScimGroupRoles"] != 1 {
				return fmt.Errorf("Unexpected reorder writes: %v", fixture.writes)
			}
			return nil
		},
	})
}

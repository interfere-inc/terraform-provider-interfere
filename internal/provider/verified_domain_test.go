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

type verifiedDomainFixture struct {
	mu      sync.Mutex
	domains []map[string]any
	writes  map[string]int
}

func (f *verifiedDomainFixture) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		root := "/v3/workspaces/existing"
		var input struct {
			Args map[string]string `json:"args"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Error(err)
			return
		}
		find := func(clerkID string) map[string]any {
			for _, domain := range f.domains {
				if domain["clerkDomainId"] == clerkID {
					return domain
				}
			}
			t.Errorf("Unknown domain: %s", clerkID)
			return map[string]any{}
		}
		switch r.URL.Path {
		case root + "/queries/domains.forOrg":
			if err := json.NewEncoder(w).Encode(f.domains); err != nil {
				t.Error(err)
			}
			return
		case root + "/actions/workspaces.addDomain":
			f.domains = append(f.domains, map[string]any{
				"id": input.Args["id"], "name": input.Args["name"], "type": "clerk", "clerkDomainId": "dmn_1",
				"enrollmentMode": "manual_invitation", "verifiedAt": nil, "affiliationEmailAddress": nil, "deletedAt": nil,
				"organizationId": workspaceID, "domainMetadata": nil, "status": nil, "createdAt": 1, "updatedAt": 1,
			})
		case root + "/actions/workspaces.updateDomainEnrollmentMode":
			find(input.Args["clerkDomainId"])["enrollmentMode"] = input.Args["enrollmentMode"]
		case root + "/actions/workspaces.removeDomain":
			removed := find(input.Args["clerkDomainId"])
			var kept []map[string]any
			for _, domain := range f.domains {
				if fmt.Sprint(domain["id"]) != fmt.Sprint(removed["id"]) {
					kept = append(kept, domain)
				}
			}
			f.domains = kept
		default:
			t.Errorf("Unexpected domain call: %s", r.URL.Path)
			w.WriteHeader(404)
			return
		}
		f.writes[r.URL.Path[len(root)+1:]]++
		if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
			t.Error(err)
		}
	}
}

func (f *verifiedDomainFixture) verify() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, domain := range f.domains {
		domain["verifiedAt"] = 1893456000000
		domain["affiliationEmailAddress"] = "admin@example.com"
	}
}

func TestVerifiedDomainLifecycle(t *testing.T) {
	fixture := &verifiedDomainFixture{domains: []map[string]any{}, writes: map[string]int{}}
	server := httptest.NewServer(fixture.handler(t))
	defer server.Close()
	config := func(mode string) string {
		return fmt.Sprintf(`
provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
resource "interfere_verified_domain" "test" {
 workspace_slug = "existing"
 name = "example.com"
%s
}
`, server.URL, mode)
	}
	const name = "interfere_verified_domain.test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
		Steps: []resource.TestStep{
			{Config: config(` enrollment_mode = "automatic_invitation"`), ExpectError: regexp.MustCompile("A new domain starts unverified")},
			{Config: config(""), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "enrollment_mode", "manual_invitation"),
				resource.TestCheckResourceAttr(name, "verification_status", "pending"),
				resource.TestCheckNoResourceAttr(name, "verified_at"),
				resource.TestCheckResourceAttr(name, "clerk_domain_id", "dmn_1"),
			)},
			{Config: config(` enrollment_mode = "automatic_suggestion"`), ExpectError: regexp.MustCompile("Verify the domain in the dashboard")},
			{PreConfig: fixture.verify, Config: config(` enrollment_mode = "automatic_suggestion"`), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "enrollment_mode", "automatic_suggestion"),
				resource.TestCheckResourceAttr(name, "verification_status", "verified"),
				resource.TestCheckResourceAttr(name, "verified_at", "1893456000000"),
				resource.TestCheckResourceAttr(name, "affiliation_email_address", "admin@example.com"),
			)},
			{ResourceName: name, ImportState: true, ImportStateId: "existing/example.com", ImportStateVerify: true},
		},
		CheckDestroy: func(*terraform.State) error {
			fixture.mu.Lock()
			defer fixture.mu.Unlock()
			if len(fixture.domains) != 0 {
				return fmt.Errorf("Destroy must remove the domain: %v", fixture.domains)
			}
			if fixture.writes["actions/workspaces.addDomain"] != 1 || fixture.writes["actions/workspaces.updateDomainEnrollmentMode"] != 1 || fixture.writes["actions/workspaces.removeDomain"] != 1 {
				return fmt.Errorf("Unexpected domain writes: %v", fixture.writes)
			}
			return nil
		},
	})
}

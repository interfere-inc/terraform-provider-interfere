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

type networkRuleFixture struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type networkRulesServer struct {
	mu      sync.Mutex
	rules   []networkRuleFixture
	saves   int
	lockout bool
}

func newNetworkRulesServer(t *testing.T, state *networkRulesServer) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state.mu.Lock()
		defer state.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/v3/workspaces/existing/network/rules" || r.Header.Get("Authorization") != "Bearer fixture-token" {
			t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
			return
		}
		switch r.Method {
		case http.MethodGet:
			if err := json.NewEncoder(w).Encode(map[string]any{"rules": state.rules}); err != nil {
				t.Error(err)
			}
		case http.MethodPut:
			var input struct {
				Rules []networkRuleFixture `json:"rules"`
			}
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Rules == nil {
				t.Error("Network rules must be sent as a list")
				w.WriteHeader(400)
				return
			}
			if state.lockout && len(input.Rules) > 0 {
				w.WriteHeader(409)
				if err := json.NewEncoder(w).Encode(map[string]any{"defined": true, "code": "NETWORK_POLICY_LOCKOUT", "status": 409, "message": "Locked out"}); err != nil {
					t.Error(err)
				}
				return
			}
			state.rules = input.Rules
			state.saves++
			if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
				t.Error(err)
			}
		default:
			t.Errorf("Unexpected method: %s", r.Method)
			w.WriteHeader(405)
		}
	}))
}

func networkRulesConfig(url string, rules ...[2]string) string {
	entries := make([]string, 0, len(rules))
	for _, rule := range rules {
		entries = append(entries, fmt.Sprintf("{ label = %q, value = %q }", rule[0], rule[1]))
	}
	return fmt.Sprintf(`
provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
resource "interfere_network_rules" "test" {
 workspace_slug = "existing"
 rules = [%s]
}
`, url, strings.Join(entries, ", "))
}

func networkRulesFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())}
}

func TestNetworkRulesLifecycle(t *testing.T) {
	state := &networkRulesServer{rules: []networkRuleFixture{{ID: "11111111-1111-4111-8111-111111111111", Label: "Dashboard", Value: "192.0.2.0/24"}}}
	server := newNetworkRulesServer(t, state)
	defer server.Close()
	office, runner, vpn := [2]string{"Office", "203.0.113.0/24"}, [2]string{"Runner", "198.51.100.7/32"}, [2]string{"VPN", "AS13335"}
	var officeID string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: networkRulesFactories(),
		Steps: []resource.TestStep{
			{
				Config: networkRulesConfig(server.URL, office, runner),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("interfere_network_rules.test", "rules.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("interfere_network_rules.test", "rules.*", map[string]string{"label": "Office", "value": "203.0.113.0/24"}),
					func(*terraform.State) error {
						state.mu.Lock()
						defer state.mu.Unlock()
						for _, rule := range state.rules {
							if rule.Value == "192.0.2.0/24" {
								return fmt.Errorf("Rules outside the configuration must be replaced")
							}
							if rule.Value == office[1] {
								officeID = rule.ID
							}
						}
						return nil
					},
				),
			},
			{ResourceName: "interfere_network_rules.test", ImportState: true, ImportStateId: "existing", ImportStateVerify: true, ImportStateVerifyIdentifierAttribute: "workspace_slug"},
			{
				PreConfig: func() {
					state.mu.Lock()
					state.rules = append(state.rules, networkRuleFixture{ID: "22222222-2222-4222-8222-222222222222", Label: "Manual", Value: "192.0.2.0/24"})
					state.mu.Unlock()
				},
				Config: networkRulesConfig(server.URL, office, runner),
				Check:  resource.TestCheckResourceAttr("interfere_network_rules.test", "rules.#", "2"),
			},
			{
				Config: networkRulesConfig(server.URL, office, vpn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("interfere_network_rules.test", "rules.*", map[string]string{"label": "VPN", "value": "AS13335"}),
					func(*terraform.State) error {
						state.mu.Lock()
						defer state.mu.Unlock()
						for _, rule := range state.rules {
							if rule.Value == office[1] && rule.ID != officeID {
								return fmt.Errorf("Unchanged rules must keep their identifiers")
							}
						}
						return nil
					},
				),
			},
		},
		CheckDestroy: func(*terraform.State) error {
			state.mu.Lock()
			defer state.mu.Unlock()
			if len(state.rules) != 0 {
				return fmt.Errorf("Destroy must save an empty rules list, found %d rules", len(state.rules))
			}
			return nil
		},
	})
}

func TestNetworkRulesRejectInvalidAndLockout(t *testing.T) {
	state := &networkRulesServer{rules: []networkRuleFixture{}, lockout: true}
	server := newNetworkRulesServer(t, state)
	defer server.Close()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: networkRulesFactories(),
		Steps: []resource.TestStep{
			{Config: networkRulesConfig(server.URL, [2]string{"Office", "203.0.113.7/24"}), ExpectError: regexp.MustCompile(`203\.0\.113\.0/24`)},
			{Config: networkRulesConfig(server.URL, [2]string{"Runner", "198.51.100.7"}), ExpectError: regexp.MustCompile(`198\.51\.100\.7/32`)},
			{Config: networkRulesConfig(server.URL, [2]string{"VPN", "as13335"}), ExpectError: regexp.MustCompile(`Invalid network rule`)},
			{Config: networkRulesConfig(server.URL, [2]string{"A", "AS13335"}, [2]string{"B", "AS13335"}), ExpectError: regexp.MustCompile(`Duplicate network rule`)},
			{Config: networkRulesConfig(server.URL, [2]string{" Office", "203.0.113.0/24"}), ExpectError: regexp.MustCompile(`whitespace`)},
			{Config: networkRulesConfig(server.URL, [2]string{"Office", "203.0.113.0/24"}), ExpectError: regexp.MustCompile(`Unable to save network rules`)},
		},
	})
	if state.saves != 0 {
		t.Fatalf("Rejected rules must not be saved, saved %d times", state.saves)
	}
}

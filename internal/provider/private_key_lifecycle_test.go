package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestPrivateKeyLifecycle(t *testing.T) {
	for _, surface := range []bool{false, true} {
		t.Run(fmt.Sprintf("surface=%v", surface), func(t *testing.T) {
			f := newKeyFixture(t)
			const nextID = "0f92eedc-cad4-46c4-8d2e-5f1bc4d61fe1"
			config := func(id, grants, expiry string) string {
				surfaceConfig := ""
				if surface {
					surfaceConfig = `surface_slug = "example-surface"`
				}
				return fmt.Sprintf(`
provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
resource "interfere_private_key" "test" {
 workspace_slug = "example"
 %s
 idempotency_key = %q
 name = "Managed key"
 scopes = %s
 seconds_until_expiration = %s
}
`, f.server.URL, surfaceConfig, id, grants, expiry)
			}
			check := func(id string) resource.TestCheckFunc {
				return resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("interfere_private_key.test", "id", id),
					resource.TestCheckResourceAttr("interfere_private_key.test", "secret", "fixture-secret-"+id),
					resource.TestCheckResourceAttr("interfere_private_key.test", "version", "ak_"+id),
				)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
				Steps: []resource.TestStep{
					{Config: config(attemptID, `["org:surfaces:write", "org:surfaces:read"]`, "7200"), Check: check(attemptID)},
					{Config: config(attemptID, `["org:surfaces:read"]`, "7200"), PlanOnly: true, ExpectError: regexp.MustCompile("New key creation identity required")},
					{ResourceName: "interfere_private_key.test", ImportState: true, ImportStateId: "example/" + attemptID, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"secret"}},
					{Config: config(nextID, `["org:surfaces:read"]`, "null"), Check: check(nextID)},
				},
				CheckDestroy: func(*terraform.State) error {
					f.mu.Lock()
					defer f.mu.Unlock()
					if f.creates != 2 || f.deletes != 2 {
						return fmt.Errorf("Unexpected lifecycle counts: create=%d delete=%d", f.creates, f.deletes)
					}
					for _, key := range f.keys {
						if key["revoked"] != true {
							return fmt.Errorf("Private key was not revoked")
						}
					}
					return nil
				},
			})
		})
	}
}

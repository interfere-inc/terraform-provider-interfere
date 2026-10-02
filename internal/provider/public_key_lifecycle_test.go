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

func TestPublicKeyLifecycle(t *testing.T) {
	f := newPublicKeyFixture(t)
	config := func(id, name string) string {
		return fmt.Sprintf(`
provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
resource "interfere_public_key" "test" {
 workspace_slug = "example"
 surface_slug = "example-surface"
 idempotency_key = %q
 name = %q
}
`, f.server.URL, id, name)
	}
	const nextID = "0f92eedc-cad4-46c4-8d2e-5f1bc4d61fe1"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
		Steps: []resource.TestStep{
			{Config: config(attemptID, "Browser"), Check: resource.TestCheckResourceAttr("interfere_public_key.test", "content", "interfere_public_us_"+attemptID)},
			{ResourceName: "interfere_public_key.test", ImportState: true, ImportStateId: "example/example-surface/" + attemptID, ImportStateVerify: true},
			{PreConfig: func() { f.mu.Lock(); defer f.mu.Unlock(); f.keys[attemptID]["content"] = "interfere_public_us_rotated" }, Config: config(attemptID, "Browser"), Check: resource.TestCheckResourceAttr("interfere_public_key.test", "content", "interfere_public_us_rotated")},
			{Config: config(attemptID, "Renamed"), PlanOnly: true, ExpectError: regexp.MustCompile("New key creation identity required")},
			{Config: config(nextID, "Replacement"), Check: resource.TestCheckResourceAttr("interfere_public_key.test", "id", nextID)},
		},
		CheckDestroy: func(*terraform.State) error {
			f.mu.Lock()
			defer f.mu.Unlock()
			if f.creates != 2 || f.deletes != 2 {
				return fmt.Errorf("Unexpected lifecycle counts: create=%d delete=%d", f.creates, f.deletes)
			}
			for _, key := range f.keys {
				if key["revoked"] != true {
					return fmt.Errorf("Public key was not revoked")
				}
			}
			return nil
		},
	})
}

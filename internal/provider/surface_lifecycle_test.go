package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestSurfaceLifecycle(t *testing.T) {
	f := newSurfaceFixture(t)
	config := func(name string) string {
		return fmt.Sprintf(`
provider "interfere" {
  token = "fixture-token"
  base_url = %q
}
resource "interfere_surface" "test" {
  workspace_slug = "example"
  idempotency_key = %q
  name = %q
  type = "react"
}
`, f.server.URL, attemptID, name)
	}
	check := func(name string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("interfere_surface.test", "id", surfaceID),
			resource.TestCheckResourceAttr("interfere_surface.test", "slug", surfaceSlug),
			resource.TestCheckResourceAttr("interfere_surface.test", "name", name),
			func(state *terraform.State) error {
				for _, value := range state.RootModule().Resources["interfere_surface.test"].Primary.Attributes {
					if strings.Contains(value, "must-not-enter-state") {
						return fmt.Errorf("Creation credentials entered Terraform state")
					}
				}
				return nil
			},
		)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"interfere": providerserver.NewProtocol6WithError(New("test")()),
		},
		Steps: []resource.TestStep{
			{Config: config("Example"), Check: check("Example")},
			{ResourceName: "interfere_surface.test", ImportState: true, ImportStateId: "example/" + surfaceSlug, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"idempotency_key"}},
			{Config: config("Renamed"), Check: check("Renamed")},
			{PreConfig: func() { f.mu.Lock(); f.name = "Drift"; f.mu.Unlock() }, Config: config("Renamed"), Check: check("Renamed")},
		},
		CheckDestroy: func(_ *terraform.State) error {
			f.mu.Lock()
			defer f.mu.Unlock()
			if !f.deleted || f.creates != 1 || f.updates != 2 || f.deletes != 1 {
				return fmt.Errorf("Unexpected lifecycle: deleted=%v create=%d update=%d delete=%d", f.deleted, f.creates, f.updates, f.deletes)
			}
			return nil
		},
	})
}

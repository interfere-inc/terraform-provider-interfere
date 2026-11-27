package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestWorkspaceAndIntegrationDataSources(t *testing.T) {
	f := newConfigFixture(t)
	config := fmt.Sprintf(`
provider "interfere" {
token="fixture-token"
base_url=%q
}
data "interfere_workspace" "test" {
workspace_slug="example"
}
data "interfere_integration" "test" {
workspace_slug=data.interfere_workspace.test.workspace_slug
integration_provider="github"
}
`, f.server.URL)
	factories := map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())}
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{
		{Config: config, Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr("data.interfere_workspace.test", "name", "Example"), resource.TestCheckResourceAttr("data.interfere_integration.test", "status", "connected"))},
	}})
}

func TestTrackingDomainLifecycle(t *testing.T) {
	f := newConfigFixture(t)
	config := fmt.Sprintf(`
provider "interfere" {
token="fixture-token"
base_url=%q
}
resource "interfere_tracking_domain" "test" {
workspace_slug="example"
name="track.example.com"
idempotency_key=%q
}
`, f.server.URL, attemptID)
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())}, Steps: []resource.TestStep{
		{Config: config, Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr("interfere_tracking_domain.test", "status", "pending"), resource.TestCheckResourceAttr("interfere_tracking_domain.test", "cname_target", "proxy.interfere.domains"))},
		{ResourceName: "interfere_tracking_domain.test", ImportState: true, ImportStateId: "example/" + attemptID, ImportStateVerify: true},
		{PreConfig: func() { f.mu.Lock(); f.domainStatus = "active"; f.mu.Unlock() }, RefreshState: true, Check: resource.TestCheckResourceAttr("interfere_tracking_domain.test", "status", "active")},
	}, CheckDestroy: func(_ *terraform.State) error {
		if !f.domainDeleted {
			return fmt.Errorf("Domain was not deleted")
		}
		return nil
	}})
}

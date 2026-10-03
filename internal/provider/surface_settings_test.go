package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSurfaceSettings(t *testing.T) {
	f := newSurfaceFixture(t)
	config := func(enabled bool) string {
		return fmt.Sprintf(`
provider "interfere" {
token="fixture-token"
base_url=%q
}
resource "interfere_surface" "test" {
workspace_slug="example"
idempotency_key=%q
name="Example"
type="react"
anonymous_user_tracking=%t
}
data "interfere_surface" "test" {
workspace_slug=interfere_surface.test.workspace_slug
slug=interfere_surface.test.slug
}
`, f.server.URL, attemptID, enabled)
	}
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())}, Steps: []resource.TestStep{
		{Config: config(false), Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr("interfere_surface.test", "anonymous_user_tracking", "false"), resource.TestCheckResourceAttr("data.interfere_surface.test", "anonymous_user_tracking", "false"))},
		{Config: config(true), Check: resource.TestCheckResourceAttr("interfere_surface.test", "anonymous_user_tracking", "true")},
		{PreConfig: func() { f.mu.Lock(); f.tracking = false; f.mu.Unlock() }, Config: config(true), Check: resource.TestCheckResourceAttr("interfere_surface.test", "anonymous_user_tracking", "true")},
	}})
}

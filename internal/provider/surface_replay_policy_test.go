package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSurfaceAnonymousReplayPolicy(t *testing.T) {
	f := newSurfaceFixture(t)
	config := func(policy string) string {
		setting := ""
		if policy != "" {
			setting = fmt.Sprintf("anonymous_replay_policy=%q", policy)
		}
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
%s
}
data "interfere_surface" "test" {
workspace_slug=interfere_surface.test.workspace_slug
slug=interfere_surface.test.slug
}
`, f.server.URL, attemptID, setting)
	}
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())}, Steps: []resource.TestStep{
		{Config: config("never"), Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr("interfere_surface.test", "anonymous_replay_policy", "never"), resource.TestCheckResourceAttr("data.interfere_surface.test", "anonymous_replay_policy", "never"))},
		{Config: config("on_error"), Check: resource.TestCheckResourceAttr("interfere_surface.test", "anonymous_replay_policy", "on_error")},
		{PreConfig: func() { f.mu.Lock(); f.replayPolicy = "always"; f.mu.Unlock() }, Config: config("on_error"), Check: resource.TestCheckResourceAttr("interfere_surface.test", "anonymous_replay_policy", "on_error")},
		{Config: config(""), Check: resource.TestCheckResourceAttr("interfere_surface.test", "anonymous_replay_policy", "on_error")},
		{Config: config("sometimes"), ExpectError: regexp.MustCompile(`value must be one of`)},
	}})
}

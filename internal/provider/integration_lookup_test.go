package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIntegrationLookupSelection(t *testing.T) {
	for _, mode := range []string{"missing", "ambiguous", "selected", "unexpected"} {
		t.Run(mode, func(t *testing.T) {
			f := newConfigFixture(t)
			if mode == "missing" {
				f.integrations = nil
			}
			if mode == "ambiguous" || mode == "selected" {
				f.integrations = append(f.integrations, map[string]any{"id": "integration-2", "provider": "github", "status": "connected"})
			}
			if mode == "unexpected" {
				f.integrations[0]["provider"] = "slack"
			}
			selector := ""
			if mode == "selected" {
				selector = "integration_id=\"integration-2\""
			}
			config := fmt.Sprintf(`
provider "interfere" {
 token="fixture-token"
 base_url=%q
}
data "interfere_integration" "test" {
 workspace_slug="example"
 integration_provider="github"
 %s
}
`, f.server.URL, selector)
			step := resource.TestStep{Config: config}
			failures := map[string]string{"missing": "Integration not found", "ambiguous": "Ambiguous integration", "unexpected": "Invalid integration response"}
			if expected, ok := failures[mode]; ok {
				step.ExpectError = regexp.MustCompile(expected)
			} else {
				step.Check = resource.TestCheckResourceAttr("data.interfere_integration.test", "id", "integration-2")
			}
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())}, Steps: []resource.TestStep{step}})
		})
	}
}

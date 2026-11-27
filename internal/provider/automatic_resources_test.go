package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAutomaticTrackingDomainIdentity(t *testing.T) {
	f := newConfigFixture(t)
	config := func(name string) string {
		return fmt.Sprintf(`provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
resource "interfere_tracking_domain" "test" {
 workspace_slug = "example"
 name = %q
}`, f.server.URL, name)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
		Steps: []resource.TestStep{
			{Config: config("track.example.com")},
			{Config: config("track.example.com"), PlanOnly: true},
			{Config: config("events.example.com")},
			{Config: config("events.example.com"), PlanOnly: true},
		},
	})
	require.Len(t, f.domainIdentities, 2)
	require.NotEqual(t, f.domainIdentities[0], f.domainIdentities[1])
	require.True(t, f.domainDeleted)
}

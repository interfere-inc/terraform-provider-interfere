package provider

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func TestAutomaticKeyLifecycle(t *testing.T) {
	for _, kind := range []string{"private", "public"} {
		t.Run(kind, func(t *testing.T) {
			private, public := newKeyFixture(t), newPublicKeyFixture(t)
			private.createFailures, public.createFailures = 1, 1
			url, extra := private.server.URL, "scopes = [\"org:surfaces:read\"]"
			if kind == "public" {
				url, extra = public.server.URL, "surface_slug = \"example-surface\""
			}
			address := "interfere_" + kind + "_key.test"
			config := func(name, identity string) string {
				return fmt.Sprintf(`provider "interfere" {
 token = "fixture-token"
 base_url = %q
}
resource "interfere_%s_key" "test" {
 workspace_slug = "example"
 name = %q
 %s
 %s
 lifecycle { create_before_destroy = true }
}`, url, kind, name, extra, identity)
			}
			var original, replacement string
			capture := func(target *string) resource.TestCheckFunc {
				return func(state *terraform.State) error {
					attributes := state.RootModule().Resources[address].Primary.Attributes
					*target = attributes["id"]
					if _, err := uuid.Parse(*target); err != nil {
						return err
					}
					if attributes["idempotency_key"] != *target {
						return fmt.Errorf("Creation identity differs from the key ID")
					}
					return nil
				}
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())},
				Steps: []resource.TestStep{
					{Config: config("Original", fmt.Sprintf("idempotency_key = %q", attemptID)), Check: capture(&original)},
					{Config: config("Original", ""), PlanOnly: true},
					{Config: config("Replacement", ""), Check: capture(&replacement)},
					{Config: config("Replacement", ""), PlanOnly: true},
					{Config: config("Replacement", ""), Taint: []string{address}},
				},
			})
			require.Equal(t, attemptID, original)
			require.NotEqual(t, original, replacement)
			keys, identities, creates, deletes := private.keys, private.createIdentities, private.creates, private.deletes
			if kind == "public" {
				keys, identities, creates, deletes = public.keys, public.createIdentities, public.creates, public.deletes
			}
			require.Equal(t, 3, creates)
			require.Equal(t, 3, deletes)
			require.Len(t, identities, 4)
			require.Equal(t, identities[0], identities[1])
			require.NotEqual(t, identities[1], identities[2])
			require.NotEqual(t, identities[2], identities[3])
			for _, key := range keys {
				require.Equal(t, true, key["revoked"])
			}
		})
	}
}

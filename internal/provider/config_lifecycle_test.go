package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestSurfaceMappingsLifecycle(t *testing.T) {
	f := newConfigFixture(t)
	config := func(repository, directory, project string) string {
		return fmt.Sprintf(`
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
resource "interfere_surface_repository" "test" {
workspace_slug="example"
surface_slug=%q
integration_id=data.interfere_integration.test.id
repository_id=%q
working_directory=%q
}
resource "interfere_surface_destination" "test" {
workspace_slug="example"
surface_slug=%q
integration_id="destination-1"
project_id=%q
}
`, f.server.URL, surfaceSlug, repository, directory, surfaceSlug, project)
	}
	factories := map[string]func() (tfprotov6.ProviderServer, error){"interfere": providerserver.NewProtocol6WithError(New("test")())}
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{
		{Config: config("repo-1", "apps/web", "project-1"), Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr("interfere_surface_repository.test", "id", surfaceID), resource.TestCheckResourceAttr("interfere_surface_destination.test", "project_id", "project-1"), resource.TestCheckResourceAttr("data.interfere_workspace.test", "name", "Example"), resource.TestCheckResourceAttr("data.interfere_integration.test", "status", "connected"))},
		{ResourceName: "interfere_surface_repository.test", ImportState: true, ImportStateId: "example/" + surfaceSlug, ImportStateVerify: true},
		{ResourceName: "interfere_surface_destination.test", ImportState: true, ImportStateId: "example/" + surfaceSlug, ImportStateVerify: true},
		{Config: config("repo-2", "apps/api", "project-2"), Check: resource.TestCheckResourceAttr("interfere_surface_repository.test", "repository_id", "repo-2")},
		{PreConfig: func() { f.mu.Lock(); f.repository = "drift"; f.project = "drift"; f.mu.Unlock() }, Config: config("repo-2", "apps/api", "project-2")},
		{PreConfig: func() {
			f.mu.Lock()
			f.source = nil
			f.repository = nil
			f.directory = nil
			f.destination = nil
			f.project = nil
			f.mu.Unlock()
		}, Config: config("repo-2", "apps/api", "project-2")},
	}, CheckDestroy: func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		if f.source != nil || f.destination != nil {
			return fmt.Errorf("Mappings were not unlinked")
		}
		return nil
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

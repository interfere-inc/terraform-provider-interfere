package provider

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
)

func TestConfigurationResourcesPreserveStateOnReadFailure(t *testing.T) {
	for _, status := range []int{401, 403, 404, 503} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			f := newConfigFixture(t)
			f.readStatus = status
			c := client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())
			cases := map[string]struct {
				resource resource.Resource
				model    any
			}{
				"repository":  {&surfaceRepositoryResource{client: c}, surfaceRepositoryModel{ID: types.StringValue(surfaceID), WorkspaceSlug: types.StringValue("example"), SurfaceSlug: types.StringValue(surfaceSlug), IntegrationID: types.StringValue("integration-1"), MappingID: types.StringValue("repo-1")}},
				"destination": {&surfaceDestinationResource{client: c}, surfaceDestinationModel{ID: types.StringValue(surfaceID), WorkspaceSlug: types.StringValue("example"), SurfaceSlug: types.StringValue(surfaceSlug), IntegrationID: types.StringValue("destination-1"), MappingID: types.StringValue("project-1")}},
				"domain":      {&trackingDomainResource{client: c}, trackingDomainModel{ID: types.StringValue(attemptID), WorkspaceSlug: types.StringValue("example"), Name: types.StringValue("track.example.com"), IdempotencyKey: types.StringValue(attemptID)}},
			}
			for name, tc := range cases {
				t.Run(name, func(t *testing.T) {
					var schema resource.SchemaResponse
					tc.resource.Schema(context.Background(), resource.SchemaRequest{}, &schema)
					state := tfsdk.State{Schema: schema.Schema}
					if ds := state.Set(context.Background(), tc.model); ds.HasError() {
						t.Fatal(ds)
					}
					response := resource.ReadResponse{State: state}
					tc.resource.Read(context.Background(), resource.ReadRequest{State: state}, &response)
					if !response.Diagnostics.HasError() || !response.State.Raw.Equal(state.Raw) {
						t.Fatalf("Expected preserved state and error: %v", response.Diagnostics)
					}
					deletion := resource.DeleteResponse{State: state}
					tc.resource.Delete(context.Background(), resource.DeleteRequest{State: state}, &deletion)
					if !deletion.Diagnostics.HasError() {
						t.Fatal("Deletion accepted after a failed read")
					}
				})
			}
		})
	}
}

func TestMappingRefusesToAdoptDifferentMapping(t *testing.T) {
	f := newConfigFixture(t)
	f.source = "other-integration"
	f.repository = "other-repository"
	r := &surfaceRepositoryResource{client: client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())}
	data := surfaceRepositoryModel{WorkspaceSlug: types.StringValue("example"), SurfaceSlug: types.StringValue(surfaceSlug), IntegrationID: types.StringValue("integration-1"), MappingID: types.StringValue("repo-1")}
	if err := r.write(context.Background(), &data, true); err == nil {
		t.Fatal("Creation silently replaced an existing mapping")
	}
	if f.source != "other-integration" || f.repository != "other-repository" {
		t.Fatal("Existing mapping was changed")
	}
}

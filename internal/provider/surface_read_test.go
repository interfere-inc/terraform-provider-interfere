package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
)

func fixtureResource(t *testing.T, f *surfaceFixture) (*surfaceResource, tfsdk.State) {
	t.Helper()
	r := &surfaceResource{client: client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())}
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	diagnostics := state.Set(context.Background(), surfaceModel{
		Id: types.StringValue(surfaceID), Slug: types.StringValue(surfaceSlug),
		Name: types.StringValue("Example"), Type: types.StringValue("react"),
		WorkspaceSlug: types.StringValue("example"), IdempotencyKey: types.StringValue(attemptID),
	})
	if diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	return r, state
}

func TestReadPreservesStateOnAmbiguousOrFailedResponses(t *testing.T) {
	cases := map[string]struct {
		status int
		body   any
	}{
		"ambiguous null": {0, json.RawMessage(`null`)},
		"incomplete":     {0, map[string]any{}},
		"unauthorized":   {401, map[string]any{"message": "fixture-token"}},
		"forbidden":      {403, map[string]any{"message": "fixture-token"}},
		"not found":      {404, map[string]any{"message": "fixture-token"}},
		"unavailable":    {503, map[string]any{"message": "fixture-token"}},
		"reused slug":    {0, map[string]any{"id": "different", "slug": surfaceSlug, "name": "Other", "type": "react", "deletedAt": nil}},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			f := newSurfaceFixture(t)
			f.readStatus, f.readBody = test.status, test.body
			r, state := fixtureResource(t, f)
			response := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &response)
			if !response.Diagnostics.HasError() || !response.State.Raw.Equal(state.Raw) {
				t.Fatalf("Expected error with preserved state: %v", response.Diagnostics)
			}
		})
	}
}

func TestReadRemovesExplicitlyDeletedSurface(t *testing.T) {
	f := newSurfaceFixture(t)
	f.deleted = true
	r, state := fixtureResource(t, f)
	response := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &response)
	if response.Diagnostics.HasError() || !response.State.Raw.IsNull() {
		t.Fatalf("Expected deleted surface to leave state: %v", response.Diagnostics)
	}
}

func TestReadRemovesConfirmedMissingSurface(t *testing.T) {
	f := newSurfaceFixture(t)
	f.readStatus = 404
	f.readBody = map[string]any{"code": "WORKSPACE_RESOURCE_NOT_FOUND"}
	r, state := fixtureResource(t, f)
	response := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &response)
	if response.Diagnostics.HasError() || !response.State.Raw.IsNull() {
		t.Fatalf("Expected confirmed missing surface to leave state: %v", response.Diagnostics)
	}
}

func TestDeleteAcceptsConfirmedMissingSurface(t *testing.T) {
	f := newSurfaceFixture(t)
	f.readStatus = 404
	f.readBody = map[string]any{"code": "WORKSPACE_RESOURCE_NOT_FOUND"}
	r, state := fixtureResource(t, f)
	response := resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("Expected deletion of missing surface to succeed: %v", response.Diagnostics)
	}
}

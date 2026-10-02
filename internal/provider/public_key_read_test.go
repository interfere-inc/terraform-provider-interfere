package provider

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
)

func fixturePublicKey(t *testing.T, f *publicKeyFixture) (*publicKeyResource, tfsdk.State) {
	t.Helper()
	r := &publicKeyResource{client: client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())}
	model := publicKeyModel{Id: types.StringValue(attemptID), IdempotencyKey: types.StringValue(attemptID), Name: types.StringValue("Browser"), WorkspaceSlug: types.StringValue("example"), SurfaceSlug: types.StringValue("example-surface"), Content: types.StringValue("interfere_public_us_original")}
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatal(diags)
	}
	return r, state
}

func TestPublicKeyFailurePreservesState(t *testing.T) {
	for name, scenario := range map[string]struct {
		status int
		body   any
	}{
		"null":                 {0, json.RawMessage(`null`)},
		"incomplete":           {0, map[string]any{}},
		"permission":           {403, map[string]string{"message": "credential-from-error"}},
		"outage":               {503, map[string]string{"message": "credential-from-error"}},
		"unclassified missing": {404, map[string]string{"message": "credential-from-error"}},
		"wrong identity":       {0, map[string]any{"id": "other", "name": "Browser", "surfaceSlug": "example-surface", "content": "credential-from-error", "revoked": false}},
		"wrong surface":        {0, map[string]any{"id": attemptID, "name": "Browser", "surfaceSlug": "other", "content": "credential-from-error", "revoked": false}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newPublicKeyFixture(t)
			f.readStatus, f.readBody = scenario.status, scenario.body
			r, state := fixturePublicKey(t, f)
			response := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &response)
			if !response.Diagnostics.HasError() || !response.State.Raw.Equal(state.Raw) {
				t.Fatalf("Expected error and preserved state: %v", response.Diagnostics)
			}
			for _, d := range response.Diagnostics {
				if strings.Contains(d.Detail(), "credential-from-error") {
					t.Fatal("Diagnostic exposed response content")
				}
			}
		})
	}
}

func TestPublicKeyKnownInactiveLeavesState(t *testing.T) {
	for _, missing := range []bool{false, true} {
		f := newPublicKeyFixture(t)
		if missing {
			f.readStatus = 404
			f.readBody = map[string]string{"code": "SURFACE_KEY_NOT_FOUND"}
		} else {
			f.readBody = map[string]any{"id": attemptID, "name": "Browser", "surfaceSlug": "example-surface", "content": "interfere_public_us_original", "revoked": true}
		}
		r, state := fixturePublicKey(t, f)
		response := resource.ReadResponse{State: state}
		r.Read(context.Background(), resource.ReadRequest{State: state}, &response)
		if response.Diagnostics.HasError() || !response.State.Raw.IsNull() {
			t.Fatalf("Expected inactive key to leave state: %v", response.Diagnostics)
		}
	}
}

func TestPublicKeyDeleteRequiresAcknowledgement(t *testing.T) {
	f := newPublicKeyFixture(t)
	f.keys[attemptID] = map[string]any{"id": attemptID, "name": "Browser", "surfaceSlug": "example-surface", "content": "interfere_public_us_original", "revoked": false}
	f.acknowledge = false
	r, state := fixturePublicKey(t, f)
	response := resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, &response)
	if !response.Diagnostics.HasError() || !response.State.Raw.Equal(state.Raw) {
		t.Fatalf("Expected error and preserved state: %v", response.Diagnostics)
	}
}

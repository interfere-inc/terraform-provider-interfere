package provider

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
)

func fixturePrivateKey(t *testing.T, f *keyFixture) (*privateKeyResource, tfsdk.State) {
	t.Helper()
	r := &privateKeyResource{client: client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())}
	var model privateKeyModel
	model.Id = types.StringValue(attemptID)
	model.IdempotencyKey = model.Id
	model.Version = types.StringValue("ak_original")
	model.Name = types.StringValue("Managed key")
	model.WorkspaceSlug = types.StringValue("example")
	model.Scopes = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("org:surfaces:read")})
	model.Secret = types.StringValue("fixture-secret")
	var response resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &response)
	if !response.Schema.Attributes["secret"].(schema.StringAttribute).Sensitive {
		t.Fatal("Secret must be sensitive")
	}
	state := tfsdk.State{Schema: response.Schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatal(diags)
	}
	return r, state
}

func TestPrivateKeyReadPreservesStateAndSecretOnFailure(t *testing.T) {
	for name, scenario := range map[string]struct {
		status int
		body   any
	}{
		"null":                 {0, json.RawMessage(`null`)},
		"incomplete":           {0, map[string]any{}},
		"permission":           {403, map[string]string{"message": "fixture-secret"}},
		"outage":               {503, map[string]string{"message": "fixture-secret"}},
		"unclassified missing": {404, map[string]string{"message": "fixture-secret"}},
		"rotation":             {0, map[string]any{"id": attemptID, "version": "ak_rotated", "name": "Managed key", "scopes": []string{"org:surfaces:read"}, "surfaceSlug": nil, "secondsUntilExpiration": nil, "expiresAt": nil, "revoked": false}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newKeyFixture(t)
			f.readStatus, f.readBody = scenario.status, scenario.body
			r, state := fixturePrivateKey(t, f)
			response := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &response)
			if !response.Diagnostics.HasError() || !response.State.Raw.Equal(state.Raw) {
				t.Fatalf("Expected preserved state and an error: %v", response.Diagnostics)
			}
			for _, diagnostic := range response.Diagnostics {
				if strings.Contains(diagnostic.Detail(), "fixture-secret") {
					t.Fatal("Diagnostic leaked a secret")
				}
			}
		})
	}
}

func TestPrivateKeyReadRemovesKnownInactiveKeys(t *testing.T) {
	for _, missing := range []bool{false, true} {
		f := newKeyFixture(t)
		if missing {
			f.readStatus = 404
			f.readBody = map[string]string{"code": "SURFACE_KEY_NOT_FOUND"}
		} else {
			f.readBody = map[string]any{"id": attemptID, "version": "ak_original", "name": "Managed key", "scopes": []string{"org:surfaces:read"}, "surfaceSlug": nil, "secondsUntilExpiration": nil, "expiresAt": nil, "revoked": true}
		}
		r, state := fixturePrivateKey(t, f)
		response := resource.ReadResponse{State: state}
		r.Read(context.Background(), resource.ReadRequest{State: state}, &response)
		if response.Diagnostics.HasError() || !response.State.Raw.IsNull() {
			t.Fatalf("Expected inactive key to leave state: %v", response.Diagnostics)
		}
	}
}

package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

const surfaceID = "60c82b45-5ea6-48b8-9ce4-f1409e17973a"
const attemptID = "0f92eedc-cad4-46c4-8d2e-5f1bc4d61fe0"
const surfaceSlug = "example-surface"

type surfaceFixture struct {
	mu                             sync.Mutex
	server                         *httptest.Server
	name, surfaceType              string
	deleted                        bool
	tracking                       bool
	creates, updates, deletes      int
	readBody                       any
	readStatus, writeStatus        int
	acknowledge, automaticIdentity bool
	createIdentities               []string
}

func newSurfaceFixture(t *testing.T) *surfaceFixture {
	t.Helper()
	f := &surfaceFixture{name: "Example", acknowledge: true, tracking: true, surfaceType: "react"}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer fixture-token" || r.Method != http.MethodPost {
			t.Error("Expected an authenticated POST")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		write := func(value any) {
			if err := json.NewEncoder(w).Encode(value); err != nil {
				t.Error(err)
			}
		}
		handlers := map[string]func(){
			"/v3/workspaces/example/surfaces": func() {
				f.creates++
				var body sdk.CreateSurfaceRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if _, err := uuid.Parse(body.IdempotencyKey); err != nil || (!f.automaticIdentity && body.IdempotencyKey != attemptID) {
					t.Errorf("Unexpected creation identity: %s", body.IdempotencyKey)
				}
				f.createIdentities = append(f.createIdentities, body.IdempotencyKey)
				if body.APIKey == nil || len(body.APIKey.Scopes) != 1 || body.APIKey.Scopes[0] != sdk.CreateSurfaceRequestAPIKeyScopesItemReleaseWrite {
					t.Error("Expected explicit release-only creation grants")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				var grants map[string]any
				if err := json.Unmarshal([]byte(body.APIKey.String()), &grants); err != nil {
					t.Error(err)
				}
				if expiration, present := grants["secondsUntilExpiration"]; !present || expiration != nil {
					t.Error("Expected explicit null creation expiry")
				}
				if f.writeStatus != 0 {
					w.WriteHeader(f.writeStatus)
					write(map[string]any{"message": "secret-fixture-token", "status": f.writeStatus})
					return
				}
				f.name, f.deleted, f.surfaceType = body.Name, false, string(body.Type)
				write(map[string]any{
					"surface":   map[string]any{"id": surfaceID, "slug": surfaceSlug, "name": f.name, "type": string(body.Type)},
					"apiKey":    map[string]any{"id": "key", "name": "Default", "secret": "must-not-enter-state"},
					"publicKey": "must-not-enter-state",
				})
			},
			"/v3/workspaces/example/queries/surfaces.getBySlugIncludeDeleted": func() {
				var body sdk.GetBySlugIncludeDeletedSurfacesRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if body.Args == nil || body.Args.SurfaceSlug != surfaceSlug {
					t.Error("Unexpected read identifier")
				}
				if f.readStatus != 0 {
					w.WriteHeader(f.readStatus)
					write(f.readBody)
					return
				}
				if f.readBody != nil {
					write(f.readBody)
					return
				}
				var deletedAt any
				if f.deleted {
					deletedAt = 1234567890
				}
				write(map[string]any{"id": surfaceID, "slug": surfaceSlug, "name": f.name, "type": f.surfaceType, "deletedAt": deletedAt, "anonymousUserTracking": f.tracking, "sourceIntegrationId": nil, "sourceMappingId": nil, "sourceWorkingDirectory": nil, "destinationIntegrationId": nil, "destinationMappingId": nil})
			},
			"/v3/workspaces/example/actions/surfaces.setAnonymousUserTracking": func() {
				var body sdk.SetAnonymousUserTrackingSurfacesRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if f.acknowledge {
					f.tracking = body.Args.Enabled
				}
				write(map[string]any{"success": f.acknowledge})
			},
			"/v3/workspaces/example/actions/surfaces.updateName": func() {
				f.updates++
				var body sdk.UpdateNameSurfacesRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if body.Args == nil || body.Args.SurfaceSlug != surfaceSlug {
					t.Error("Unexpected update identifier")
					return
				}
				f.name = body.Args.Name
				write(map[string]any{"success": f.acknowledge})
			},
			"/v3/workspaces/example/actions/surfaces.delete": func() {
				f.deletes++
				var body sdk.DeleteSurfacesRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if body.Args == nil || body.Args.SurfaceSlug != surfaceSlug {
					t.Error("Unexpected delete identifier")
				}
				f.deleted = f.acknowledge
				write(map[string]any{"success": f.acknowledge})
			},
		}
		if handler, ok := handlers[r.URL.Path]; ok {
			handler()
			return
		}
		t.Errorf("Unexpected path: %s", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(f.server.Close)
	return f
}

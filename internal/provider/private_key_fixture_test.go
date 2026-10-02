package provider

import (
	"encoding/json"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/keys"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type keyFixture struct {
	mu               sync.Mutex
	server           *httptest.Server
	keys             map[string]map[string]any
	creates, deletes int
	readStatus       int
	readBody         any
}

func newKeyFixture(t *testing.T) *keyFixture {
	t.Helper()
	f := &keyFixture{keys: map[string]map[string]any{}}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer fixture-token" {
			t.Error("Missing bearer token")
			w.WriteHeader(401)
			return
		}
		write := func(value any) {
			if err := json.NewEncoder(w).Encode(value); err != nil {
				t.Error(err)
			}
		}
		root := "/v3/workspaces/example/api-keys"
		if r.URL.Path == root && r.Method == http.MethodPost {
			payload, err := io.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
				return
			}
			var body keys.CreatePrivateRequest
			if err := json.Unmarshal(payload, &body); err != nil {
				t.Error(err)
				return
			}
			var raw map[string]any
			if err := json.Unmarshal(payload, &raw); err != nil {
				t.Error(err)
			}
			if _, present := raw["secondsUntilExpiration"]; !present {
				t.Error("Missing explicit key expiry")
			}
			if old := f.keys[body.IdempotencyKey]; old != nil {
				w.WriteHeader(409)
				write(map[string]string{"code": "SURFACE_KEY_CONFLICT"})
				return
			}
			f.creates++
			var surface any
			if body.SurfaceSlug != nil {
				surface = *body.SurfaceSlug
			}
			var expires any
			if body.SecondsUntilExpiration != nil {
				expires = 4000000000000
			}
			f.keys[body.IdempotencyKey] = map[string]any{"id": body.IdempotencyKey, "version": "ak_" + body.IdempotencyKey, "name": body.Name, "surfaceSlug": surface, "scopes": body.Scopes, "secondsUntilExpiration": body.SecondsUntilExpiration, "expiresAt": expires, "revoked": false}
			write(map[string]any{"apiKey": map[string]string{"id": body.IdempotencyKey, "name": body.Name, "secret": "fixture-secret-" + body.IdempotencyKey}})
			return
		}
		if strings.HasPrefix(r.URL.Path, root+"/") {
			id := strings.TrimPrefix(r.URL.Path, root+"/")
			key := f.keys[id]
			if r.Method == http.MethodGet {
				if f.readStatus != 0 {
					w.WriteHeader(f.readStatus)
					write(f.readBody)
					return
				}
				if f.readBody != nil {
					write(f.readBody)
					return
				}
				if key == nil {
					w.WriteHeader(404)
					write(map[string]string{"code": "SURFACE_KEY_NOT_FOUND"})
					return
				}
				write(key)
				return
			}
			if r.Method == http.MethodDelete && key != nil {
				key["revoked"] = true
				f.deletes++
				write(map[string]bool{"success": true})
				return
			}
		}
		t.Errorf("Unexpected key API request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(404)
	}))
	t.Cleanup(f.server.Close)
	return f
}

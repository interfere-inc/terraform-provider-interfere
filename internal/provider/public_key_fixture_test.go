package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/keys"
)

type publicKeyFixture struct {
	mu               sync.Mutex
	server           *httptest.Server
	keys             map[string]map[string]any
	creates, deletes int
	readStatus       int
	readBody         any
	acknowledge      bool
}

func newPublicKeyFixture(t *testing.T) *publicKeyFixture {
	t.Helper()
	f := &publicKeyFixture{keys: map[string]map[string]any{}, acknowledge: true}
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
		root := "/v3/workspaces/example/surfaces/example-surface/publicKeys"
		if r.URL.Path == root && r.Method == http.MethodPost {
			var body keys.CreatePublicRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
				return
			}
			key := f.keys[body.IdempotencyKey]
			if key != nil && (key["revoked"] == true || key["name"] != body.Name) {
				w.WriteHeader(409)
				write(map[string]string{"code": "SURFACE_KEY_CONFLICT"})
				return
			}
			if key == nil {
				f.creates++
				key = map[string]any{"id": body.IdempotencyKey, "name": body.Name, "surfaceSlug": "example-surface", "content": "interfere_public_us_" + body.IdempotencyKey, "revoked": false}
				f.keys[body.IdempotencyKey] = key
			}
			write(map[string]any{"content": key["content"], "name": key["name"]})
			return
		}
		if strings.HasPrefix(r.URL.Path, root+"/") {
			key := f.keys[strings.TrimPrefix(r.URL.Path, root+"/")]
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
				f.deletes++
				if f.acknowledge {
					key["revoked"] = true
				}
				write(map[string]bool{"success": f.acknowledge})
				return
			}
		}
		t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(404)
	}))
	t.Cleanup(f.server.Close)
	return f
}

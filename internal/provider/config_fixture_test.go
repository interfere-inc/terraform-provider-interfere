package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type configFixture struct {
	mu                                                  sync.Mutex
	server                                              *httptest.Server
	source, repository, directory, destination, project any
	domain, domainDeleted                               bool
	domainStatus                                        string
	readStatus                                          int
	acknowledge                                         bool
	integrations                                        []map[string]any
}

func newConfigFixture(t *testing.T) *configFixture {
	t.Helper()
	f := &configFixture{acknowledge: true, domainStatus: "pending", integrations: []map[string]any{{"id": "integration-1", "provider": "github", "status": "connected"}}}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer fixture-token" || r.Method != http.MethodPost {
			t.Error("Expected authenticated POST")
			w.WriteHeader(401)
			return
		}
		var body struct {
			Args map[string]any `json:"args"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		write := func(v any) {
			if err := json.NewEncoder(w).Encode(v); err != nil {
				t.Error(err)
			}
		}
		if strings.Contains(r.URL.Path, "/queries/") && f.readStatus != 0 {
			w.WriteHeader(f.readStatus)
			write(map[string]any{"message": "fixture-token"})
			return
		}
		handlers := map[string]func(){
			"queries/surfaces.getBySlugIncludeDeleted": func() {
				if body.Args["surfaceSlug"] != surfaceSlug {
					t.Error("Wrong surface")
				}
				write(map[string]any{"id": surfaceID, "slug": surfaceSlug, "name": "Example", "type": "react", "deletedAt": nil, "anonymousUserTracking": true, "sourceIntegrationId": f.source, "sourceMappingId": f.repository, "sourceWorkingDirectory": f.directory, "destinationIntegrationId": f.destination, "destinationMappingId": f.project})
			},
			"queries/organizations.current": func() {
				write(map[string]any{"id": surfaceID, "slug": "example", "name": "Example", "dataResidencyLocation": "us"})
			},
			"queries/integrations.installations": func() { write(f.integrations) },
			"queries/domains.byIdIncludeDeleted": func() {
				if body.Args["domainId"] != attemptID {
					t.Error("Wrong domain identity")
				}
				if !f.domain {
					write(nil)
					return
				}
				var deleted any
				if f.domainDeleted {
					deleted = 1234
				}
				write(map[string]any{"id": attemptID, "name": "track.example.com", "type": "proxy", "deletedAt": deleted, "status": f.domainStatus, "domainMetadata": map[string]any{"cfHostnameId": "cf-domain", "cfSslStatus": f.domainStatus, "lastError": nil}})
			},
			"actions/integrations.linkSurfaceToRepository": func() {
				if f.acknowledge {
					f.source = body.Args["integrationId"]
					f.repository = body.Args["repositoryId"]
					f.directory = body.Args["workingDirectory"]
				}
				write(map[string]any{"success": f.acknowledge})
			},
			"actions/integrations.linkSurfaceToDestination": func() {
				if f.acknowledge {
					f.destination = body.Args["integrationId"]
					f.project = body.Args["projectId"]
				}
				write(map[string]any{"success": f.acknowledge})
			},
			"actions/integrations.unlinkSurfaceMapping": func() {
				if f.acknowledge {
					f.source = nil
					f.repository = nil
					f.directory = nil
				}
				write(map[string]any{"success": f.acknowledge})
			},
			"actions/integrations.unlinkSurfaceDestinationMapping": func() {
				if f.acknowledge {
					f.destination = nil
					f.project = nil
				}
				write(map[string]any{"success": f.acknowledge})
			},
			"actions/organizations.addProxyDomain": func() {
				if body.Args["id"] != attemptID || body.Args["name"] != "track.example.com" {
					t.Error("Wrong domain creation")
				}
				if f.acknowledge {
					f.domain = true
					f.domainDeleted = false
				}
				write(map[string]any{"success": f.acknowledge})
			},
			"actions/organizations.removeProxyDomain": func() {
				if body.Args["name"] != "track.example.com" || body.Args["id"] != attemptID {
					t.Error("Wrong domain removal")
				}
				if f.acknowledge {
					f.domainDeleted = true
				}
				write(map[string]any{"success": f.acknowledge})
			},
		}
		path := strings.TrimPrefix(r.URL.Path, "/v3/workspaces/example/")
		if handle, ok := handlers[path]; ok {
			handle()
			return
		}
		t.Errorf("Unexpected request: %s", r.URL.Path)
		w.WriteHeader(404)
	}))
	t.Cleanup(f.server.Close)
	return f
}

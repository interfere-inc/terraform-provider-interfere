OPENAPI ?=
FERN ?= fern

.PHONY: generate build test check docs

generate:
	test -n "$(OPENAPI)"
	cp "$(OPENAPI)" fern/openapi.json
	$(FERN) generate --local --force

build:
	go build -trimpath -o bin/terraform-provider-interfere .

test:
	go test ./internal/provider -count=1

check:
	test -z "$$(gofmt -l main.go internal/provider)"
	go vet ./internal/provider
	terraform fmt -check -recursive examples

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.25.0 generate --rendered-provider-name Interfere

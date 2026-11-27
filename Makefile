FERN ?= bunx --package fern-api@5.193.2 fern
OPENAPI_GENERATOR ?= go run github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@v0.3.0
SPEC ?= openapi.json
FRAMEWORK_GENERATOR ?= go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@v0.4.1

.PHONY: generate generate-sdk generate-schemas check-generated check-build build test check docs
.NOTPARALLEL: generate check-generated

generate: generate-sdk generate-schemas docs

generate-sdk:
	$(FERN) generate --local --force

generate-schemas:
	mkdir -p bin
	$(OPENAPI_GENERATOR) generate --config generator_config.yml --output bin/provider-code-spec.json openapi.json
	$(FRAMEWORK_GENERATOR) generate resources --input bin/provider-code-spec.json --output internal
	$(FRAMEWORK_GENERATOR) generate data-sources --input bin/provider-code-spec.json --output internal
	shasum -a 256 openapi.json | cut -d ' ' -f 1 > openapi.sha256

check-generated: generate-schemas docs
	git diff --exit-code -- openapi.sha256 'internal/resource_*' 'internal/datasource_*' docs
	test -z "$$(git ls-files --others --exclude-standard -- 'internal/resource_*' 'internal/datasource_*' docs)"

build:
	go build -trimpath -o bin/terraform-provider-interfere .

test:
	go test ./internal/provider $(addprefix ./,$(wildcard internal/resource_* internal/datasource_*)) -count=1

check:
	test -z "$$(gofmt -l main.go internal/provider)"
	go vet ./internal/provider
	terraform fmt -check -recursive examples

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.25.0 generate --provider-name interfere --rendered-provider-name Interfere

check-build:
	@command -v go >/dev/null && command -v docker >/dev/null || { echo "check-build needs go and docker to generate the Go SDK and compile the provider" >&2; exit 1; }
	test "$$(shasum -a 256 $(SPEC) | cut -d ' ' -f 1)" = "$$(cat openapi.sha256)"
	if [ "$(SPEC)" != openapi.json ]; then cp $(SPEC) openapi.json; fi
	$(MAKE) generate-sdk generate-schemas
	git diff --exit-code -- openapi.sha256 'internal/resource_*' 'internal/datasource_*'
	go build ./...
	go vet ./internal/provider
	go test ./internal/provider -count=1

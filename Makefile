FERN ?= fern
OPENAPI_GENERATOR ?= go run github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@v0.3.0
FRAMEWORK_GENERATOR ?= go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@v0.4.1

.PHONY: generate generate-sdk generate-schemas check-generated build test check docs
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

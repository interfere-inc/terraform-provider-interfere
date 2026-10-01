OPENAPI ?=
FERN ?= fern
TOOLS := $(CURDIR)/.tools

.PHONY: tools generate build test check docs

tools:
	GOBIN=$(TOOLS) go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@v0.3.0
	GOBIN=$(TOOLS) go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@v0.4.1

generate: tools
	test -n "$(OPENAPI)"
	cp "$(OPENAPI)" fern/openapi.json
	$(TOOLS)/tfplugingen-openapi generate --config generator_config.yml --output provider-code-spec.json fern/openapi.json
	$(TOOLS)/tfplugingen-framework generate resources --input provider-code-spec.json --output internal
	$(FERN) generate --local --force

build:
	go build -trimpath -o bin/terraform-provider-interfere .

test:
	go test ./internal/provider -count=1

check:
	test -z "$$(gofmt -l main.go internal/provider)"
	go vet ./internal/provider ./internal/resource_*
	terraform fmt -check -recursive examples

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.25.0 generate --rendered-provider-name Interfere

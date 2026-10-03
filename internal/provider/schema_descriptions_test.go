package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

func TestEveryPublishedAttributeHasDocumentation(t *testing.T) {
	ctx := context.Background()
	provider := &interfereProvider{}
	for _, create := range provider.Resources(ctx) {
		instance := create()
		var metadata resource.MetadataResponse
		instance.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "interfere"}, &metadata)
		var response resource.SchemaResponse
		instance.Schema(ctx, resource.SchemaRequest{}, &response)
		require.False(t, response.Diagnostics.HasError())
		require.Zero(t, response.Schema.Version, metadata.TypeName)
		for name, attribute := range response.Schema.Attributes {
			require.NotEmpty(t, attribute.GetDescription(), "%s.%s", metadata.TypeName, name)
		}
	}
	for _, create := range provider.DataSources(ctx) {
		instance := create()
		var metadata datasource.MetadataResponse
		instance.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "interfere"}, &metadata)
		var response datasource.SchemaResponse
		instance.Schema(ctx, datasource.SchemaRequest{}, &response)
		require.False(t, response.Diagnostics.HasError())
		for name, attribute := range response.Schema.Attributes {
			require.NotEmpty(t, attribute.GetDescription(), "%s.%s", metadata.TypeName, name)
		}
	}
}

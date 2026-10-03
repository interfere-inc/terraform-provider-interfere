package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceSchemaPreservesV030Contract(t *testing.T) {
	ctx := context.Background()
	var managed resource.SchemaResponse
	NewWorkspaceResource().Schema(ctx, resource.SchemaRequest{}, &managed)
	var lookup datasource.SchemaResponse
	NewWorkspaceDataSource().Schema(ctx, datasource.SchemaRequest{}, &lookup)
	require.Zero(t, managed.Schema.Version)
	require.Len(t, managed.Schema.Attributes, 4)
	require.Len(t, lookup.Schema.Attributes, 4)
	for _, field := range []string{"id", "name", "workspace_slug", "data_residency_location"} {
		attribute := managed.Schema.Attributes[field]
		require.Equal(t, types.StringType, attribute.GetType())
		require.Equal(t, field == "name" || field == "workspace_slug", attribute.IsRequired())
		require.Equal(t, field == "id" || field == "data_residency_location", attribute.IsComputed())
		require.False(t, attribute.IsOptional())
		require.False(t, attribute.IsSensitive())
		require.NotEmpty(t, attribute.GetDescription())
		read := lookup.Schema.Attributes[field]
		require.Equal(t, types.StringType, read.GetType())
		require.Equal(t, field == "workspace_slug", read.IsRequired())
		require.Equal(t, field != "workspace_slug", read.IsComputed())
		require.False(t, read.IsOptional())
		require.NotEmpty(t, read.GetDescription())
	}
}

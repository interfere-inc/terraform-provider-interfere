package provider

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
	"github.com/stretchr/testify/require"
)

func TestCreationRecoveryAfterExhaustedRetries(t *testing.T) {
	ctx := context.Background()
	f := newPublicKeyFixture(t)
	f.createFailures = 3
	r := &publicKeyResource{client: client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())}
	var schema resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schema)
	data := publicKeyModel{
		Id: types.StringUnknown(), IdempotencyKey: types.StringUnknown(),
		WorkspaceSlug: types.StringValue("example"), SurfaceSlug: types.StringValue(surfaceSlug),
		Name: types.StringValue("Recovery"), Content: types.StringUnknown(),
	}
	plan := tfsdk.Plan{Schema: schema.Schema}
	require.False(t, plan.Set(ctx, data).HasError())
	response := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema, Raw: plan.Raw}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &response)
	require.True(t, response.Diagnostics.HasError())
	require.Len(t, f.createIdentities, 3)
	require.Equal(t, 1, f.creates)
	for _, id := range f.createIdentities {
		require.Equal(t, f.createIdentities[0], id)
	}
	detail := response.Diagnostics.Errors()[0].Detail()
	require.NotContains(t, detail, "fixture-secret-must-not-leak")
	identity := regexp.MustCompile("Creation UUID: ([a-f0-9-]+)").FindStringSubmatch(detail)
	require.Len(t, identity, 2)
	require.Equal(t, f.createIdentities[0], identity[1])
	data.IdempotencyKey = types.StringValue(identity[1])
	require.False(t, plan.Set(ctx, data).HasError())
	retry := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema, Raw: plan.Raw}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &retry)
	require.False(t, retry.Diagnostics.HasError(), retry.Diagnostics)
	require.Equal(t, 1, f.creates)
	require.False(t, retry.State.Get(ctx, &data).HasError())
	require.Equal(t, identity[1], data.Id.ValueString())
	require.True(t, strings.HasSuffix(data.Content.ValueString(), identity[1]))
}

func TestCreationDoesNotRetryPermissionFailures(t *testing.T) {
	ctx := context.Background()
	f := newSurfaceFixture(t)
	f.writeStatus = 403
	r := &surfaceResource{client: client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())}
	var schema resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schema)
	plan := tfsdk.Plan{Schema: schema.Schema}
	require.False(t, plan.Set(ctx, surfaceModel{IdempotencyKey: types.StringValue(attemptID), WorkspaceSlug: types.StringValue("example"), Name: types.StringValue("Example"), Type: types.StringValue("react")}).HasError())
	response := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema, Raw: plan.Raw}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &response)
	require.True(t, response.Diagnostics.HasError())
	require.Equal(t, 1, f.creates)
}

func TestPrivateKeyRecoveryRetainsTheOriginalSecret(t *testing.T) {
	ctx := context.Background()
	f := newKeyFixture(t)
	f.createFailures = 3
	r := &privateKeyResource{client: client.NewClient(option.WithToken("fixture-token"), option.WithBaseURL(f.server.URL), option.WithoutRetries())}
	var schema resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schema)
	data := privateKeyModel{
		IdempotencyKey: types.StringUnknown(), WorkspaceSlug: types.StringValue("example"),
		Name:   types.StringValue("Recovery"),
		Scopes: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("org:surfaces:read")}),
	}
	plan := tfsdk.Plan{Schema: schema.Schema}
	require.False(t, plan.Set(ctx, data).HasError())
	failed := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema, Raw: plan.Raw}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &failed)
	require.True(t, failed.Diagnostics.HasError())
	require.Len(t, f.createIdentities, 3)
	identity := f.createIdentities[0]
	detail := failed.Diagnostics.Errors()[0].Detail()
	require.Contains(t, detail, "Creation UUID: "+identity)
	require.NotContains(t, detail, "fixture-secret")
	require.False(t, plan.SetAttribute(ctx, path.Root("idempotency_key"), types.StringValue(identity)).HasError())
	recovered := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema, Raw: plan.Raw}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &recovered)
	require.False(t, recovered.Diagnostics.HasError(), recovered.Diagnostics)
	require.False(t, recovered.State.Get(ctx, &data).HasError())
	require.Equal(t, 1, f.creates)
	require.Equal(t, "fixture-secret-"+identity, data.Secret.ValueString())
	require.Equal(t, identity, data.Id.ValueString())
	for _, attempted := range f.createIdentities {
		require.Equal(t, identity, attempted)
	}
}

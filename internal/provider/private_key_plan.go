package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *privateKeyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var prior, planned privateKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planned)...)
	if resp.Diagnostics.HasError() || planned.IdempotencyKey.IsUnknown() || !prior.IdempotencyKey.Equal(planned.IdempotencyKey) {
		return
	}
	changed := !prior.Name.Equal(planned.Name) || !prior.Scopes.Equal(planned.Scopes) || !prior.SecondsUntilExpiration.Equal(planned.SecondsUntilExpiration) || !prior.SurfaceSlug.Equal(planned.SurfaceSlug) || !prior.WorkspaceSlug.Equal(planned.WorkspaceSlug)
	if changed {
		resp.Diagnostics.AddError("New key creation identity required", "Changing private-key configuration replaces the key. Set idempotency_key to a fresh UUID before applying so the existing key is not revoked before a conflicting creation attempt.")
	}
}

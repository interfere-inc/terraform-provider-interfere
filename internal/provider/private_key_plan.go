package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *privateKeyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	planCreationIdentity(ctx, req, resp, "name", "scopes", "seconds_until_expiration", "surface_slug", "workspace_slug")
}

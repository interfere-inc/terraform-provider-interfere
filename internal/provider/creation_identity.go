package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func creationIdentityAttribute(replaceWhenConfigured bool) schema.StringAttribute {
	attribute := schema.StringAttribute{
		Optional: true, Computed: true,
		Description: "Creation UUID, generated automatically when omitted. Explicit UUIDs remain supported for recovering an uncertain creation. Reuse that UUID to retry; use a fresh UUID for replacement.",
		Validators:  []validator.String{stringvalidator.RegexMatches(uuidPattern, "Must be a UUID."), stringvalidator.NoneOf(uuid.Nil.String())},
	}
	if replaceWhenConfigured {
		attribute.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()}
	}
	return attribute
}

func creationIdentity(planned types.String) types.String {
	if planned.IsNull() || planned.IsUnknown() {
		return types.StringValue(uuid.NewString())
	}
	return planned
}

func creationError(identity types.String, detail string) string {
	return fmt.Sprintf("%s Creation UUID: %s. If the result is uncertain, set idempotency_key to this UUID before retrying to recover the same creation.", detail, identity.ValueString())
}

func planCreationIdentity(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, replacementFields ...string) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var configured, prior, planned types.String
	keyPath := path.Root("idempotency_key")
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, keyPath, &configured)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, keyPath, &prior)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, keyPath, &planned)...)
	var before, after map[string]tftypes.Value
	if err := req.State.Raw.As(&before); err != nil {
		resp.Diagnostics.AddError("Unable to read creation state", err.Error())
	}
	if err := req.Plan.Raw.As(&after); err != nil {
		resp.Diagnostics.AddError("Unable to read creation plan", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	changed := false
	for _, field := range replacementFields {
		changed = changed || !before[field].Equal(after[field])
	}
	if configured.IsNull() {
		if changed {
			prior = types.StringUnknown()
		}
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, keyPath, prior)...)
		return
	}
	if changed && !planned.IsUnknown() && prior.Equal(planned) {
		resp.Diagnostics.AddError("New key creation identity required", "Configuration changes replace this resource. Omit idempotency_key to generate a fresh UUID automatically, or configure a new UUID before applying.")
	}
}

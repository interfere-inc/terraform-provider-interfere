package provider

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func (r *networkRulesResource) read(ctx context.Context, workspaceSlug string) ([]*sdk.ReadNetworkRulesResponseRulesItem, error) {
	remote, err := r.client.Workspaces.ReadNetworkRules(ctx, &sdk.ReadNetworkRulesRequest{WorkspaceSlug: workspaceSlug})
	if err != nil {
		return nil, errors.New(networkRulesError(err))
	}
	if remote == nil {
		return nil, errors.New("The API returned an empty network rules response. Terraform state was preserved.")
	}
	for _, rule := range remote.GetRules() {
		if rule == nil || rule.ID == "" || rule.Value == "" {
			return nil, errors.New("The API returned an incomplete network rule. Terraform state was preserved.")
		}
	}
	return remote.GetRules(), nil
}

func (r *networkRulesResource) save(ctx context.Context, workspaceSlug string, desired []networkRuleModel) error {
	remote, err := r.read(ctx, workspaceSlug)
	if err != nil {
		return err
	}
	existing := make(map[string]string, len(remote))
	for _, rule := range remote {
		existing[rule.Value] = rule.ID
	}
	rules := make([]*sdk.SaveNetworkRulesRequestRulesItem, 0, len(desired))
	for _, rule := range desired {
		id, ok := existing[rule.Value.ValueString()]
		if !ok {
			id = uuid.NewString()
		}
		rules = append(rules, &sdk.SaveNetworkRulesRequestRulesItem{ID: id, Label: rule.Label.ValueString(), Value: rule.Value.ValueString()})
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].Value < rules[j].Value })
	result, err := r.client.Workspaces.SaveNetworkRules(ctx, &sdk.SaveNetworkRulesRequest{WorkspaceSlug: workspaceSlug, Rules: rules})
	if err != nil {
		return errors.New(networkRulesError(err))
	}
	if result == nil || !result.Success {
		return errors.New("The API did not acknowledge the network rules update. Terraform state was preserved.")
	}
	return nil
}

func networkRulesError(err error) string {
	var conflict *sdk.ConflictError
	if errors.As(err, &conflict) && apiErrorCode(conflict.Body) == "NETWORK_POLICY_LOCKOUT" {
		return "The API rejected these rules because they would block the network this request comes from. Add a rule matching the Terraform runner's egress IP address, range or ASN."
	}
	var forbidden *sdk.ForbiddenError
	if errors.As(err, &forbidden) {
		switch apiErrorCode(forbidden.Body) {
		case "NETWORK_PLAN_REQUIRED":
			return "Network restrictions require an Enterprise plan. An empty rules list is always allowed."
		case "NETWORK_ACCESS_DENIED":
			return "The workspace network rules block the network this request comes from. An administrator on an allowed network must add the Terraform runner's egress before Terraform can manage the rules."
		}
	}
	return apiError(err)
}

func apiErrorCode(body any) string {
	if fields, ok := body.(map[string]interface{}); ok {
		if code, ok := fields["code"].(string); ok {
			return code
		}
	}
	return ""
}

func networkRulesValue(remote []*sdk.ReadNetworkRulesResponseRulesItem) (types.Set, diag.Diagnostics) {
	elements := make([]attr.Value, 0, len(remote))
	for _, rule := range remote {
		elements = append(elements, types.ObjectValueMust(networkRuleAttributeTypes, map[string]attr.Value{
			"label": types.StringValue(rule.Label),
			"value": types.StringValue(rule.Value),
		}))
	}
	return types.SetValue(types.ObjectType{AttrTypes: networkRuleAttributeTypes}, elements)
}

func (r *networkRulesResource) apply(ctx context.Context, data *networkRulesModel) diag.Diagnostics {
	var diags diag.Diagnostics
	var desired []networkRuleModel
	diags.Append(data.Rules.ElementsAs(ctx, &desired, false)...)
	if diags.HasError() {
		return diags
	}
	slug := data.WorkspaceSlug.ValueString()
	if err := r.save(ctx, slug, desired); err != nil {
		diags.AddError("Unable to save network rules", err.Error())
		return diags
	}
	remote, err := r.read(ctx, slug)
	if err != nil {
		diags.AddError("Unable to read saved network rules", err.Error())
		return diags
	}
	rules, valueDiags := networkRulesValue(remote)
	diags.Append(valueDiags...)
	data.Rules = rules
	return diags
}

func (r *networkRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data networkRulesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(r.apply(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *networkRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data networkRulesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(r.apply(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *networkRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data networkRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data.WorkspaceSlug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read network rules", err.Error())
		return
	}
	rules, diags := networkRulesValue(remote)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Rules = rules
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *networkRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data networkRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.save(ctx, data.WorkspaceSlug.ValueString(), nil); err != nil {
		resp.Diagnostics.AddError("Unable to remove network rules", err.Error())
	}
}

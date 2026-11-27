package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_network_rules"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

const networkRulesDescription = "Manage the complete list of networks allowed to reach a workspace's dashboard and API. Each apply replaces every rule in the workspace, including rules added in the dashboard. Destroying the resource saves an empty list, which removes all network restrictions. The network the apply runs from must match a rule, or the API rejects the change; if the runner's egress later leaves the list, Terraform loses access to the workspace until an administrator on an allowed network fixes the rules. Import using the workspace slug."

var networkRuleAttributeTypes = map[string]attr.Type{"label": types.StringType, "value": types.StringType}

type networkRulesResource struct{ client *client.Client }

type networkRulesModel struct {
	WorkspaceSlug types.String `tfsdk:"workspace_slug"`
	Rules         types.Set    `tfsdk:"rules"`
}

type networkRuleModel struct {
	Label types.String `tfsdk:"label"`
	Value types.String `tfsdk:"value"`
}

func NewNetworkRulesResource() resource.Resource { return &networkRulesResource{} }

func (r *networkRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_rules"
}

func (r *networkRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	configured, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Invalid provider configuration", "Expected an Interfere API client.")
		return
	}
	r.client = configured
}

func (r *networkRulesResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	generated := resource_network_rules.NetworkRulesResourceSchema(ctx)
	rules := generated.Attributes["rules"].(schema.ListNestedAttribute)
	label := rules.NestedObject.Attributes["label"].(schema.StringAttribute)
	label.Validators = []validator.String{stringvalidator.LengthAtMost(100), networkRuleLabelValidator{}}
	value := rules.NestedObject.Attributes["value"].(schema.StringAttribute)
	value.Validators = []validator.String{networkRuleValueValidator{}}
	resp.Schema = schema.Schema{
		Description: networkRulesDescription,
		Attributes: map[string]schema.Attribute{
			"workspace_slug": requiredString(generated.Attributes["workspace_slug"], stringplanmodifier.RequiresReplace()),
			"rules": schema.SetNestedAttribute{
				Required:    true,
				Description: rules.Description + " Every apply replaces the whole list.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{"label": label, "value": value},
				},
				Validators: []validator.Set{setvalidator.SizeAtMost(100), uniqueNetworkRuleValues{}},
			},
		},
	}
}

func (r *networkRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !slugPattern.MatchString(req.ID) {
		resp.Diagnostics.AddError("Invalid network rules import ID", "Use the existing workspace slug.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_slug"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rules"), types.SetValueMust(types.ObjectType{AttrTypes: networkRuleAttributeTypes}, []attr.Value{}))...)
}

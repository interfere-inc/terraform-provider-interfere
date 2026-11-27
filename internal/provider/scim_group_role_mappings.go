package provider

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_scim_group_role_mappings"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/scim"
)

var scimRoles = []string{"org:owner", "org:admin", "org:member", "org:product", "org:support", "org:viewer", "org:billing", "org:auditor"}

type scimGroupRoleMappingsResource struct{ client *client.Client }
type scimGroupRoleMapping struct {
	SCIMGroupID types.String `tfsdk:"scim_group_id"`
	Role        types.String `tfsdk:"role"`
}
type scimGroupRoleMappingsModel struct {
	ID            types.String           `tfsdk:"id"`
	WorkspaceSlug types.String           `tfsdk:"workspace_slug"`
	DirectoryID   types.String           `tfsdk:"directory_id"`
	Mappings      []scimGroupRoleMapping `tfsdk:"mappings"`
}

func NewSCIMGroupRoleMappingsResource() resource.Resource { return &scimGroupRoleMappingsResource{} }
func (r *scimGroupRoleMappingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_group_role_mappings"
}

func (r *scimGroupRoleMappingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *scimGroupRoleMappingsResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_scim_group_role_mappings.ScimGroupRoleMappingsResourceSchema(ctx)
	resp.Schema.Description = "Manage which workspace role each group in a SCIM directory grants, and their order. This resource owns every role mapping in the directory: groups left out of mappings have no role. When a member is in several mapped groups, the earliest mapping wins. Groups appear after the identity provider provisions them, and the directory must have group role mapping turned on. The directory itself, its connection and its SCIM bearer token stay outside Terraform; the token is never read or stored in state. Destroy removes the role from every mapped group. Import with workspace-slug/directory-id."
	attributes := resp.Schema.Attributes
	flattenAttributes(attributes, "args")
	attributes["workspace_slug"] = requiredString(attributes["workspace_slug"], stringplanmodifier.RequiresReplace())
	attributes["directory_id"] = requiredString(attributes["clerk_scim_directory_id"], stringplanmodifier.RequiresReplace())
	delete(attributes, "clerk_scim_directory_id")
	attributes["id"] = schema.StringAttribute{
		Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		Description: "Identifier of the SCIM directory, the same as directory_id.",
	}
	attributes["mappings"] = schema.ListNestedAttribute{
		Required:    true,
		Description: "Group role mappings in precedence order; earlier mappings win when a member is in several groups. Each group may appear once. An empty list removes every mapping.",
		Validators:  []validator.List{listvalidator.SizeAtMost(100)},
		NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"scim_group_id": schema.StringAttribute{Required: true, Description: "Identifier of the group in the SCIM directory, as the identity provider provisioned it."},
			"role":          schema.StringAttribute{Required: true, Description: "Workspace role granted to members of the group: org:owner, org:admin, org:member (Developer), org:product, org:support, org:viewer, org:billing or org:auditor. You can grant only a role whose permissions you hold.", Validators: []validator.String{stringvalidator.OneOf(scimRoles...)}},
		}},
	}
}

func (r *scimGroupRoleMappingsResource) groups(ctx context.Context, data scimGroupRoleMappingsModel) ([]*sdk.QuerySCIMGroupsForOrgResponseItem, error) {
	remote, err := r.client.SCIM.Groups.List(ctx, &scim.ListGroupsRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString()})
	if err != nil {
		return nil, errors.New(apiError(err))
	}
	groups := make([]*sdk.QuerySCIMGroupsForOrgResponseItem, 0, len(remote))
	for _, group := range remote {
		if group != nil && group.ClerkSCIMDirectoryID == data.DirectoryID.ValueString() && group.DeletedAt == nil {
			groups = append(groups, group)
		}
	}
	if len(groups) == 0 {
		return nil, errResourceNotFound
	}
	return groups, nil
}

func mappedGroups(groups []*sdk.QuerySCIMGroupsForOrgResponseItem) []*sdk.QuerySCIMGroupsForOrgResponseItem {
	mapped := slices.DeleteFunc(slices.Clone(groups), func(group *sdk.QuerySCIMGroupsForOrgResponseItem) bool { return group.RoleKey == nil })
	slices.SortStableFunc(mapped, func(a, b *sdk.QuerySCIMGroupsForOrgResponseItem) int {
		if a.Precedence == nil || b.Precedence == nil || *a.Precedence == *b.Precedence {
			return strings.Compare(a.SCIMGroupDisplayName, b.SCIMGroupDisplayName)
		}
		return *a.Precedence - *b.Precedence
	})
	return mapped
}

func (data *scimGroupRoleMappingsModel) refresh(groups []*sdk.QuerySCIMGroupsForOrgResponseItem) {
	data.ID = data.DirectoryID
	data.Mappings = []scimGroupRoleMapping{}
	for _, group := range mappedGroups(groups) {
		data.Mappings = append(data.Mappings, scimGroupRoleMapping{SCIMGroupID: types.StringValue(group.SCIMGroupID), Role: types.StringValue(*group.RoleKey)})
	}
}

func (r *scimGroupRoleMappingsResource) setRole(ctx context.Context, slug string, group *sdk.QuerySCIMGroupsForOrgResponseItem, role *string) error {
	args := &scim.SetRoleGroupsRequestArgs{}
	args.SetID(group.ID)
	var key *scim.SetRoleGroupsRequestArgsRoleKey
	if role != nil {
		parsed, err := scim.NewSetRoleGroupsRequestArgsRoleKeyFromString(*role)
		if err != nil {
			return err
		}
		key = &parsed
	}
	args.SetRoleKey(key)
	result, err := r.client.SCIM.Groups.SetRole(ctx, &scim.SetRoleGroupsRequest{WorkspaceSlug: slug, Args: args})
	return acknowledged(result != nil && result.Success, err, fmt.Sprintf("role mapping for group %q", group.SCIMGroupID))
}

func (r *scimGroupRoleMappingsResource) apply(ctx context.Context, data *scimGroupRoleMappingsModel) error {
	slug := data.WorkspaceSlug.ValueString()
	groups, err := r.groups(ctx, *data)
	if errors.Is(err, errResourceNotFound) {
		return errors.New("The directory has no provisioned groups in this workspace. Check directory_id and wait for the identity provider to provision groups.")
	}
	if err != nil {
		return err
	}
	byGroup := map[string]*sdk.QuerySCIMGroupsForOrgResponseItem{}
	for _, group := range groups {
		byGroup[group.SCIMGroupID] = group
	}
	desired := map[string]string{}
	for _, mapping := range data.Mappings {
		id := mapping.SCIMGroupID.ValueString()
		if _, duplicate := desired[id]; duplicate {
			return fmt.Errorf("Group %q appears more than once in mappings.", id)
		}
		if byGroup[id] == nil {
			return fmt.Errorf("Group %q has not been provisioned to this directory.", id)
		}
		desired[id] = mapping.Role.ValueString()
	}
	for _, group := range groups {
		if _, keep := desired[group.SCIMGroupID]; !keep && group.RoleKey != nil {
			if err := r.setRole(ctx, slug, group, nil); err != nil {
				return err
			}
		}
	}
	for _, mapping := range data.Mappings {
		group := byGroup[mapping.SCIMGroupID.ValueString()]
		role := mapping.Role.ValueString()
		if group.RoleKey == nil || *group.RoleKey != role {
			if err := r.setRole(ctx, slug, group, &role); err != nil {
				return err
			}
		}
	}
	if len(data.Mappings) == 0 {
		return nil
	}
	if groups, err = r.groups(ctx, *data); err != nil {
		return err
	}
	order := make([]string, 0, len(data.Mappings))
	for _, mapping := range data.Mappings {
		order = append(order, byGroup[mapping.SCIMGroupID.ValueString()].ID)
	}
	current := []string{}
	for _, group := range mappedGroups(groups) {
		current = append(current, group.ID)
	}
	if slices.Equal(current, order) {
		return nil
	}
	result, err := r.client.SCIM.Groups.Reorder(ctx, &scim.ReorderGroupsRequest{WorkspaceSlug: slug, Args: &scim.ReorderGroupsRequestArgs{ClerkSCIMDirectoryID: data.DirectoryID.ValueString(), IDs: order}})
	return acknowledged(result != nil && result.Success, err, "group role order")
}

func (r *scimGroupRoleMappingsResource) write(ctx context.Context, data *scimGroupRoleMappingsModel) error {
	if err := r.apply(ctx, data); err != nil {
		return err
	}
	groups, err := r.groups(ctx, *data)
	if err != nil {
		return err
	}
	planned := data.Mappings
	data.refresh(groups)
	if !slices.EqualFunc(planned, data.Mappings, func(a, b scimGroupRoleMapping) bool {
		return a.SCIMGroupID.Equal(b.SCIMGroupID) && a.Role.Equal(b.Role)
	}) {
		return errors.New("The API acknowledged the role mappings but returned a different set or order. Check that role mapping is turned on for the directory.")
	}
	return nil
}

func (r *scimGroupRoleMappingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data scimGroupRoleMappingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.write(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Unable to set group role mappings", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *scimGroupRoleMappingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data scimGroupRoleMappingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.write(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Unable to update group role mappings", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *scimGroupRoleMappingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data scimGroupRoleMappingsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	groups, err := r.groups(ctx, data)
	if errors.Is(err, errResourceNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read group role mappings", err.Error())
		return
	}
	data.refresh(groups)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *scimGroupRoleMappingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data scimGroupRoleMappingsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	groups, err := r.groups(ctx, data)
	if errors.Is(err, errResourceNotFound) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to remove group role mappings", err.Error())
		return
	}
	for _, group := range groups {
		if group.RoleKey == nil {
			continue
		}
		if err := r.setRole(ctx, data.WorkspaceSlug.ValueString(), group, nil); err != nil {
			resp.Diagnostics.AddError("Unable to remove group role mappings", err.Error())
			return
		}
	}
}

func (r *scimGroupRoleMappingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	slug, directory, found := strings.Cut(req.ID, "/")
	if !found || !slugPattern.MatchString(slug) || directory == "" || strings.Contains(directory, "/") {
		resp.Diagnostics.AddError("Invalid import identifier", "Use workspace-slug/directory-id.")
		return
	}
	data := scimGroupRoleMappingsModel{WorkspaceSlug: types.StringValue(slug), DirectoryID: types.StringValue(directory)}
	groups, err := r.groups(ctx, data)
	if errors.Is(err, errResourceNotFound) {
		resp.Diagnostics.AddError("Directory not found", "The directory has no provisioned groups in this workspace.")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to import group role mappings", err.Error())
		return
	}
	data.refresh(groups)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

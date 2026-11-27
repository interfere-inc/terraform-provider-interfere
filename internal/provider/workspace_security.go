package provider

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

var authMethods = []string{"enterprise_sso", "oauth_google", "oauth_github", "oauth_vercel", "passkey", "email_code", "password"}

var alwaysEnabledAuthMethods = []string{"email_code", "passkey"}

type alwaysEnabledAuthMethodsValidator struct{}

func (alwaysEnabledAuthMethodsValidator) Description(context.Context) string {
	return "Must include email_code and passkey."
}

func (v alwaysEnabledAuthMethodsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (alwaysEnabledAuthMethodsValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var methods []types.String
	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &methods, false)...)
	for _, method := range methods {
		if method.IsUnknown() {
			return
		}
	}
	for _, required := range alwaysEnabledAuthMethods {
		if !slices.ContainsFunc(methods, func(method types.String) bool { return method.ValueString() == required }) {
			resp.Diagnostics.AddAttributeError(req.Path, "Missing sign-in method", fmt.Sprintf("Email code and passkey sign-in are always enabled; include %q.", required))
		}
	}
}

type supportAccessExpiryModifier struct{}

func (supportAccessExpiryModifier) Description(context.Context) string {
	return "Keeps the current expiry when support access is not configured; otherwise an omitted expiry means no expiry."
}

func (m supportAccessExpiryModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (supportAccessExpiryModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if !req.ConfigValue.IsNull() {
		return
	}
	var granted types.Bool
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("impersonation_granted"), &granted)...)
	if !granted.IsNull() {
		resp.PlanValue = types.Int64Null()
		return
	}
	if !req.State.Raw.IsNull() {
		resp.PlanValue = req.StateValue
	}
}

func known[T interface {
	IsNull() bool
	IsUnknown() bool
}](value T) bool {
	return !value.IsNull() && !value.IsUnknown()
}

func remoteAuthMethods(remote *sdk.QueryWorkspacesCurrentResponse) []string {
	if remote.EnabledAuthMethods == nil {
		return nil
	}
	methods := make([]string, 0, len(remote.EnabledAuthMethods))
	for _, method := range remote.EnabledAuthMethods {
		methods = append(methods, string(method))
	}
	return methods
}

func authMethodSet(methods []string) (types.Set, diag.Diagnostics) {
	if methods == nil {
		return types.SetNull(types.StringType), nil
	}
	values := make([]attr.Value, 0, len(methods))
	for _, method := range methods {
		values = append(values, types.StringValue(method))
	}
	return types.SetValue(types.StringType, values)
}

func remoteSupportExpiry(remote *sdk.QueryWorkspacesCurrentResponse) types.Int64 {
	if remote.ImpersonationExpiresAt == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*remote.ImpersonationExpiresAt))
}

func setWorkspaceSecurity(data *workspaceModel, remote *sdk.QueryWorkspacesCurrentResponse) diag.Diagnostics {
	data.AllowMemberWorkspaceCreation = types.BoolValue(remote.AllowMemberWorkspaceCreation)
	data.ImpersonationGranted = types.BoolValue(remote.ImpersonationGranted)
	data.ImpersonationExpiresAt = remoteSupportExpiry(remote)
	methods, diags := authMethodSet(remoteAuthMethods(remote))
	data.EnabledAuthMethods = methods
	return diags
}

func acknowledged(success bool, err error, setting string) error {
	if err != nil {
		return errors.New(apiError(err))
	}
	if !success {
		return fmt.Errorf("The API did not acknowledge the %s update. Terraform state was preserved.", setting)
	}
	return nil
}

func (r *workspaceResource) applySecurity(ctx context.Context, data *workspaceModel, remote *sdk.QueryWorkspacesCurrentResponse) error {
	slug := data.WorkspaceSlug.ValueString()
	if known(data.AllowMemberWorkspaceCreation) {
		if data.AllowMemberWorkspaceCreation.ValueBool() != remote.AllowMemberWorkspaceCreation {
			result, err := r.client.Workspaces.UpdateMemberWorkspaceCreation(ctx, &sdk.UpdateMemberWorkspaceCreationWorkspacesRequest{
				WorkspaceSlug: slug,
				Args:          &sdk.UpdateMemberWorkspaceCreationWorkspacesRequestArgs{AllowMemberWorkspaceCreation: data.AllowMemberWorkspaceCreation.ValueBool()},
			})
			if err := acknowledged(result != nil && result.Success, err, "member workspace creation"); err != nil {
				return err
			}
		}
	} else {
		data.AllowMemberWorkspaceCreation = types.BoolValue(remote.AllowMemberWorkspaceCreation)
	}

	granted := remote.ImpersonationGranted
	if known(data.ImpersonationGranted) {
		granted = data.ImpersonationGranted.ValueBool()
	}
	expiry := remoteSupportExpiry(remote)
	if !data.ImpersonationExpiresAt.IsUnknown() {
		expiry = data.ImpersonationExpiresAt
	}
	if granted != remote.ImpersonationGranted || !expiry.Equal(remoteSupportExpiry(remote)) {
		args := &sdk.SetAccessSupportRequestArgs{}
		args.SetImpersonationGranted(granted)
		var expiresAt *int
		if !expiry.IsNull() {
			value := int(expiry.ValueInt64())
			expiresAt = &value
		}
		args.SetImpersonationExpiresAt(expiresAt)
		result, err := r.client.Support.SetAccess(ctx, &sdk.SetAccessSupportRequest{WorkspaceSlug: slug, Args: args})
		if err := acknowledged(result != nil && result.Success, err, "support access"); err != nil {
			return err
		}
	}
	data.ImpersonationGranted = types.BoolValue(granted)
	data.ImpersonationExpiresAt = expiry

	current := remoteAuthMethods(remote)
	if !known(data.EnabledAuthMethods) {
		methods, diags := authMethodSet(current)
		if diags.HasError() {
			return errors.New("The API returned invalid sign-in methods. Terraform state was preserved.")
		}
		data.EnabledAuthMethods = methods
		return nil
	}
	var desired []string
	if diags := data.EnabledAuthMethods.ElementsAs(ctx, &desired, false); diags.HasError() {
		return errors.New("Invalid sign-in methods in the plan.")
	}
	if current != nil && len(current) == len(desired) && !slices.ContainsFunc(desired, func(method string) bool { return !slices.Contains(current, method) }) {
		return nil
	}
	methods := make([]sdk.UpdateAuthMethodsWorkspacesRequestArgsMethodsItem, 0, len(desired))
	for _, method := range authMethods {
		if slices.Contains(desired, method) {
			methods = append(methods, sdk.UpdateAuthMethodsWorkspacesRequestArgsMethodsItem(method))
		}
	}
	result, err := r.client.Workspaces.UpdateAuthMethods(ctx, &sdk.UpdateAuthMethodsWorkspacesRequest{
		WorkspaceSlug: slug,
		Args:          &sdk.UpdateAuthMethodsWorkspacesRequestArgs{Methods: methods},
	})
	return acknowledged(result != nil && result.Success, err, "sign-in methods")
}

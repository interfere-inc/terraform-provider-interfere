package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_verified_domain"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

var pathEnrollmentMode = path.Root("enrollment_mode")

var enrollmentModes = []string{"manual_invitation", "automatic_invitation", "automatic_suggestion"}

type verifiedDomainResource struct{ client *client.Client }
type verifiedDomainModel struct {
	ID                      types.String `tfsdk:"id"`
	WorkspaceSlug           types.String `tfsdk:"workspace_slug"`
	Name                    types.String `tfsdk:"name"`
	EnrollmentMode          types.String `tfsdk:"enrollment_mode"`
	VerificationStatus      types.String `tfsdk:"verification_status"`
	VerifiedAt              types.Int64  `tfsdk:"verified_at"`
	AffiliationEmailAddress types.String `tfsdk:"affiliation_email_address"`
	ClerkDomainID           types.String `tfsdk:"clerk_domain_id"`
}

func NewVerifiedDomainResource() resource.Resource { return &verifiedDomainResource{} }
func (r *verifiedDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_verified_domain"
}

func (r *verifiedDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verifiedDomainResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_verified_domain.VerifiedDomainResourceSchema(ctx)
	resp.Schema.Description = "Manage a workspace email domain and how people with an address at it join the workspace. Creation returns while verification is pending; a workspace admin verifies the domain in the dashboard, and verification_status refreshes on the next plan. Automatic enrollment modes require a verified domain, so create the domain with manual_invitation and change enrollment_mode after verification. Domains managed by an enterprise SSO connection cannot be managed here. Import with workspace-slug/domain-name."
	attributes := resp.Schema.Attributes
	flattenAttributes(attributes, "args")
	identifier := attributes["id"].(schema.StringAttribute)
	attributes["id"] = schema.StringAttribute{Computed: true, Description: identifier.Description, MarkdownDescription: identifier.MarkdownDescription, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}}
	attributes["workspace_slug"] = requiredString(attributes["workspace_slug"], stringplanmodifier.RequiresReplace())
	name := requiredString(attributes["name"], stringplanmodifier.RequiresReplace())
	name.Validators = []validator.String{stringvalidator.LengthAtMost(253), stringvalidator.RegexMatches(trackingHostnamePattern, "Use a lowercase domain name without a scheme, path or trailing dot.")}
	attributes["name"] = name
	attributes["enrollment_mode"] = schema.StringAttribute{
		Optional: true, Computed: true,
		Default:     stringdefault.StaticString("manual_invitation"),
		Description: "How people with an email address at this domain join the workspace: manual_invitation (an admin invites them), automatic_invitation (they are invited automatically) or automatic_suggestion (they are offered to request to join). Automatic modes require a verified domain. Defaults to manual_invitation.",
		Validators:  []validator.String{stringvalidator.OneOf(enrollmentModes...)},
	}
	attributes["verification_status"] = schema.StringAttribute{
		Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		Description: "Domain verification status: pending or verified. Verification happens outside Terraform and is read on refresh.",
	}
	attributes["verified_at"] = schema.Int64Attribute{
		Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		Description: "When the domain was verified, in Unix milliseconds. Null while verification is pending.",
	}
	attributes["affiliation_email_address"] = schema.StringAttribute{
		Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		Description: "Email address used to verify the domain. Null until verification starts.",
	}
	attributes["clerk_domain_id"] = schema.StringAttribute{
		Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		Description: "Identifier of the domain in the sign-in service.",
	}
}

func (r *verifiedDomainResource) list(ctx context.Context, slug string) (sdk.QueryDomainsForOrgResponse, error) {
	remote, err := r.client.Domains.List(ctx, &sdk.ListDomainsRequest{WorkspaceSlug: slug})
	if err != nil {
		return nil, errors.New(apiError(err))
	}
	return remote, nil
}

func (r *verifiedDomainResource) find(ctx context.Context, slug string, match func(*sdk.QueryDomainsForOrgResponseItem) bool) (*sdk.QueryDomainsForOrgResponseItem, error) {
	domains, err := r.list(ctx, slug)
	if err != nil {
		return nil, err
	}
	for _, domain := range domains {
		if domain == nil || domain.Type == nil || *domain.Type != sdk.QueryDomainsForOrgResponseItemTypeClerk || !match(domain) {
			continue
		}
		if domain.EnrollmentMode != nil && *domain.EnrollmentMode == "enterprise_sso" {
			return nil, errors.New("This domain is managed by an enterprise SSO connection and cannot be managed by Terraform.")
		}
		if domain.ClerkDomainID == nil || *domain.ClerkDomainID == "" {
			return nil, errors.New("The domain has no sign-in service identifier yet. Retry after it finishes registering; Terraform state was preserved.")
		}
		return domain, nil
	}
	return nil, errResourceNotFound
}

func (r *verifiedDomainResource) read(ctx context.Context, data verifiedDomainModel) (*sdk.QueryDomainsForOrgResponseItem, error) {
	remote, err := r.find(ctx, data.WorkspaceSlug.ValueString(), func(domain *sdk.QueryDomainsForOrgResponseItem) bool { return domain.ID == data.ID.ValueString() })
	if err != nil {
		return nil, err
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() && remote.Name != data.Name.ValueString() {
		return nil, errors.New("The domain identity now has a different name. Terraform state was preserved.")
	}
	return remote, nil
}

func (data *verifiedDomainModel) refresh(remote *sdk.QueryDomainsForOrgResponseItem) {
	data.ID = types.StringValue(remote.ID)
	data.Name = types.StringValue(remote.Name)
	data.EnrollmentMode = types.StringValue("manual_invitation")
	if remote.EnrollmentMode != nil {
		data.EnrollmentMode = types.StringValue(*remote.EnrollmentMode)
	}
	data.VerificationStatus = types.StringValue("pending")
	data.VerifiedAt = types.Int64Null()
	if remote.VerifiedAt != nil {
		data.VerificationStatus = types.StringValue("verified")
		data.VerifiedAt = types.Int64Value(int64(*remote.VerifiedAt))
	}
	data.AffiliationEmailAddress = types.StringPointerValue(remote.AffiliationEmailAddress)
	data.ClerkDomainID = types.StringPointerValue(remote.ClerkDomainID)
}

func automaticEnrollment(mode string) bool {
	return mode == "automatic_invitation" || mode == "automatic_suggestion"
}

func (r *verifiedDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data verifiedDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if automaticEnrollment(data.EnrollmentMode.ValueString()) {
		resp.Diagnostics.AddAttributeError(pathEnrollmentMode, "Domain verification required", "A new domain starts unverified. Create it with enrollment_mode = \"manual_invitation\", verify it in the dashboard, then change enrollment_mode.")
		return
	}
	id := uuid.NewString()
	result, err := r.client.Workspaces.AddDomain(ctx, &sdk.AddDomainWorkspacesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.AddDomainWorkspacesRequestArgs{ID: id, Name: data.Name.ValueString()}})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create domain", fmt.Sprintf("%s If the domain was added anyway, import it with %s/%s.", apiError(err), data.WorkspaceSlug.ValueString(), data.Name.ValueString()))
		return
	}
	if result == nil || !result.Success {
		resp.Diagnostics.AddError("Creation not acknowledged", "The API did not confirm domain creation.")
		return
	}
	data.ID = types.StringValue(id)
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read created domain", fmt.Sprintf("%s Import it with %s/%s.", err.Error(), data.WorkspaceSlug.ValueString(), data.Name.ValueString()))
		return
	}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *verifiedDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data verifiedDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if errors.Is(err, errResourceNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read domain", err.Error())
		return
	}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *verifiedDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data verifiedDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read domain", err.Error())
		return
	}
	mode := data.EnrollmentMode.ValueString()
	if automaticEnrollment(mode) && remote.VerifiedAt == nil {
		resp.Diagnostics.AddAttributeError(pathEnrollmentMode, "Domain verification required", "Verify the domain in the dashboard before choosing an automatic enrollment mode. Terraform state was preserved.")
		return
	}
	if remote.EnrollmentMode == nil || *remote.EnrollmentMode != mode {
		enrollment, err := sdk.NewUpdateDomainEnrollmentModeWorkspacesRequestArgsEnrollmentModeFromString(mode)
		if err != nil {
			resp.Diagnostics.AddAttributeError(pathEnrollmentMode, "Invalid enrollment mode", err.Error())
			return
		}
		result, err := r.client.Workspaces.UpdateDomainEnrollmentMode(ctx, &sdk.UpdateDomainEnrollmentModeWorkspacesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.UpdateDomainEnrollmentModeWorkspacesRequestArgs{ClerkDomainID: *remote.ClerkDomainID, EnrollmentMode: enrollment}})
		if err := acknowledged(result != nil && result.Success, err, "enrollment mode"); err != nil {
			resp.Diagnostics.AddError("Unable to update domain", err.Error())
			return
		}
		remote, err = r.read(ctx, data)
		if err != nil {
			resp.Diagnostics.AddError("Unable to read updated domain", err.Error())
			return
		}
		if remote.EnrollmentMode == nil || *remote.EnrollmentMode != mode {
			resp.Diagnostics.AddError("Enrollment mode not applied", "The API acknowledged the update but the domain kept its previous enrollment mode. Terraform state was preserved.")
			return
		}
	}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *verifiedDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data verifiedDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.read(ctx, data)
	if errors.Is(err, errResourceNotFound) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete domain", err.Error())
		return
	}
	result, err := r.client.Workspaces.RemoveDomain(ctx, &sdk.RemoveDomainWorkspacesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.RemoveDomainWorkspacesRequestArgs{ClerkDomainID: *remote.ClerkDomainID}})
	if err := acknowledged(result != nil && result.Success, err, "domain removal"); err != nil {
		resp.Diagnostics.AddError("Unable to delete domain", err.Error())
	}
}

func (r *verifiedDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	slug, name, found := strings.Cut(req.ID, "/")
	if !found || !slugPattern.MatchString(slug) || !trackingHostnamePattern.MatchString(name) {
		resp.Diagnostics.AddError("Invalid import identifier", "Use workspace-slug/domain-name, such as acme/example.com.")
		return
	}
	remote, err := r.find(ctx, slug, func(domain *sdk.QueryDomainsForOrgResponseItem) bool { return domain.Name == name })
	if errors.Is(err, errResourceNotFound) {
		resp.Diagnostics.AddError("Domain not found", "No active verified domain with this name exists in the workspace.")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to import domain", err.Error())
		return
	}
	data := verifiedDomainModel{WorkspaceSlug: types.StringValue(slug)}
	data.refresh(remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

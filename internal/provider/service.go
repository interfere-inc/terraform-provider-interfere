package provider

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
)

type serviceResource struct{ client *client.Client }
type serviceModel struct {
	ID              types.String `tfsdk:"id"`
	WorkspaceSlug   types.String `tfsdk:"workspace_slug"`
	Name            types.String `tfsdk:"name"`
	Namespace       types.String `tfsdk:"namespace"`
	Type            types.String `tfsdk:"type"`
	WebsiteSettings types.Object `tfsdk:"website_settings"`
}

func NewServiceResource() resource.Resource { return &serviceResource{} }

func (r *serviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *serviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *serviceResource) read(ctx context.Context, data *serviceModel) error {
	remote, err := r.client.Services.Get(ctx, &sdk.GetServicesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.GetServicesRequestArgs{ServiceID: data.ID.ValueString()}})
	if err != nil {
		return errors.New(apiError(err))
	}
	if remote == nil || remote.ID != data.ID.ValueString() || remote.Name == "" || !validServiceType(remote.Type) {
		return errors.New("The API returned no service or an unexpected service identity. Verify workspace access; Terraform state was preserved.")
	}
	settings, err := serviceSettingsValue(ctx, remote.WebsiteSettings)
	if err != nil {
		return err
	}
	data.Name = types.StringValue(remote.Name)
	data.Namespace = types.StringPointerValue(remote.Namespace)
	data.Type = types.StringValue(remote.Type)
	data.WebsiteSettings = settings
	return nil
}

func (r *serviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data serviceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.read(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Unable to read service", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceResource) apply(ctx context.Context, data *serviceModel) error {
	current := *data
	if err := r.read(ctx, &current); err != nil {
		return err
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() && !data.Type.Equal(current.Type) {
		result, err := r.client.Services.SetType(ctx, &sdk.SetTypeServicesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.SetTypeServicesRequestArgs{ServiceID: data.ID.ValueString(), Type: sdk.SetTypeServicesRequestArgsType(data.Type.ValueString())}})
		if err != nil {
			return errors.New(apiError(err))
		}
		if result == nil || !result.Success {
			return errors.New("The API did not acknowledge the service type update. Terraform state was preserved.")
		}
	}
	if !data.WebsiteSettings.IsNull() && !data.WebsiteSettings.IsUnknown() && !data.WebsiteSettings.Equal(current.WebsiteSettings) {
		input, err := serviceSettingsRequest(ctx, *data)
		if err != nil {
			return err
		}
		result, err := r.client.Services.SetWebsiteSettings(ctx, input)
		if err != nil {
			return errors.New(apiError(err))
		}
		if result == nil || !result.Success {
			return errors.New("The API did not acknowledge the service website settings update. Terraform state was preserved.")
		}
	}
	return r.read(ctx, data)
}

func (r *serviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data serviceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.apply(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Unable to manage service", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data serviceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.apply(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Unable to update service", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}

func (r *serviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || !slugPattern.MatchString(parts[0]) || !uuidPattern.MatchString(parts[1]) {
		resp.Diagnostics.AddError("Invalid service import ID", "Use workspace-slug/service-uuid.")
		return
	}
	data := serviceModel{WorkspaceSlug: types.StringValue(parts[0]), ID: types.StringValue(parts[1])}
	if err := r.read(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Unable to import service", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func validServiceType(value string) bool {
	return value == "unclassified" || value == "other" || value == "backend" || value == "website"
}

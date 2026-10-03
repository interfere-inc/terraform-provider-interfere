package provider

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/client"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/option"
)

type interfereProvider struct{ version string }

type providerModel struct {
	Token   types.String `tfsdk:"token"`
	BaseURL types.String `tfsdk:"base_url"`
	Headers types.Map    `tfsdk:"headers"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider { return &interfereProvider{version: version} }
}

func (p *interfereProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "interfere"
	resp.Version = p.version
}

func (p *interfereProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage workspace settings, surfaces, integration mappings, tracking domains, and API keys using a workspace API key, session, or delegated OAuth access token.",
		Attributes: map[string]schema.Attribute{
			"headers":  schema.MapAttribute{Optional: true, Sensitive: true, ElementType: types.StringType, Description: "Additional HTTP headers for an authenticated API proxy, such as Cloudflare Access. Values must be known before planning."},
			"token":    schema.StringAttribute{Optional: true, Sensitive: true, Description: "Workspace API key, session, or delegated OAuth access token. Defaults to INTERFERE_TOKEN. Grant the workspace-basics, surface, or workspace-auth permissions required by the configured resources. Release-only keys are not supported."},
			"base_url": schema.StringAttribute{Optional: true, Description: "API URL. Defaults to https://api.interfere.com. Plain HTTP is supported only for local development on loopback addresses."},
		},
	}
}

func (p *interfereProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.Token.IsUnknown() || data.BaseURL.IsUnknown() || data.Headers.IsUnknown() {
		resp.Diagnostics.AddError("Unknown provider configuration", "The token, base_url, and headers must be known before managing resources.")
		return
	}
	headers := http.Header{}
	if !data.Headers.IsNull() {
		var configured map[string]string
		resp.Diagnostics.Append(data.Headers.ElementsAs(ctx, &configured, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for name, value := range configured {
			headers.Set(name, value)
		}
	}
	token := data.Token.ValueString()
	if data.Token.IsNull() {
		token = os.Getenv("INTERFERE_TOKEN")
	}
	if strings.TrimSpace(token) == "" {
		resp.Diagnostics.AddError("Missing access token", "Set INTERFERE_TOKEN or configure token with a workspace API key, session, or delegated OAuth access token.")
		return
	}
	baseURL := data.BaseURL.ValueString()
	if data.BaseURL.IsNull() {
		baseURL = "https://api.interfere.com"
	}
	endpoint, err := url.Parse(baseURL)
	if err != nil || endpoint.Host == "" || endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != "" {
		resp.Diagnostics.AddError("Invalid API URL", "base_url must be an absolute API URL without credentials, a query, or a fragment.")
		return
	}
	local := endpoint.Hostname() == "localhost" || endpoint.Hostname() == "127.0.0.1" || endpoint.Hostname() == "::1"
	if endpoint.Scheme != "https" && !(endpoint.Scheme == "http" && local) {
		resp.Diagnostics.AddError("Invalid API URL", "Use HTTPS, or HTTP with a loopback address for local development.")
		return
	}
	configured := client.NewClient(
		option.WithHTTPHeader(headers),
		option.WithBaseURL(strings.TrimRight(baseURL, "/")),
		option.WithToken(token),
		option.WithoutRetries(),
		option.WithHTTPClient(&http.Client{
			Timeout:       30 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse },
		}),
	)
	resp.ResourceData = configured
	resp.DataSourceData = configured
}

func (p *interfereProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewSurfaceResource, NewWorkspaceResource, NewPrivateKeyResource, NewPublicKeyResource, NewSurfaceRepositoryResource, NewSurfaceDestinationResource, NewTrackingDomainResource}
}

func (p *interfereProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewWorkspaceDataSource, NewSurfaceDataSource, NewIntegrationDataSource}
}

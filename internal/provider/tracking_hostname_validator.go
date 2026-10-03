package provider

import (
	"context"
	"net"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var trackingHostnamePattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+$`)

type trackingHostnameValidator struct{}

func (trackingHostnameValidator) Description(context.Context) string {
	return "Use a lowercase DNS hostname without a scheme, path, trailing dot, IP address, or www prefix."
}

func (v trackingHostnameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v trackingHostnameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	if len(value) > 253 || !trackingHostnamePattern.MatchString(value) || strings.HasPrefix(value, "www.") || net.ParseIP(value) != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid tracking hostname", v.Description(ctx))
	}
}

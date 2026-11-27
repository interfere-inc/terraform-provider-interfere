package provider

import (
	"context"
	"fmt"
	"net/netip"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var asnPattern = regexp.MustCompile(`^AS([1-9][0-9]{0,9})$`)

func canonicalNetworkRule(value string) (string, bool) {
	if match := asnPattern.FindStringSubmatch(value); match != nil {
		number, err := strconv.ParseUint(match[1], 10, 64)
		return value, err == nil && number <= 4_294_967_294
	}
	prefix, err := netip.ParsePrefix(value)
	if err != nil || prefix.Addr().Is4In6() || prefix.Addr().Zone() != "" {
		if address, addressErr := netip.ParseAddr(value); addressErr == nil && !address.Is4In6() && address.Zone() == "" {
			return netip.PrefixFrom(address, address.BitLen()).String(), false
		}
		return "", false
	}
	canonical := prefix.Masked().String()
	return canonical, canonical == value
}

type networkRuleValueValidator struct{}

func (networkRuleValueValidator) Description(context.Context) string {
	return "Use a canonical CIDR range with its prefix length and host bits cleared, such as 203.0.113.0/24, 198.51.100.7/32 or 2001:db8::/32, or an ASN written as AS followed by its number, such as AS13335."
}

func (v networkRuleValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v networkRuleValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	canonical, ok := canonicalNetworkRule(value)
	if ok {
		return
	}
	detail := v.Description(ctx)
	if canonical != "" {
		detail = fmt.Sprintf("Write %q as %q. %s", value, canonical, detail)
	}
	resp.Diagnostics.AddAttributeError(req.Path, "Invalid network rule", detail)
}

type networkRuleLabelValidator struct{}

func (networkRuleLabelValidator) Description(context.Context) string {
	return "Labels must not start or end with whitespace."
}

func (v networkRuleLabelValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v networkRuleLabelValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	if strings.TrimSpace(value) != value {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid network rule label", v.Description(ctx))
	}
}

type uniqueNetworkRuleValues struct{}

func (uniqueNetworkRuleValues) Description(context.Context) string {
	return "Each network may appear in only one rule."
}

func (v uniqueNetworkRuleValues) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v uniqueNetworkRuleValues) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	seen := map[string]bool{}
	for _, element := range req.ConfigValue.Elements() {
		rule, ok := element.(types.Object)
		if !ok || rule.IsUnknown() {
			continue
		}
		value, ok := rule.Attributes()["value"].(types.String)
		if !ok || value.IsNull() || value.IsUnknown() {
			continue
		}
		if seen[value.ValueString()] {
			resp.Diagnostics.AddAttributeError(req.Path, "Duplicate network rule", fmt.Sprintf("%s appears in more than one rule. %s", value.ValueString(), v.Description(ctx)))
		}
		seen[value.ValueString()] = true
	}
}

package provider

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/resource_service"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

type serviceWebsiteSettings struct {
	AnonymousReplayPolicy string          `json:"anonymousReplayPolicy"`
	AnonymousUserTracking bool            `json:"anonymousUserTracking"`
	EnvironmentFiltering  bool            `json:"environmentFiltering"`
	Plugins               map[string]bool `json:"plugins,omitempty"`
}

var servicePluginNames = map[string]string{"device": "device", "errors": "errors", "logs": "logs", "page_events": "pageEvents", "rage_click": "rageClick", "replay": "replay"}

func serviceSettingTypes(ctx context.Context) map[string]attr.Type {
	result := resource_service.WebsiteSettingsValue{}.AttributeTypes(ctx)
	result["plugins"] = types.ObjectType{AttrTypes: resource_service.PluginsValue{}.AttributeTypes(ctx)}
	return result
}

func serviceSettingsValue(ctx context.Context, remote any) (types.Object, error) {
	kinds := serviceSettingTypes(ctx)
	value := types.ObjectNull(kinds)
	payload, err := json.Marshal(remote)
	if err != nil {
		return value, errors.New("The API returned invalid website settings.")
	}
	var settings *serviceWebsiteSettings
	if err := json.Unmarshal(payload, &settings); err != nil {
		return value, errors.New("The API returned invalid website settings.")
	}
	if settings == nil {
		return value, nil
	}
	if settings.AnonymousReplayPolicy != "always" && settings.AnonymousReplayPolicy != "on_error" && settings.AnonymousReplayPolicy != "never" {
		return value, errors.New("The API returned an invalid anonymous replay policy.")
	}
	plugins := types.ObjectNull(resource_service.PluginsValue{}.AttributeTypes(ctx))
	if settings.Plugins != nil {
		values := map[string]attr.Value{}
		for field, name := range servicePluginNames {
			values[field] = types.BoolNull()
			if enabled, ok := settings.Plugins[name]; ok {
				values[field] = types.BoolValue(enabled)
			}
		}
		converted, diagnostics := types.ObjectValue(resource_service.PluginsValue{}.AttributeTypes(ctx), values)
		if diagnostics.HasError() {
			return value, errors.New("The API returned invalid website plugin settings.")
		}
		plugins = converted
	}
	value, diagnostics := types.ObjectValue(kinds, map[string]attr.Value{
		"anonymous_replay_policy": types.StringValue(settings.AnonymousReplayPolicy),
		"anonymous_user_tracking": types.BoolValue(settings.AnonymousUserTracking),
		"environment_filtering":   types.BoolValue(settings.EnvironmentFiltering),
		"plugins":                 plugins,
	})
	if diagnostics.HasError() {
		return value, errors.New("The API returned invalid website settings.")
	}
	return value, nil
}

func serviceSettingsRequest(_ context.Context, data serviceModel) (*sdk.SetWebsiteSettingsServicesRequest, error) {
	values := data.WebsiteSettings.Attributes()
	settings := serviceWebsiteSettings{
		AnonymousReplayPolicy: values["anonymous_replay_policy"].(types.String).ValueString(),
		AnonymousUserTracking: values["anonymous_user_tracking"].(types.Bool).ValueBool(),
		EnvironmentFiltering:  values["environment_filtering"].(types.Bool).ValueBool(),
	}
	plugins := values["plugins"].(types.Object)
	if !plugins.IsNull() {
		settings.Plugins = map[string]bool{}
		for field, name := range servicePluginNames {
			value := plugins.Attributes()[field].(types.Bool)
			if !value.IsNull() {
				settings.Plugins[name] = value.ValueBool()
			}
		}
	}
	website := map[string]any{"anonymousReplayPolicy": settings.AnonymousReplayPolicy, "anonymousUserTracking": settings.AnonymousUserTracking, "environmentFiltering": settings.EnvironmentFiltering}
	if !plugins.IsNull() {
		website["plugins"] = settings.Plugins
	}
	payload, err := json.Marshal(map[string]any{"args": map[string]any{"serviceId": data.ID.ValueString(), "websiteSettings": website}})
	if err != nil {
		return nil, errors.New("Unable to encode service settings.")
	}
	request := &sdk.SetWebsiteSettingsServicesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString()}
	if err := json.Unmarshal(payload, request); err != nil {
		return nil, errors.New("Unable to encode the service settings request.")
	}
	request.WorkspaceSlug = data.WorkspaceSlug.ValueString()
	return request, nil
}

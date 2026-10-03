package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

func (r *surfaceResource) applySettings(ctx context.Context, data *surfaceModel) error {
	remote, err := r.read(ctx, *data)
	if err != nil {
		return err
	}
	if remote.DeletedAt != nil {
		return errors.New("The surface was deleted. Refresh the plan before retrying.")
	}
	if !data.AnonymousUserTracking.IsNull() && !data.AnonymousUserTracking.IsUnknown() && data.AnonymousUserTracking.ValueBool() != remote.AnonymousUserTracking {
		result, err := r.client.Surfaces.SetAnonymousUserTracking(ctx, &sdk.SetAnonymousUserTrackingSurfacesRequest{WorkspaceSlug: data.WorkspaceSlug.ValueString(), Args: &sdk.SetAnonymousUserTrackingSurfacesRequestArgs{SurfaceSlug: data.Slug.ValueString(), Enabled: data.AnonymousUserTracking.ValueBool()}})
		if err != nil {
			return errors.New(apiError(err))
		}
		if result == nil || !result.Success {
			return errors.New("The API did not acknowledge the settings update.")
		}
	} else {
		data.AnonymousUserTracking = types.BoolValue(remote.AnonymousUserTracking)
	}
	return nil
}

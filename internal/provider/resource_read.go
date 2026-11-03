package provider

import (
	"errors"
	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk"
)

var errResourceNotFound = errors.New("Resource no longer exists in this workspace")

func resourceReadError(err error) error {
	var missing *sdk.NotFoundError
	if errors.As(err, &missing) {
		if body, ok := missing.Body.(map[string]interface{}); ok && body["code"] == "WORKSPACE_RESOURCE_NOT_FOUND" {
			return errResourceNotFound
		}
	}
	return errors.New(apiError(err))
}

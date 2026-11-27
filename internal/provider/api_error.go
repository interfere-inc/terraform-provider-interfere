package provider

import (
	"errors"
	"fmt"

	"github.com/interfere-inc/terraform-provider-interfere/internal/sdk/core"
)

func apiError(err error) string {
	var response *core.APIError
	if errors.As(err, &response) {
		return fmt.Sprintf("The Interfere API returned HTTP %d. Check workspace access and the resource configuration. The response body is omitted because it may contain credentials.", response.StatusCode)
	}
	return "The API request failed or returned an invalid response. Check connectivity and credentials. For an uncertain create result, retry with the same idempotency_key."
}

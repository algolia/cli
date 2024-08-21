package provisionning

import (
	"fmt"
	"net/http"
)

// GetAvailableHostingRegions returns the list of available regions.
func (c *Client) GetAvailableHostingRegions(params HostingRegionsRequest) (*HostingRegionResponse, error) {
	var res *HostingRegionResponse

	urlParams := make(map[string]string)
	if params.PlanName != "" {
		urlParams["plan_name"] = params.PlanName
	}
	if params.ParentApplicationID != "" {
		urlParams["parent_application_id"] = params.ParentApplicationID
	}

	err := c.request(&res, http.MethodGet, "regions", nil, urlParams)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

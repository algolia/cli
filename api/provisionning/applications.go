package provisionning

import (
	"fmt"
	"net/http"
)

// ListApplications returns the list of applications.
func (c *Client) ListApplications() ([]Application, error) {
	var res []Application
	err := c.request(&res, http.MethodGet, "applications", nil, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetApplication returns the application with the given ID.
func (c *Client) GetApplication(id string) (*Application, error) {
	var res Application
	path := fmt.Sprintf("applications/%s", id)
	err := c.request(&res, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateFreeApplication creates a new free application on the requested region.
func (c *Client) CreateFreeApplication(params FreeApplicationCreationRequest) (*Application, error) {
	// Region is required
	if params.Region == "" {
		return nil, fmt.Errorf("region is required")
	}

	req := map[string]string{
		"region_code": params.Region,
		"name":        params.Name,
		"description": params.Description,
	}

	var res Application
	err := c.request(&res, http.MethodPost, "applications", nil, req)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetApplicationCreationStatus returns the status of the application creation.
func (c *Client) GetApplicationCreationStatus(id string) (string, error) {
	var res struct {
		Status string `json:"status"`
	}
	path := fmt.Sprintf("application/%s/status", id)
	err := c.request(&res, http.MethodGet, path, nil, nil)
	if err != nil {
		return "", err
	}

	return res.Status, nil
}

// AddApplicationMember adds a member to the application.
func (c *Client) AddApplicationMember(applicationID, email string) error {
	req := map[string]string{
		"email": email,
	}

	path := fmt.Sprintf("applications/%s/team-members", applicationID)
	err := c.request(nil, http.MethodPost, path, nil, req)
	if err != nil {
		return err
	}

	return nil
}

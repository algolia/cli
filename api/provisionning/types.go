package provisionning

import (
	"fmt"
)

type ErrResponse struct {
	Errors []struct {
		Status string `json:"status"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

func (e ErrResponse) Error() string {
	if e.Errors == nil {
		return ""
	}
	if len(e.Errors) == 0 {
		return ""
	}
	return fmt.Sprintf("[%s] %s: %s", e.Errors[0].Status, e.Errors[0].Title, e.Errors[0].Detail)
}

type PaginatedResponse struct {
	Meta struct {
		TotalCount int `jsonapi:"total_count"`
		PerPage    int `jsonapi:"per_page"`
		Page       int `jsonapi:"page"`
		TotalPages int `jsonapi:"total_pages"`
	} `jsonapi:"meta"`
}

type Application struct {
	ID          string    `jsonapi:"primary,application"`
	Name        string    `jsonapi:"attr,name"`
	Description string    `jsonapi:"attr,description"`
	Permissions []string  `jsonapi:"attr,permissions"`
	Plan        Plan      `jsonapi:"attr,plan"`
	LogRegion   string    `jsonapi:"attr,log_region"`
	Clusters    []Cluster `jsonapi:"attr,clusters"`
	IsBlocked   bool      `jsonapi:"attr,is_blocked"`
	Status      string    `jsonapi:"attr,status"`
}

type PaginatedApplications struct {
	Applications []Application `jsonapi:"data"`
}

type Plan struct {
	Name    string `jsonapi:"attr,name"`
	Version int    `jsonapi:"attr,version"`
}

type Cluster struct {
	Name         string `jsonapi:"attr,name"`
	LocationName string `jsonapi:"attr,location_name"`
	LocationCode string `jsonapi:"attr,location_code"`
}

type TeamMember struct {
	ID        string `jsonapi:"primary,team_member"`
	Email     string `jsonapi:"attr,email"`
	ExpiresOn string `jsonapi:"attr,expires_on"`
}

type FreeApplicationCreationRequest struct {
	Region      string `json:"region"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type HostingRegionsRequest struct {
	PlanName            string `json:"plan_name"`
	ParentApplicationID string `json:"parent_application_id"`
}

type HostingRegion struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Tier int    `json:"tier"`
}

type HostingRegionResponse struct {
	Regions []HostingRegion `json:"regions"`
	Details string          `json:"details"`
}

type APIKey struct {
	ID                  string   `jsonapi:"primary,api_key"`
	Name                string   `jsonapi:"attr,name"`
	Value               string   `jsonapi:"attr,value"`
	ACL                 []string `jsonapi:"attr,acl"`
	Indexes             []string `jsonapi:"attr,indexes"`
	Referers            []string `jsonapi:"attr,referers"`
	Description         string   `jsonapi:"attr,description"`
	MaxHitsPerQuery     int      `jsonapi:"attr,max_hits_per_query"`
	MaxHitsPerIPPerHour int      `jsonapi:"attr,max_hits_per_ip_per_hour"`
	Validity            int64    `jsonapi:"attr,validity"`
}

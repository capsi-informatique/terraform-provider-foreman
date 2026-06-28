package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HanseMerkur/terraform-provider-utils/log"
)

const (
	OrganizationEndpointPrefix = "organizations"
)

// ForemanOrganization represents a Foreman organization
type ForemanOrganization struct {
	ForemanObject

	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	ParentId    int                  `json:"parent_id,omitempty"`
	Parameters  []ForemanKVParameter `json:"organization_parameters_attributes,omitempty"`
}

func (c *Client) CreateOrganization(ctx context.Context, o *ForemanOrganization) (*ForemanOrganization, error) {
	log.Tracef("foreman/api/organization.go#Create")

	reqEndpoint := fmt.Sprintf("/%s", OrganizationEndpointPrefix)

	jsonBytes, jsonEncErr := c.WrapJSON("organization", o)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("organizationJSONBytes: [%s]", jsonBytes)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodPost, reqEndpoint, bytes.NewBuffer(jsonBytes))
	if reqErr != nil {
		return nil, reqErr
	}

	var created ForemanOrganization
	if err := c.SendAndParse(req, &created); err != nil {
		return nil, err
	}

	log.Debugf("createdOrganization: [%+v]", created)
	return &created, nil
}

func (c *Client) ReadOrganization(ctx context.Context, id int) (*ForemanOrganization, error) {
	log.Tracef("foreman/api/organization.go#Read")

	reqEndpoint := fmt.Sprintf("/%s/%d", OrganizationEndpointPrefix, id)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodGet, reqEndpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	var read ForemanOrganization
	if err := c.SendAndParse(req, &read); err != nil {
		return nil, err
	}

	log.Debugf("readOrganization: [%+v]", read)
	return &read, nil
}

func (c *Client) UpdateOrganization(ctx context.Context, o *ForemanOrganization) (*ForemanOrganization, error) {
	log.Tracef("foreman/api/organization.go#Update")

	reqEndpoint := fmt.Sprintf("/%s/%d", OrganizationEndpointPrefix, o.Id)

	jsonBytes, jsonEncErr := c.WrapJSON("organization", o)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("organizationJSONBytes: [%s]", jsonBytes)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodPut, reqEndpoint, bytes.NewBuffer(jsonBytes))
	if reqErr != nil {
		return nil, reqErr
	}

	var updated ForemanOrganization
	if err := c.SendAndParse(req, &updated); err != nil {
		return nil, err
	}

	log.Debugf("updatedOrganization: [%+v]", updated)
	return &updated, nil
}

func (c *Client) DeleteOrganization(ctx context.Context, id int) error {
	log.Tracef("foreman/api/organization.go#Delete")

	reqEndpoint := fmt.Sprintf("/%s/%d", OrganizationEndpointPrefix, id)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodDelete, reqEndpoint, nil)
	if reqErr != nil {
		return reqErr
	}

	return c.SendAndParse(req, nil)
}

func (c *Client) QueryOrganization(ctx context.Context, o *ForemanOrganization) (QueryResponse, error) {
	log.Tracef("foreman/api/organization.go#Query")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", OrganizationEndpointPrefix)
	req, reqErr := c.NewRequestWithContext(ctx, http.MethodGet, reqEndpoint, nil)
	if reqErr != nil {
		return queryResponse, reqErr
	}

	reqQuery := req.URL.Query()
	reqQuery.Set("search", "name=\""+o.Name+"\"")
	req.URL.RawQuery = reqQuery.Encode()

	if err := c.SendAndParse(req, &queryResponse); err != nil {
		return queryResponse, err
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	results := []ForemanOrganization{}
	resultsBytes, jsonEncErr := json.Marshal(queryResponse.Results)
	if jsonEncErr != nil {
		return queryResponse, jsonEncErr
	}
	if err := json.Unmarshal(resultsBytes, &results); err != nil {
		return queryResponse, err
	}

	iArr := make([]interface{}, len(results))
	for idx, val := range results {
		iArr[idx] = val
	}
	queryResponse.Results = iArr

	return queryResponse, nil
}

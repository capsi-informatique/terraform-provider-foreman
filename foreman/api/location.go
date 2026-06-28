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
	LocationEndpointPrefix = "locations"
)

// ForemanLocation represents a Foreman location
type ForemanLocation struct {
	ForemanObject

	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	ParentId    int                  `json:"parent_id,omitempty"`
	Parameters  []ForemanKVParameter `json:"location_parameters_attributes,omitempty"`
}

func (c *Client) CreateLocation(ctx context.Context, l *ForemanLocation) (*ForemanLocation, error) {
	log.Tracef("foreman/api/location.go#Create")

	reqEndpoint := fmt.Sprintf("/%s", LocationEndpointPrefix)

	jsonBytes, jsonEncErr := c.WrapJSON("location", l)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("locationJSONBytes: [%s]", jsonBytes)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodPost, reqEndpoint, bytes.NewBuffer(jsonBytes))
	if reqErr != nil {
		return nil, reqErr
	}

	var created ForemanLocation
	if err := c.SendAndParse(req, &created); err != nil {
		return nil, err
	}

	log.Debugf("createdLocation: [%+v]", created)
	return &created, nil
}

func (c *Client) ReadLocation(ctx context.Context, id int) (*ForemanLocation, error) {
	log.Tracef("foreman/api/location.go#Read")

	reqEndpoint := fmt.Sprintf("/%s/%d", LocationEndpointPrefix, id)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodGet, reqEndpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	var read ForemanLocation
	if err := c.SendAndParse(req, &read); err != nil {
		return nil, err
	}

	log.Debugf("readLocation: [%+v]", read)
	return &read, nil
}

func (c *Client) UpdateLocation(ctx context.Context, l *ForemanLocation) (*ForemanLocation, error) {
	log.Tracef("foreman/api/location.go#Update")

	reqEndpoint := fmt.Sprintf("/%s/%d", LocationEndpointPrefix, l.Id)

	jsonBytes, jsonEncErr := c.WrapJSON("location", l)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("locationJSONBytes: [%s]", jsonBytes)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodPut, reqEndpoint, bytes.NewBuffer(jsonBytes))
	if reqErr != nil {
		return nil, reqErr
	}

	var updated ForemanLocation
	if err := c.SendAndParse(req, &updated); err != nil {
		return nil, err
	}

	log.Debugf("updatedLocation: [%+v]", updated)
	return &updated, nil
}

func (c *Client) DeleteLocation(ctx context.Context, id int) error {
	log.Tracef("foreman/api/location.go#Delete")

	reqEndpoint := fmt.Sprintf("/%s/%d", LocationEndpointPrefix, id)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodDelete, reqEndpoint, nil)
	if reqErr != nil {
		return reqErr
	}

	return c.SendAndParse(req, nil)
}

func (c *Client) QueryLocation(ctx context.Context, l *ForemanLocation) (QueryResponse, error) {
	log.Tracef("foreman/api/location.go#Query")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", LocationEndpointPrefix)
	req, reqErr := c.NewRequestWithContext(ctx, http.MethodGet, reqEndpoint, nil)
	if reqErr != nil {
		return queryResponse, reqErr
	}

	reqQuery := req.URL.Query()
	reqQuery.Set("search", "name=\""+l.Name+"\"")
	req.URL.RawQuery = reqQuery.Encode()

	if err := c.SendAndParse(req, &queryResponse); err != nil {
		return queryResponse, err
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	results := []ForemanLocation{}
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

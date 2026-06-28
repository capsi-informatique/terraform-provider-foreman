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
	ConfigGroupEndpointPrefix = "config_groups"
)

// ForemanConfigGroup represents a Foreman config group (groups of Puppet classes)
type ForemanConfigGroup struct {
	ForemanObject

	PuppetClassIds []int `json:"puppetclass_ids,omitempty"`
}

func (c *Client) CreateConfigGroup(ctx context.Context, cg *ForemanConfigGroup) (*ForemanConfigGroup, error) {
	log.Tracef("foreman/api/config_group.go#Create")

	reqEndpoint := fmt.Sprintf("/%s", ConfigGroupEndpointPrefix)

	jsonBytes, jsonEncErr := c.WrapJSON("config_group", cg)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("configGroupJSONBytes: [%s]", jsonBytes)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodPost, reqEndpoint, bytes.NewBuffer(jsonBytes))
	if reqErr != nil {
		return nil, reqErr
	}

	var created ForemanConfigGroup
	if err := c.SendAndParse(req, &created); err != nil {
		return nil, err
	}

	log.Debugf("createdConfigGroup: [%+v]", created)
	return &created, nil
}

func (c *Client) ReadConfigGroup(ctx context.Context, id int) (*ForemanConfigGroup, error) {
	log.Tracef("foreman/api/config_group.go#Read")

	reqEndpoint := fmt.Sprintf("/%s/%d", ConfigGroupEndpointPrefix, id)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodGet, reqEndpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	var read ForemanConfigGroup
	if err := c.SendAndParse(req, &read); err != nil {
		return nil, err
	}

	log.Debugf("readConfigGroup: [%+v]", read)
	return &read, nil
}

func (c *Client) UpdateConfigGroup(ctx context.Context, cg *ForemanConfigGroup) (*ForemanConfigGroup, error) {
	log.Tracef("foreman/api/config_group.go#Update")

	reqEndpoint := fmt.Sprintf("/%s/%d", ConfigGroupEndpointPrefix, cg.Id)

	jsonBytes, jsonEncErr := c.WrapJSON("config_group", cg)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("configGroupJSONBytes: [%s]", jsonBytes)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodPut, reqEndpoint, bytes.NewBuffer(jsonBytes))
	if reqErr != nil {
		return nil, reqErr
	}

	var updated ForemanConfigGroup
	if err := c.SendAndParse(req, &updated); err != nil {
		return nil, err
	}

	log.Debugf("updatedConfigGroup: [%+v]", updated)
	return &updated, nil
}

func (c *Client) DeleteConfigGroup(ctx context.Context, id int) error {
	log.Tracef("foreman/api/config_group.go#Delete")

	reqEndpoint := fmt.Sprintf("/%s/%d", ConfigGroupEndpointPrefix, id)

	req, reqErr := c.NewRequestWithContext(ctx, http.MethodDelete, reqEndpoint, nil)
	if reqErr != nil {
		return reqErr
	}

	return c.SendAndParse(req, nil)
}

func (c *Client) QueryConfigGroup(ctx context.Context, cg *ForemanConfigGroup) (QueryResponse, error) {
	log.Tracef("foreman/api/config_group.go#Query")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", ConfigGroupEndpointPrefix)
	req, reqErr := c.NewRequestWithContext(ctx, http.MethodGet, reqEndpoint, nil)
	if reqErr != nil {
		return queryResponse, reqErr
	}

	reqQuery := req.URL.Query()
	reqQuery.Set("search", "name=\""+cg.Name+"\"")
	req.URL.RawQuery = reqQuery.Encode()

	if err := c.SendAndParse(req, &queryResponse); err != nil {
		return queryResponse, err
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	results := []ForemanConfigGroup{}
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

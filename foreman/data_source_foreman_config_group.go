package foreman

import (
	"context"
	"fmt"

	"github.com/HanseMerkur/terraform-provider-utils/autodoc"
	"github.com/HanseMerkur/terraform-provider-utils/helper"
	"github.com/HanseMerkur/terraform-provider-utils/log"
	"github.com/terraform-coop/terraform-provider-foreman/foreman/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceForemanConfigGroup() *schema.Resource {
	r := resourceForemanConfigGroup()
	ds := helper.DataSourceSchemaFromResourceSchema(r.Schema)

	ds["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: fmt.Sprintf(
			"The name of the config group. "+
				"%s \"base\"",
			autodoc.MetaExample,
		),
	}

	return &schema.Resource{
		ReadContext: dataSourceForemanConfigGroupRead,
		Schema:      ds,
	}
}

func dataSourceForemanConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("data_source_foreman_config_group.go#Read")

	client := meta.(*api.Client)
	cg := buildForemanConfigGroup(d)

	log.Debugf("ForemanConfigGroup: [%+v]", cg)

	queryResponse, err := client.QueryConfigGroup(ctx, cg)
	if err != nil {
		return diag.FromErr(err)
	}

	if queryResponse.Subtotal == 0 {
		return diag.Errorf("Data source config_group returned no results")
	} else if queryResponse.Subtotal > 1 {
		return diag.Errorf("Data source config_group returned more than 1 result")
	}

	result, ok := queryResponse.Results[0].(api.ForemanConfigGroup)
	if !ok {
		return diag.Errorf(
			"Data source results contain unexpected type. Expected "+
				"[api.ForemanConfigGroup], got [%T]",
			queryResponse.Results[0],
		)
	}

	log.Debugf("ForemanConfigGroup: [%+v]", result)

	setResourceDataFromForemanConfigGroup(d, &result)
	return nil
}

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

func dataSourceForemanLocation() *schema.Resource {
	r := resourceForemanLocation()
	ds := helper.DataSourceSchemaFromResourceSchema(r.Schema)

	ds["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: fmt.Sprintf(
			"The name of the location. "+
				"%s \"DC1\"",
			autodoc.MetaExample,
		),
	}

	return &schema.Resource{
		ReadContext: dataSourceForemanLocationRead,
		Schema:      ds,
	}
}

func dataSourceForemanLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("data_source_foreman_location.go#Read")

	client := meta.(*api.Client)
	l := buildForemanLocation(d)

	log.Debugf("ForemanLocation: [%+v]", l)

	queryResponse, err := client.QueryLocation(ctx, l)
	if err != nil {
		return diag.FromErr(err)
	}

	if queryResponse.Subtotal == 0 {
		return diag.Errorf("Data source location returned no results")
	} else if queryResponse.Subtotal > 1 {
		return diag.Errorf("Data source location returned more than 1 result")
	}

	result, ok := queryResponse.Results[0].(api.ForemanLocation)
	if !ok {
		return diag.Errorf(
			"Data source results contain unexpected type. Expected "+
				"[api.ForemanLocation], got [%T]",
			queryResponse.Results[0],
		)
	}

	log.Debugf("ForemanLocation: [%+v]", result)

	setResourceDataFromForemanLocation(d, &result)
	return nil
}

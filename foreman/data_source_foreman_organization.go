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

func dataSourceForemanOrganization() *schema.Resource {
	r := resourceForemanOrganization()
	ds := helper.DataSourceSchemaFromResourceSchema(r.Schema)

	ds["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: fmt.Sprintf(
			"The name of the organization. "+
				"%s \"MyOrg\"",
			autodoc.MetaExample,
		),
	}

	return &schema.Resource{
		ReadContext: dataSourceForemanOrganizationRead,
		Schema:      ds,
	}
}

func dataSourceForemanOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("data_source_foreman_organization.go#Read")

	client := meta.(*api.Client)
	o := buildForemanOrganization(d)

	log.Debugf("ForemanOrganization: [%+v]", o)

	queryResponse, err := client.QueryOrganization(ctx, o)
	if err != nil {
		return diag.FromErr(err)
	}

	if queryResponse.Subtotal == 0 {
		return diag.Errorf("Data source organization returned no results")
	} else if queryResponse.Subtotal > 1 {
		return diag.Errorf("Data source organization returned more than 1 result")
	}

	result, ok := queryResponse.Results[0].(api.ForemanOrganization)
	if !ok {
		return diag.Errorf(
			"Data source results contain unexpected type. Expected "+
				"[api.ForemanOrganization], got [%T]",
			queryResponse.Results[0],
		)
	}

	log.Debugf("ForemanOrganization: [%+v]", result)

	setResourceDataFromForemanOrganization(d, &result)
	return nil
}

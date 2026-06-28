package foreman

import (
	"context"
	"fmt"
	"strconv"

	"github.com/HanseMerkur/terraform-provider-utils/autodoc"
	"github.com/HanseMerkur/terraform-provider-utils/log"
	"github.com/terraform-coop/terraform-provider-foreman/foreman/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceForemanOrganization() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceForemanOrganizationCreate,
		ReadContext:   resourceForemanOrganizationRead,
		UpdateContext: resourceForemanOrganizationUpdate,
		DeleteContext: resourceForemanOrganizationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{

			autodoc.MetaAttribute: {
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s Foreman representation of an organization.",
					autodoc.MetaSummary,
				),
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"The name of the organization. "+
						"%s \"MyOrg\"",
					autodoc.MetaExample,
				),
			},

			"title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Full title of the organization (includes parent hierarchy).",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the organization.",
			},

			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the parent organization.",
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A map of parameters that will be saved as organization parameters.",
			},
		},
	}
}

func buildForemanOrganization(d *schema.ResourceData) *api.ForemanOrganization {
	log.Tracef("resource_foreman_organization.go#buildForemanOrganization")

	org := api.ForemanOrganization{}
	obj := buildForemanObject(d)
	org.ForemanObject = *obj

	if attr, ok := d.GetOk("description"); ok {
		org.Description = attr.(string)
	}
	if attr, ok := d.GetOk("parent_id"); ok {
		org.ParentId = attr.(int)
	}
	if attr, ok := d.GetOk("parameters"); ok {
		org.Parameters = api.ToKV(attr.(map[string]interface{}))
	}

	return &org
}

func setResourceDataFromForemanOrganization(d *schema.ResourceData, o *api.ForemanOrganization) {
	log.Tracef("resource_foreman_organization.go#setResourceDataFromForemanOrganization")

	d.SetId(strconv.Itoa(o.Id))
	d.Set("name", o.Name)
	d.Set("title", o.Title)
	d.Set("description", o.Description)
	d.Set("parent_id", o.ParentId)
	d.Set("parameters", api.FromKV(o.Parameters))
}

func resourceForemanOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_organization.go#Create")

	client := meta.(*api.Client)
	o := buildForemanOrganization(d)

	log.Debugf("ForemanOrganization: [%+v]", o)

	created, err := client.CreateOrganization(ctx, o)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debugf("Created ForemanOrganization: [%+v]", created)

	setResourceDataFromForemanOrganization(d, created)
	return nil
}

func resourceForemanOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_organization.go#Read")

	client := meta.(*api.Client)
	o := buildForemanOrganization(d)

	log.Debugf("ForemanOrganization: [%+v]", o)

	read, err := client.ReadOrganization(ctx, o.Id)
	if err != nil {
		return diag.FromErr(api.CheckDeleted(d, err))
	}

	log.Debugf("Read ForemanOrganization: [%+v]", read)

	setResourceDataFromForemanOrganization(d, read)
	return nil
}

func resourceForemanOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_organization.go#Update")

	client := meta.(*api.Client)
	o := buildForemanOrganization(d)

	log.Debugf("ForemanOrganization: [%+v]", o)

	updated, err := client.UpdateOrganization(ctx, o)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debugf("Updated ForemanOrganization: [%+v]", updated)

	setResourceDataFromForemanOrganization(d, updated)
	return nil
}

func resourceForemanOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_organization.go#Delete")

	client := meta.(*api.Client)
	o := buildForemanOrganization(d)

	log.Debugf("ForemanOrganization: [%+v]", o)

	return diag.FromErr(api.CheckDeleted(d, client.DeleteOrganization(ctx, o.Id)))
}

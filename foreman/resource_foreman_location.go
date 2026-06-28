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

func resourceForemanLocation() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceForemanLocationCreate,
		ReadContext:   resourceForemanLocationRead,
		UpdateContext: resourceForemanLocationUpdate,
		DeleteContext: resourceForemanLocationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{

			autodoc.MetaAttribute: {
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s Foreman representation of a location.",
					autodoc.MetaSummary,
				),
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"The name of the location. "+
						"%s \"DC1\"",
					autodoc.MetaExample,
				),
			},

			"title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Full title of the location (includes parent hierarchy).",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the location.",
			},

			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the parent location.",
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A map of parameters that will be saved as location parameters.",
			},
		},
	}
}

func buildForemanLocation(d *schema.ResourceData) *api.ForemanLocation {
	log.Tracef("resource_foreman_location.go#buildForemanLocation")

	loc := api.ForemanLocation{}
	obj := buildForemanObject(d)
	loc.ForemanObject = *obj

	if attr, ok := d.GetOk("description"); ok {
		loc.Description = attr.(string)
	}
	if attr, ok := d.GetOk("parent_id"); ok {
		loc.ParentId = attr.(int)
	}
	if attr, ok := d.GetOk("parameters"); ok {
		loc.Parameters = api.ToKV(attr.(map[string]interface{}))
	}

	return &loc
}

func setResourceDataFromForemanLocation(d *schema.ResourceData, l *api.ForemanLocation) {
	log.Tracef("resource_foreman_location.go#setResourceDataFromForemanLocation")

	d.SetId(strconv.Itoa(l.Id))
	d.Set("name", l.Name)
	d.Set("title", l.Title)
	d.Set("description", l.Description)
	d.Set("parent_id", l.ParentId)
	d.Set("parameters", api.FromKV(l.Parameters))
}

func resourceForemanLocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_location.go#Create")

	client := meta.(*api.Client)
	l := buildForemanLocation(d)

	log.Debugf("ForemanLocation: [%+v]", l)

	created, err := client.CreateLocation(ctx, l)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debugf("Created ForemanLocation: [%+v]", created)

	setResourceDataFromForemanLocation(d, created)
	return nil
}

func resourceForemanLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_location.go#Read")

	client := meta.(*api.Client)
	l := buildForemanLocation(d)

	log.Debugf("ForemanLocation: [%+v]", l)

	read, err := client.ReadLocation(ctx, l.Id)
	if err != nil {
		return diag.FromErr(api.CheckDeleted(d, err))
	}

	log.Debugf("Read ForemanLocation: [%+v]", read)

	setResourceDataFromForemanLocation(d, read)
	return nil
}

func resourceForemanLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_location.go#Update")

	client := meta.(*api.Client)
	l := buildForemanLocation(d)

	log.Debugf("ForemanLocation: [%+v]", l)

	updated, err := client.UpdateLocation(ctx, l)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debugf("Updated ForemanLocation: [%+v]", updated)

	setResourceDataFromForemanLocation(d, updated)
	return nil
}

func resourceForemanLocationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_location.go#Delete")

	client := meta.(*api.Client)
	l := buildForemanLocation(d)

	log.Debugf("ForemanLocation: [%+v]", l)

	return diag.FromErr(api.CheckDeleted(d, client.DeleteLocation(ctx, l.Id)))
}

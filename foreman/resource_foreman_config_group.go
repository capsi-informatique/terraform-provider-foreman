package foreman

import (
	"context"
	"fmt"
	"strconv"

	"github.com/HanseMerkur/terraform-provider-utils/autodoc"
	"github.com/HanseMerkur/terraform-provider-utils/conv"
	"github.com/HanseMerkur/terraform-provider-utils/log"
	"github.com/terraform-coop/terraform-provider-foreman/foreman/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceForemanConfigGroup() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceForemanConfigGroupCreate,
		ReadContext:   resourceForemanConfigGroupRead,
		UpdateContext: resourceForemanConfigGroupUpdate,
		DeleteContext: resourceForemanConfigGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{

			autodoc.MetaAttribute: {
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s Foreman representation of a config group (group of Puppet classes).",
					autodoc.MetaSummary,
				),
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"The name of the config group. "+
						"%s \"base\"",
					autodoc.MetaExample,
				),
			},

			"puppetclass_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Set of Puppet class IDs associated with this config group.",
			},
		},
	}
}

func buildForemanConfigGroup(d *schema.ResourceData) *api.ForemanConfigGroup {
	log.Tracef("resource_foreman_config_group.go#buildForemanConfigGroup")

	cg := api.ForemanConfigGroup{}
	obj := buildForemanObject(d)
	cg.ForemanObject = *obj

	if attr, ok := d.GetOk("puppetclass_ids"); ok {
		cg.PuppetClassIds = conv.InterfaceSliceToIntSlice(attr.(*schema.Set).List())
	}

	return &cg
}

func setResourceDataFromForemanConfigGroup(d *schema.ResourceData, cg *api.ForemanConfigGroup) {
	log.Tracef("resource_foreman_config_group.go#setResourceDataFromForemanConfigGroup")

	d.SetId(strconv.Itoa(cg.Id))
	d.Set("name", cg.Name)
	d.Set("puppetclass_ids", cg.PuppetClassIds)
}

func resourceForemanConfigGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_config_group.go#Create")

	client := meta.(*api.Client)
	cg := buildForemanConfigGroup(d)

	log.Debugf("ForemanConfigGroup: [%+v]", cg)

	created, err := client.CreateConfigGroup(ctx, cg)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debugf("Created ForemanConfigGroup: [%+v]", created)

	setResourceDataFromForemanConfigGroup(d, created)
	return nil
}

func resourceForemanConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_config_group.go#Read")

	client := meta.(*api.Client)
	cg := buildForemanConfigGroup(d)

	log.Debugf("ForemanConfigGroup: [%+v]", cg)

	read, err := client.ReadConfigGroup(ctx, cg.Id)
	if err != nil {
		return diag.FromErr(api.CheckDeleted(d, err))
	}

	log.Debugf("Read ForemanConfigGroup: [%+v]", read)

	setResourceDataFromForemanConfigGroup(d, read)
	return nil
}

func resourceForemanConfigGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_config_group.go#Update")

	client := meta.(*api.Client)
	cg := buildForemanConfigGroup(d)

	log.Debugf("ForemanConfigGroup: [%+v]", cg)

	updated, err := client.UpdateConfigGroup(ctx, cg)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debugf("Updated ForemanConfigGroup: [%+v]", updated)

	setResourceDataFromForemanConfigGroup(d, updated)
	return nil
}

func resourceForemanConfigGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_config_group.go#Delete")

	client := meta.(*api.Client)
	cg := buildForemanConfigGroup(d)

	log.Debugf("ForemanConfigGroup: [%+v]", cg)

	return diag.FromErr(api.CheckDeleted(d, client.DeleteConfigGroup(ctx, cg.Id)))
}

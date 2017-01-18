package cloudflare

import (
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudFlarePageRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFlarePageRuleCreate,
		Read:   resourceCloudFlarePageRuleRead,
		Update: resourceCloudFlarePageRuleUpdate,
		Delete: resourceCloudFlarePageRuleDelete,

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"target": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"actions": &schema.Schema{
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeMap,
					ValidateFunc: validatePageRuleAction,
				},
			},

			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  1,
				Optional: true,
			},

			"status": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "active",
				Optional:     true,
				ValidateFunc: validatePageRuleStatus,
			},
		},
	}
}

func resourceCloudFlarePageRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	actions := d.Get("actions").([]map[string]interface{})
	newPageRuleActions := make([]cloudflare.PageRuleAction, 0, len(actions))

	newPageRuleTargets := []cloudflare.PageRuleTarget{
		cloudflare.PageRuleTarget{
			Target: "url",
			Constraint: struct {
				Operator string `json:"operator"`
				Value    string `json:"value"`
			}{
				Operator: "matches",
				Value:    d.Get("target").(string),
			},
		},
	}

	for _, action := range actions {
		newPageRuleAction := cloudflare.PageRuleAction{
			ID:    action["action"].(string),
			Value: action["value"].(interface{}),
		}
		newPageRuleActions = append(newPageRuleActions, newPageRuleAction)
	}

	newPageRule := cloudflare.PageRule{
		Targets:  newPageRuleTargets,
		Actions:  newPageRuleActions,
		Priority: d.Get("priority").(int),
		Status:   d.Get("status").(string),
	}

	zoneName := d.Get("domain").(string)

	zoneId, err := client.ZoneIDByName(zoneName)
	if err != nil {
		return fmt.Errorf("Error finding zone %q: %s", zoneName, err)
	}

	d.Set("zone_id", zoneId)
	log.Printf("[DEBUG] CloudFlare Page Rule create configuration: %#v", newPageRule)

	err = client.CreatePageRule(zoneId, newPageRule)
	if err != nil {
		return fmt.Errorf("Failed to create page rule: %s", err)
	}

	return resourceCloudFlarePageRuleRead(d, meta)
}

func resourceCloudFlarePageRuleRead(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Page Rule Read not implemented.")
}

func resourceCloudFlarePageRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Page Rule Update not implemented.")
}

func resourceCloudFlarePageRuleDelete(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Page Rule Delete not implemented.")
}

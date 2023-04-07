// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceNetworkSecurityUrlLists() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkSecurityUrlListsCreate,
		Read:   resourceNetworkSecurityUrlListsRead,
		Update: resourceNetworkSecurityUrlListsUpdate,
		Delete: resourceNetworkSecurityUrlListsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkSecurityUrlListsImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The location of the url lists.`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Short name of the UrlList resource to be created.
This value should be 1-63 characters long, containing only letters, numbers, hyphens, and underscores, and should not start with a number. E.g. 'urlList'.`,
			},
			"values": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `FQDNs and URLs.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Free-text description of the resource.`,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. Time when the security policy was created.
A timestamp in RFC3339 UTC 'Zulu' format, with nanosecond resolution and up to nine fractional digits.
Examples: '2014-10-02T15:01:23Z' and '2014-10-02T15:01:23.045123456Z'`,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. Time when the security policy was updated.
A timestamp in RFC3339 UTC 'Zulu' format, with nanosecond resolution and up to nine fractional digits.
Examples: '2014-10-02T15:01:23Z' and '2014-10-02T15:01:23.045123456Z'.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceNetworkSecurityUrlListsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandNetworkSecurityUrlListsDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	valuesProp, err := expandNetworkSecurityUrlListsValues(d.Get("values"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("values"); !isEmptyValue(reflect.ValueOf(valuesProp)) && (ok || !reflect.DeepEqual(v, valuesProp)) {
		obj["values"] = valuesProp
	}

	url, err := replaceVars(d, config, "{{NetworkSecurityBasePath}}projects/{{project}}/locations/{{location}}/urlLists?urlListId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new UrlLists: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for UrlLists: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating UrlLists: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{location}}/urlLists/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = NetworkSecurityOperationWaitTime(
		config, res, project, "Creating UrlLists", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create UrlLists: %s", err)
	}

	log.Printf("[DEBUG] Finished creating UrlLists %q: %#v", d.Id(), res)

	return resourceNetworkSecurityUrlListsRead(d, meta)
}

func resourceNetworkSecurityUrlListsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{NetworkSecurityBasePath}}projects/{{project}}/locations/{{location}}/urlLists/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for UrlLists: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("NetworkSecurityUrlLists %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading UrlLists: %s", err)
	}

	if err := d.Set("create_time", flattenNetworkSecurityUrlListsCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading UrlLists: %s", err)
	}
	if err := d.Set("update_time", flattenNetworkSecurityUrlListsUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading UrlLists: %s", err)
	}
	if err := d.Set("description", flattenNetworkSecurityUrlListsDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading UrlLists: %s", err)
	}
	if err := d.Set("values", flattenNetworkSecurityUrlListsValues(res["values"], d, config)); err != nil {
		return fmt.Errorf("Error reading UrlLists: %s", err)
	}

	return nil
}

func resourceNetworkSecurityUrlListsUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for UrlLists: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	descriptionProp, err := expandNetworkSecurityUrlListsDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	valuesProp, err := expandNetworkSecurityUrlListsValues(d.Get("values"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("values"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, valuesProp)) {
		obj["values"] = valuesProp
	}

	url, err := replaceVars(d, config, "{{NetworkSecurityBasePath}}projects/{{project}}/locations/{{location}}/urlLists/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating UrlLists %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("values") {
		updateMask = append(updateMask, "values")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating UrlLists %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating UrlLists %q: %#v", d.Id(), res)
	}

	err = NetworkSecurityOperationWaitTime(
		config, res, project, "Updating UrlLists", userAgent,
		d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return err
	}

	return resourceNetworkSecurityUrlListsRead(d, meta)
}

func resourceNetworkSecurityUrlListsDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for UrlLists: %s", err)
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{NetworkSecurityBasePath}}projects/{{project}}/locations/{{location}}/urlLists/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting UrlLists %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "UrlLists")
	}

	err = NetworkSecurityOperationWaitTime(
		config, res, project, "Deleting UrlLists", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting UrlLists %q: %#v", d.Id(), res)
	return nil
}

func resourceNetworkSecurityUrlListsImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/urlLists/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{location}}/urlLists/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNetworkSecurityUrlListsCreateTime(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenNetworkSecurityUrlListsUpdateTime(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenNetworkSecurityUrlListsDescription(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenNetworkSecurityUrlListsValues(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandNetworkSecurityUrlListsDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandNetworkSecurityUrlListsValues(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
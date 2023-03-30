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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceComputeOrganizationSecurityPolicyAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeOrganizationSecurityPolicyAssociationCreate,
		Read:   resourceComputeOrganizationSecurityPolicyAssociationRead,
		Delete: resourceComputeOrganizationSecurityPolicyAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeOrganizationSecurityPolicyAssociationImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The resource that the security policy is attached to.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name for an association.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The security policy ID of the association.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The display name of the security policy of the association.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceComputeOrganizationSecurityPolicyAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandComputeOrganizationSecurityPolicyAssociationName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	attachmentIdProp, err := expandComputeOrganizationSecurityPolicyAssociationAttachmentId(d.Get("attachment_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("attachment_id"); !isEmptyValue(reflect.ValueOf(attachmentIdProp)) && (ok || !reflect.DeepEqual(v, attachmentIdProp)) {
		obj["attachmentId"] = attachmentIdProp
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}{{policy_id}}/addAssociation")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new OrganizationSecurityPolicyAssociation: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating OrganizationSecurityPolicyAssociation: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{policy_id}}/association/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// `parent` is needed to poll the asynchronous operations but its available only on the policy.

	policyUrl := fmt.Sprintf("{{ComputeBasePath}}%s", d.Get("policy_id"))
	url, err = replaceVars(d, config, policyUrl)
	if err != nil {
		return err
	}

	policyRes, err := SendRequest(config, "GET", "", url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeOrganizationSecurityPolicy %q", d.Get("policy_id")))
	}

	parent := flattenComputeOrganizationSecurityPolicyParent(policyRes["parent"], d, config)
	var opRes map[string]interface{}
	err = ComputeOrgOperationWaitTimeWithResponse(
		config, res, &opRes, parent.(string), "Creating OrganizationSecurityPolicyAssociation", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create OrganizationSecurityPolicyAssociation: %s", err)
	}

	log.Printf("[DEBUG] Finished creating OrganizationSecurityPolicyAssociation %q: %#v", d.Id(), res)

	return resourceComputeOrganizationSecurityPolicyAssociationRead(d, meta)
}

func resourceComputeOrganizationSecurityPolicyAssociationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}{{policy_id}}/getAssociation?name={{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(transformSecurityPolicyAssociationReadError(err), d, fmt.Sprintf("ComputeOrganizationSecurityPolicyAssociation %q", d.Id()))
	}

	if err := d.Set("name", flattenComputeOrganizationSecurityPolicyAssociationName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSecurityPolicyAssociation: %s", err)
	}
	if err := d.Set("attachment_id", flattenComputeOrganizationSecurityPolicyAssociationAttachmentId(res["attachmentId"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSecurityPolicyAssociation: %s", err)
	}
	if err := d.Set("display_name", flattenComputeOrganizationSecurityPolicyAssociationDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSecurityPolicyAssociation: %s", err)
	}

	return nil
}

func resourceComputeOrganizationSecurityPolicyAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := replaceVars(d, config, "{{ComputeBasePath}}{{policy_id}}/removeAssociation?name={{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting OrganizationSecurityPolicyAssociation %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "OrganizationSecurityPolicyAssociation")
	}

	// `parent` is needed to poll the asynchronous operations but its available only on the policy.

	policyUrl := fmt.Sprintf("{{ComputeBasePath}}%s", d.Get("policy_id"))
	url, err = replaceVars(d, config, policyUrl)
	if err != nil {
		return err
	}

	policyRes, err := SendRequest(config, "GET", "", url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeOrganizationSecurityPolicy %q", d.Get("policy_id")))
	}

	parent := flattenComputeOrganizationSecurityPolicyParent(policyRes["parent"], d, config)
	var opRes map[string]interface{}
	err = ComputeOrgOperationWaitTimeWithResponse(
		config, res, &opRes, parent.(string), "Creating OrganizationSecurityPolicyAssociation", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create OrganizationSecurityPolicyAssociation: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting OrganizationSecurityPolicyAssociation %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeOrganizationSecurityPolicyAssociationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"(?P<policy_id>.+)/association/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{policy_id}}/association/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeOrganizationSecurityPolicyAssociationName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeOrganizationSecurityPolicyAssociationAttachmentId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeOrganizationSecurityPolicyAssociationDisplayName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandComputeOrganizationSecurityPolicyAssociationName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeOrganizationSecurityPolicyAssociationAttachmentId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

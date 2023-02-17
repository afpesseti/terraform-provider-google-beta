package google

import (
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	compute "google.golang.org/api/compute/v0.beta"
	"log"
	"time"
)

// Change
func resourceComputeNetworkEdgeSecurityServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNetworkEdgeSecurityServicesCreate,
		Read:   resourceComputeNetworkEdgeSecurityServicesRead,
		Update: resourceComputeNetworkEdgeSecurityServicesUpdate,
		Delete: resourceComputeNetworkEdgeSecurityServicesDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkEdgeSecurityServicesImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCEName,
				Description:  `Name of the resource. Provided by the client when the resource is created.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `An optional description of this resource. Provide this property when you create the resource.`,
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Fingerprint of this resource. A hash of the contents stored in this object. This field is used in optimistic locking. This field will be ignored when inserting a NetworkEdgeSecurityService.`,
			},

			"security_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The resource URL for the network edge security service associated with this network edge security service.`,
			},
			"region": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `URL of the region where the resource resides.`,
			},
		},

		UseJSONNumber: true,
	}
}

func resourceComputeNetworkEdgeSecurityServicesCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)
	networkEdgeSecurityServices := &compute.NetworkEdgeSecurityService{
		Name:        sp,
		Description: d.Get("description").(string),
	}

	if v, ok := d.GetOk("region"); ok {
		networkEdgeSecurityServices.Region = v.(string)
	}

	if v, ok := d.GetOk("security_policy"); ok {
		networkEdgeSecurityServices.SecurityPolicy = v.(string)
	}

	log.Printf("[DEBUG] NetworkEdgeSecurityService insert request: %#v", networkEdgeSecurityServices)

	client := config.NewComputeClient(userAgent)

	op, err := client.NetworkEdgeSecurityServices.Insert(project, region, networkEdgeSecurityServices).Do()

	if err != nil {
		return errwrap.Wrapf("Error creating NetworkEdgeSecurityService: {{err}}", err)
	}

	id, err := replaceVars(d, config, "projects/{{project}}/regions/{{region}}/networkEdgeSecurityServices/{{name}}")
	if err != nil {
		fmt.Print("-------------------------------------------Problem to INSERT!!-----------------------------------------\n")
		fmt.Print(err)
		fmt.Print("\n")
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Creating NetworkEdgeSecurityService %q", sp), userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		fmt.Print("-------------------------------------------Problem to INSERT 02!!-----------------------------------------\n")
		fmt.Print(err)
		fmt.Print("\n")
		return err
	}

	fmt.Print("-------------------------------------------Insert sucessfully!!-----------------------------------------\n")
	fmt.Print("\n")

	return resourceComputeNetworkEdgeSecurityServicesRead(d, meta)
}

func resourceComputeNetworkEdgeSecurityServicesRead(d *schema.ResourceData, meta interface{}) error {
	fmt.Print("-------------------------------------------Prepare To READ!!-----------------------------------------\n")
	fmt.Print("\n")

	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	client := config.NewComputeClient(userAgent)

	networkEdgeSecurityServices, err := client.NetworkEdgeSecurityServices.Get(project, region, sp).Do()
	if err != nil {
		fmt.Print("-------------------------------------------Problem To READ!!-----------------------------------------\n")
		fmt.Print("\n")
		fmt.Print(err)
		return handleNotFoundError(err, d, fmt.Sprintf("NetworkEdgeSecurityServices %q", d.Id()))
	}

	if err := d.Set("name", networkEdgeSecurityServices.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", networkEdgeSecurityServices.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("region", networkEdgeSecurityServices.Region); err != nil {
		fmt.Printf("Error setting region: %s", err)
	}
	if err := d.Set("fingerprint", networkEdgeSecurityServices.Fingerprint); err != nil {
		return fmt.Errorf("Error setting fingerprint: %s", err)
	}
	if err := d.Set("security_policy", networkEdgeSecurityServices.SecurityPolicy); err != nil {
		return fmt.Errorf("Error setting security policy: %s", err)
	}

	fmt.Print("-------------------------------------------Read sucessfully!!-----------------------------------------\n")
	fmt.Print("\n")

	return nil
}

func resourceComputeNetworkEdgeSecurityServicesUpdate(d *schema.ResourceData, meta interface{}) error {
	fmt.Print("-------------------------------------------Prepare To Update!!-----------------------------------------\n")
	fmt.Print("\n")

	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	networkEdgeSecurityServices := &compute.NetworkEdgeSecurityService{
		Fingerprint: d.Get("fingerprint").(string),
	}

	if d.HasChange("description") {
		networkEdgeSecurityServices.Description = d.Get("description").(string)
		networkEdgeSecurityServices.ForceSendFields = append(networkEdgeSecurityServices.ForceSendFields, "Description")
	}

	if d.HasChange("region") {
		networkEdgeSecurityServices.Region = d.Get("region").(string)
		networkEdgeSecurityServices.ForceSendFields = append(networkEdgeSecurityServices.ForceSendFields, "Region")
	}

	if d.HasChange("security_policy") {
		networkEdgeSecurityServices.SecurityPolicy = d.Get("security_policy").(string)
		networkEdgeSecurityServices.ForceSendFields = append(networkEdgeSecurityServices.ForceSendFields, "SecurityPolicy")
	}

	if len(networkEdgeSecurityServices.ForceSendFields) > 0 {
		client := config.NewComputeClient(userAgent)

		op, err := client.NetworkEdgeSecurityServices.Patch(project, region, sp, networkEdgeSecurityServices).Do()

		if err != nil {
			fmt.Print("-------------------------------------------Problem To Update!!-----------------------------------------\n")
			fmt.Print(err)
			fmt.Print("\n")
			return errwrap.Wrapf(fmt.Sprintf("Error updating NetworkEdgeSecurityServices %q: {{err}}", sp), err)
		}

		err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating NetworkEdgeSecurityServices %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			fmt.Print("-------------------------------------------Problem To Update 02!!-----------------------------------------\n")
			fmt.Print(err)
			fmt.Print("\n")
			return err
		}
	}

	fmt.Print("-------------------------------------------Update sucessfully!!-----------------------------------------\n")
	fmt.Print("\n")

	return resourceComputeNetworkEdgeSecurityServicesRead(d, meta)
}

func resourceComputeNetworkEdgeSecurityServicesDelete(d *schema.ResourceData, meta interface{}) error {
	fmt.Print("-------------------------------------------Prepare To DELETE!!-----------------------------------------\n")
	fmt.Print("\n")

	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	client := config.NewComputeClient(userAgent)

	// Delete the SecurityPolicy
	op, err := client.NetworkEdgeSecurityServices.Delete(project, region, d.Get("name").(string)).Do()
	if err != nil {
		fmt.Print("-------------------------------------------Problem To DELETE!!-----------------------------------------\n")
		fmt.Print(err)
		fmt.Print("\n")
		return errwrap.Wrapf("Error deleting NetworkEdgeSecurityServices: {{err}}", err)
	}

	err = computeOperationWaitTime(config, op, project, "Deleting NetworkEdgeSecurityServices", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		fmt.Print("-------------------------------------------Problem To DELETE computeOperationWaitTime!!-----------------------------------------\n")
		fmt.Print(err)
		fmt.Print("\n")
		return err
	}

	fmt.Print("-------------------------------------------Delete sucessfully!!-----------------------------------------\n")

	d.SetId("")
	return nil
}

func resourceNetworkEdgeSecurityServicesImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/networkEdgeSecurityServices/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		fmt.Print("-------------------------------------------Problem to Import!!-----------------------------------------\n")
		fmt.Print(err)
		fmt.Print("\n")
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/regions/{{region}}/networkEdgeSecurityServices/{{name}}")
	if err != nil {
		fmt.Print("-------------------------------------------Problem on ReplaceVars!!-----------------------------------------\n")
		fmt.Print(err)
		fmt.Print("\n")
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

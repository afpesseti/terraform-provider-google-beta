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
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalNetworkEndpointGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupExample(context),
			},
			{
				ResourceName:      "google_compute_global_network_endpoint_group.neg",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_global_network_endpoint_group" "neg" {
  name                  = "tf-test-my-lb-neg%{random_suffix}"
  default_port          = "90"
  network_endpoint_type = "INTERNET_FQDN_PORT"
}
`, context)
}

func TestAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupIpAddressExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalNetworkEndpointGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupIpAddressExample(context),
			},
			{
				ResourceName:      "google_compute_global_network_endpoint_group.neg",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupIpAddressExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_global_network_endpoint_group" "neg" {
  name                  = "tf-test-my-lb-neg%{random_suffix}"
  network_endpoint_type = "INTERNET_IP_PORT"
  default_port          = 90
}
`, context)
}

func testAccCheckComputeGlobalNetworkEndpointGroupDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_global_network_endpoint_group" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/global/networkEndpointGroups/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("ComputeGlobalNetworkEndpointGroup still exists at %s", url)
			}
		}

		return nil
	}
}

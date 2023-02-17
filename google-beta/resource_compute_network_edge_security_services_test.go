package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccComputeRegionSecurityPolicy_basic(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polLink := "google_compute_region_security_policy.policy.self_link"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSecurityPolicy_basic(spName),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSecurityPolicy_withDdosProtection(spName),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeNetworkEdgeSecurityServices_basic(spName, polLink),
			},
			{
				ResourceName:      "google_compute_network_edge_security_services.services",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeNetworkEdge_basic(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polLink := "google_compute_region_security_policy.policy.self_link"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEdgeSecurityServices_basic(spName, polLink),
			},
			{
				ResourceName:      "google_compute_network_edge_security_services.services",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionSecurityPolicy_basic(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_security_policy" "policy" {
  name        = "%s"
  description = "default rule"
  type = "CLOUD_ARMOR_NETWORK"
}
`, spName)
}

func testAccComputeRegionSecurityPolicy_withDdosProtection(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_security_policy" "policy" {
  name        = "%s"
  description = "update rule"
  type = "CLOUD_ARMOR_NETWORK"

  ddos_protection_config {
    ddos_protection = "ADVANCED"
  }
}
`, spName)
}

func testAccComputeNetworkEdgeSecurityServices_basic(spName, polLink string) string {
	return fmt.Sprintf(`
resource "google_compute_network_edge_security_services" "services" {
	name        = "%s"
	description = "basic network edge security services"

}
`, spName)
}

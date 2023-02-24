package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccComputeNetworkEdgeSecurityServices_basic(t *testing.T) {
	t.Parallel()

	polName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEdgeSecurityServices_basic(polName),
			},
			{
				ResourceName:      "google_compute_network_edge_security_services.services",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeNetworkEdgeSecurityServices_withRegionSecurityPolicy(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEdgeSecurityServices_withRegionSecurityPolicy_basic(polName),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeNetworkEdgeSecurityServices_withRegionSecurityPolicy(serviceName, polName, "google_compute_region_security_policy.policy.self_link"),
			},
			{
				ResourceName:      "google_compute_network_edge_security_services.services",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeNetworkEdgeSecurityServices_withRegionSecurityPolicy_basic(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_security_policy" "policy" {
  name        = "%s"
  description = "default rule"
  type = "CLOUD_ARMOR_NETWORK"

  ddos_protection_config {
    ddos_protection = "STANDARD"
  }
}
`, spName)
}

func testAccComputeNetworkEdgeSecurityServices_basic(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_network_edge_security_services" "services" {
	name        = "%s"
	description = "basic network edge security services"
}
`, spName)
}

func testAccComputeNetworkEdgeSecurityServices_withRegionSecurityPolicy(serviceName, polName, polLink string) string {
	return fmt.Sprintf(`
resource "google_compute_network_edge_security_services" "services" {
	name        = "%s"
	description = "basic network edge security services"
	security_policy = %s
}

resource "google_compute_region_security_policy" "policy" {
	name        = "%s"
	description = "default rule"
	type = "CLOUD_ARMOR_NETWORK"

	ddos_protection_config {
		ddos_protection = "ADVANCED"
	}
}
`, serviceName, polLink, polName)
}

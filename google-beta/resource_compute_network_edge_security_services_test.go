package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// Change
func TestAccComputeNetworkEdgeSecurityServices_withSecurityPolicies(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polLink := "google_compute_security_policy.policy.self_link"
	//polNameStandard := fmt.Sprintf("tf-test-%s", randString(t, 10))
	//clearpolNameAdvanced := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEdgeSecurityServices_withSecurityPolicies(spName, polLink),
			},
			{
				ResourceName:      "google_compute_network_edge_security_services.services",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeSecurityPolicy_withDdosProtectionConfig(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "default rule"
  type = "CLOUD_ARMOR"
}
`, spName)
}

func testAccComputeNetworkEdgeSecurityServices_withSecurityPolicies(spName, polLink string) string {
	return fmt.Sprintf(`
resource "google_compute_network_edge_security_services" "services" {
	name        = "%s"
	description = "basic network edge security services"
}
`, spName)
}

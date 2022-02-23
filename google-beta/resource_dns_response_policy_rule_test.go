package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSResponsePolicyRule_update(t *testing.T) {
	t.Parallel()

	responsePolicyRuleSuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersOiCS,
		CheckDestroy: testAccCheckDNSResponsePolicyRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsResponsePolicyRule_privateUpdate(responsePolicyRuleSuffix, "network-1"),
			},
			{
				ResourceName:      "google_dns_response_policy_rule.example-response-policy-rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsResponsePolicyRule_privateUpdate(responsePolicyRuleSuffix, "network-2"),
			},
			{
				ResourceName:      "google_dns_response_policy_rule.example-response-policy-rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsResponsePolicyRule_privateUpdate(suffix, network string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network-1" {
  provider = google-beta

  name                    = "tf-test-network-1-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network-2" {
  provider = google-beta

  name                    = "tf-test-network-2-%s"
  auto_create_subnetworks = false
}

resource "google_dns_response_policy" "response-policy" {
  provider = google-beta

  response_policy_name = "tf-test-response-policy-%s"

  networks {
    network_url = google_compute_network.%s.self_link
  }
}

resource "google_dns_response_policy_rule" "example-response-policy-rule" {
  provider = google-beta

  response_policy = google_dns_response_policy.response-policy.response_policy_name
  rule_name       = "tf-test-response-policy-rule-%s"
  dns_name        = "dns.example.com."

  local_data {
    local_datas {
      name   = "dns.example.com."
      type   = "A"
      ttl    = 300
      rrdatas = ["192.0.2.91"]
    }
  }
}  
`, suffix, suffix, suffix, network, suffix)
}
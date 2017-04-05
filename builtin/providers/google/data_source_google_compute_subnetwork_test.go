package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleSubnetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccDataSourceGoogleSubnetworkConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleSubnetworkCheck("data.google_compute_subnetwork.my_subnetwork", "google_compute_subnetwork.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleSubnetworkCheck(name string, subnetwork_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}
		network, ok := s.RootModule().Resources[subnetwork_name]
		if !ok {
			return fmt.Errorf("can't find google_compute_network.foobar in state")
		}

		subnetworkOrigin, ok := s.RootModule().Resources["google_compute_subnetwork.foobar"]
		if !ok {
			return fmt.Errorf("can't find google_compute_subnetwork.foobar in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != subnetworkOrigin.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				subnetworkOrigin.Primary.Attributes["id"],
			)
		}

		if attr["self_link"] != subnetworkOrigin.Primary.Attributes["self_link"] {
			return fmt.Errorf(
				"self_link is %s; want %s",
				attr["self_link"],
				subnetworkOrigin.Primary.Attributes["self_link"],
			)
		}

		if attr["name"] != subnetworkOrigin.Primary.Attributes["name"] {
			return fmt.Errorf("bad name %s", attr["name"])
		}

		if attr["ip_cidr_range"] != subnetworkOrigin.Primary.Attributes["ip_cidr_range"] {
			return fmt.Errorf("bad ip_cidr_range %s", attr["ip_cidr_range"])
		}
		if attr["network"] != network.Primary.Attributes["network"] {
			return fmt.Errorf("bad network_name %s", attr["network"])
		}

		if attr["description"] != subnetworkOrigin.Primary.Attributes["description"] {
			return fmt.Errorf("bad description %s", attr["description"])
		}
		return nil
	}
}

var TestAccDataSourceGoogleSubnetworkConfig = `

resource "google_compute_network" "foobar" {
	name = "network-test"
	description = "my-description"
}
resource "google_compute_subnetwork" "foobar" {
	name = "subnetwork-test"
	description = "my-description"
	ip_cidr_range = "10.0.0.0/24"
	network  = "${google_compute_network.foobar.self_link}"
}

data "google_compute_subnetwork" "my_subnetwork" {
	name = "${google_compute_subnetwork.foobar.name}"
}
`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO: Add entry data source tests, below is example code

func TestAccEntryDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEntryDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dvls_Entry.test", "id", "Entry-id"),
				),
			},
		},
	})
}

const testAccEntryDataSourceConfig = `
data "dvls_Entry" "test" {
  configurable_attribute = "Entry"
}
`

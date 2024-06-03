package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO: Add entryusercredential data source tests, below is example code

func TestAccEntryUserCredentialDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEntryUserCredentialDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dvls_EntryUserCredential.test", "id", "EntryUserCredential-id"),
				),
			},
		},
	})
}

const testAccEntryUserCredentialDataSourceConfig = `
data "dvls_EntryUserCredential" "test" {
  configurable_attribute = "EntryUserCredential"
}
`

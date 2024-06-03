package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO: Add entryusercredential resource tests, below is example code

func TestAccEntryUserCredentialResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEntryUserCredentialResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_EntryUserCredential.test", "configurable_attribute", "one"),
					resource.TestCheckResourceAttr("dvls_EntryUserCredential.test", "id", "EntryUserCredential-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "dvls_EntryUserCredential.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// EntryUserCredential code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute"},
			},
			// Update and Read testing
			{
				Config: testAccEntryUserCredentialResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_EntryUserCredential.test", "configurable_attribute", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEntryUserCredentialResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "dvls_EntryUserCredential" "test" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}

package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"os"
	"testing"
)

func TestAccDownloadFileFunction_Simple(t *testing.T) {
	_ = os.Remove("file.dat") // remove existing test file

	config := `
output "test" {
  value = provider::download::file("http://localhost:8080/file.dat", "file.dat")
}
`
	resource.Test(t, resource.TestCase{
		PreCheck: func() {},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("file.dat")),
				},
			},
		},
	})
}

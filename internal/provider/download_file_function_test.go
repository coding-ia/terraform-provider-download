package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"os"
	"regexp"
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

func TestAccDownloadFileFunction_EmptyOrBadURL(t *testing.T) {
	expectedError, _ := regexp.Compile(".*invalid url.*")
	config := `
output "test" {
  value = provider::download::file("", "file.dat")
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
				Config:      config,
				ExpectError: expectedError,
			},
		},
	})
}

func TestAccDownloadFileFunction_NoOutputFile(t *testing.T) {
	expectedError, _ := regexp.Compile(".*filename is empty.*")
	config := `
output "test" {
  value = provider::download::file("http://localhost:8080/file.dat", "")
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
				Config:      config,
				ExpectError: expectedError,
			},
		},
	})
}

func TestAccDownloadFileFunction_URLNotFound(t *testing.T) {
	expectedError, _ := regexp.Compile(".*bad status: 404 Not Found.*")
	config := `
output "test" {
  value = provider::download::file("http://localhost:8080/file2.dat", "file.dat")
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
				Config:      config,
				ExpectError: expectedError,
			},
		},
	})
}

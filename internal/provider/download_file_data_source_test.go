package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccDownloadDataSourceDownloadFile_Simple(t *testing.T) {
	config := `
data "download_file" "test" {
  url           = "http://localhost:8080/file.dat"
  output_file   = "file.dat"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.download_file.test", "output_base64sha256", "VkfwXsGJWJR9ModO63iPo5agXQurfBtx8RLOt+mzHu4="),
					resource.TestCheckResourceAttr("data.download_file.test", "output_sha256", "5647f05ec18958947d32874eeb788fa396a05d0bab7c1b71f112ceb7e9b31eee"),
					resource.TestCheckResourceAttr("data.download_file.test", "output_sha", "7d76d48d64d7ac5411d714a4bb83f37e3e5b8df6"),
					resource.TestCheckResourceAttr("data.download_file.test", "output_md5", "b2d1236c286a3c0704224fe4105eca49"),
				),
			},
		},
	})
}

func TestAccDownloadDataSourceDownloadFile_NoOutputFile(t *testing.T) {
	expectedError, _ := regexp.Compile(".*open : no such file or directory.*")
	config := `
data "download_file" "test" {
  url           = "http://localhost:8080/file.dat"
  output_file   = ""
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expectedError,
			},
		},
	})
}

func TestAccDownloadDataSourceDownloadFile_EmptyUrl(t *testing.T) {
	expectedError, _ := regexp.Compile(".*URL cannot be empty.*")
	config := `
data "download_file" "test" {
  url           = ""
  output_file   = ""
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expectedError,
			},
		},
	})
}

func TestAccDownloadDataSourceDownloadFile_Verify(t *testing.T) {
	config := `
data "download_file" "test" {
  url           = "http://localhost:8080/file.dat"
  output_file   = "file.dat"

  verify_sha256 = "5647f05ec18958947d32874eeb788fa396a05d0bab7c1b71f112ceb7e9b31eee"
  verify_sha    = "7d76d48d64d7ac5411d714a4bb83f37e3e5b8df6"
  verify_md5    = "b2d1236c286a3c0704224fe4105eca49"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.download_file.test", "output_base64sha256", "VkfwXsGJWJR9ModO63iPo5agXQurfBtx8RLOt+mzHu4="),
					resource.TestCheckResourceAttr("data.download_file.test", "output_sha256", "5647f05ec18958947d32874eeb788fa396a05d0bab7c1b71f112ceb7e9b31eee"),
					resource.TestCheckResourceAttr("data.download_file.test", "output_sha", "7d76d48d64d7ac5411d714a4bb83f37e3e5b8df6"),
					resource.TestCheckResourceAttr("data.download_file.test", "output_md5", "b2d1236c286a3c0704224fe4105eca49"),
				),
			},
		},
	})
}

func TestAccDownloadDataSourceDownloadFile_VerifyInvalidSha256(t *testing.T) {
	expectedError, _ := regexp.Compile(".*SHA256 signature mismatch.*")
	config := `
data "download_file" "test" {
  url           = "http://localhost:8080/file.dat"
  output_file   = "file.dat"

  verify_sha256 = "00000"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expectedError,
			},
		},
	})
}

func TestAccDownloadDataSourceDownloadFile_VerifyInvalidSha(t *testing.T) {
	expectedError, _ := regexp.Compile(".*SHA1 signature mismatch.*")
	config := `
data "download_file" "test" {
  url           = "http://localhost:8080/file.dat"
  output_file   = "file.dat"

  verify_sha = "00000"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expectedError,
			},
		},
	})
}

func TestAccDownloadDataSourceDownloadFile_VerifyInvalidMD5(t *testing.T) {
	expectedError, _ := regexp.Compile(".*MD5 signature mismatch.*")
	config := `
data "download_file" "test" {
  url           = "http://localhost:8080/file.dat"
  output_file   = "file.dat"

  verify_md5 = "00000"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expectedError,
			},
		},
	})
}

func TestAccDownloadDataSourceDownloadFile_InvalidUrl(t *testing.T) {
	expectedError, _ := regexp.Compile(".*bad status: 404 Not Found.*")
	config := `
data "download_file" "test" {
  url           = "http://localhost:8080/file2.dat"
  output_file   = "file.dat"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expectedError,
			},
		},
	})
}

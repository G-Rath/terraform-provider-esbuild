package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccESBuildBundleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccESBuildBundleBasicDataSourceConfig("basic.ts"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "id", "e76a68d60445036e1c6f26a0d1d9ba0a58cd1fb9e2a915fb0d38e0f1df4290ec"),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "filename", "./fixtures/basic.ts"),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "content", testAccESBuildBundleDefaultContentNode),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "target", "esnext"),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "platform", "node"),
				),
			},
		},
	})
}

func TestAccESBuildBundleDataSource_FileWithImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccESBuildBundleBasicDataSourceConfig("import.ts"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "id", "e76a68d60445036e1c6f26a0d1d9ba0a58cd1fb9e2a915fb0d38e0f1df4290ec"),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "filename", "./fixtures/import.ts"),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "content", testAccESBuildBundleDefaultContentNode),
				),
			},
		},
	})
}

func TestAccESBuildBundleDataSource_FileDoesNotExist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccESBuildBundleBasicDataSourceConfig("does-not-exist.ts"),
				ExpectError: regexp.MustCompile(`Could not resolve "\./fixtures/does-not-exist\.ts"`),
			},
		},
	})
}

func TestAccESBuildBundleDataSource_FileHasSyntaxError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccESBuildBundleBasicDataSourceConfig("syntax-error.ts"),
				ExpectError: regexp.MustCompile(`Expected "=>" but found "=="`),
			},
		},
	})
}

func TestAccESBuildBundleDataSource_Target_Valid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccESBuildBundleBasicDataSourceConfigWithTarget("basic.ts", "es2020"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "id", "e76a68d60445036e1c6f26a0d1d9ba0a58cd1fb9e2a915fb0d38e0f1df4290ec"),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "filename", "./fixtures/basic.ts"),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "content", testAccESBuildBundleDefaultContentNode),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "target", "es2020"),
				),
			},
		},
	})
}

func TestAccESBuildBundleDataSource_Target_Invalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccESBuildBundleBasicDataSourceConfigWithTarget("basic.ts", "not-a-target"),
				ExpectError: regexp.MustCompile(`unknown target "not-a-target"`),
			},
		},
	})
}

func TestAccESBuildBundleDataSource_Platform_Valid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccESBuildBundleBasicDataSourceConfigWithPlatform("basic.ts", "browser"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "id", "65c3f40b0b4ae806a1c648b9984b5249e11061c1b28b65da66dc4a01f721e1c2"),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "filename", "./fixtures/basic.ts"),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "content", testAccESBuildBundleDefaultContentBrowser),
					resource.TestCheckResourceAttr("data.esbuild_bundle.test", "platform", "browser"),
				),
			},
		},
	})
}

func TestAccESBuildBundleDataSource_Platform_Invalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccESBuildBundleBasicDataSourceConfigWithPlatform("basic.ts", "not-a-platform"),
				ExpectError: regexp.MustCompile(`unknown platform "not-a-platform"`),
			},
		},
	})
}

func testAccESBuildBundleBasicDataSourceConfig(filename string) string {
	return fmt.Sprintf(`
data "esbuild_bundle" "test" {
  filename = "${path.module}/fixtures/%s"
}
`, filename)
}

func testAccESBuildBundleBasicDataSourceConfigWithTarget(filename, target string) string {
	return fmt.Sprintf(`
data "esbuild_bundle" "test" {
  filename = "${path.module}/fixtures/%s"
  target   = "%s"
}
`, filename, target)
}

func testAccESBuildBundleBasicDataSourceConfigWithPlatform(filename, platform string) string {
	return fmt.Sprintf(`
data "esbuild_bundle" "test" {
  filename = "${path.module}/fixtures/%s"
  platform = "%s"
}
`, filename, platform)
}

const testAccESBuildBundleDefaultContentNode = `// fixtures/basic.ts
var x = "string";
console.log(x);
`

const testAccESBuildBundleDefaultContentBrowser = `(() => {
  // fixtures/basic.ts
  var x = "string";
  console.log(x);
})();
`


---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "esbuild_bundle Data Source - terraform-provider-esbuild"
subcategory: ""
description: |-
  A bundle of code compiled by ESBuild
---

# esbuild_bundle (Data Source)

A bundle of code compiled by ESBuild

## Example Usage

```terraform
data "esbuild_bundle" "node" {
  filename = "${path.module}/src/index.ts"
  platform = "node"
  target   = "es2020"
}

data "esbuild_bundle" "browser" {
  filename = "${path.module}/src/index.ts"
  platform = "browser"
}

# write the compiled content to a file
resource "local_file" "node" {
  filename = "dist/index.node.js"
  content  = data.esbuild_bundle.node.content
}

resource "local_file" "browser" {
  filename = "dist/index.browser.js"
  content  = data.esbuild_bundle.browser.content
}

# or just output them
output "content_for_node" {
  value = data.esbuild_bundle.node.content
}

output "content_for_browser" {
  value = data.esbuild_bundle.browser.content
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filename` (String) Path to the entrypoint file to be compiled and bundled

### Optional

- `platform` (String) The platform to compile the JavaScript code for
- `target` (String) The target environment to compile the JavaScript code for

### Read-Only

- `content` (String) The compiled content of the bundle
- `id` (String) SHA256 hash of the bundle content



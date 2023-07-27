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

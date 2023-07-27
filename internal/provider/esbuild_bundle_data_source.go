package provider

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ESBuildBundleDataSource{}

func NewESBuildBundleDataSource() datasource.DataSource {
	return &ESBuildBundleDataSource{}
}

// ESBuildBundleDataSource defines the data source implementation.
type ESBuildBundleDataSource struct {
}

// ESBuildBundleDataSourceModel describes the data source data model.
type ESBuildBundleDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Filename      types.String `tfsdk:"filename"`
	Platform      types.String `tfsdk:"platform"`
	Target        types.String `tfsdk:"target"`
	Content       types.String `tfsdk:"content"`
	SourceMapMode types.String `tfsdk:"sourcemap"`
}

func (d *ESBuildBundleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bundle"
}

func (d *ESBuildBundleDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A bundle of code compiled by ESBuild",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "SHA256 hash of the bundle content",
				Type:                types.StringType,
				Computed:            true,
			},
			"filename": {
				MarkdownDescription: "Path to the entrypoint file to be compiled and bundled",
				Type:                types.StringType,
				Required:            true,
			},
			"target": {
				MarkdownDescription: "The target environment to compile the JavaScript code for",
				Type:                types.StringType,
				Optional:            true,
			},
			"platform": {
				MarkdownDescription: "The platform to compile the JavaScript code for",
				Type:                types.StringType,
				Optional:            true,
			},
			"sourcemap": {
				MarkdownDescription: "The sourcemap generation setting",
				Type:                types.StringType,
				Optional:            true,
			},
			"content": {
				MarkdownDescription: "The compiled content of the bundle",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (d *ESBuildBundleDataSourceModel) getTarget() (api.Target, string, error) {
	if d.Target.IsNull() {
		return api.ESNext, "esnext", nil
	}

	validTargets := map[string]api.Target{
		"esnext": api.ESNext,
		"es5":    api.ES5,
		"es6":    api.ES2015,
		"es2015": api.ES2015,
		"es2016": api.ES2016,
		"es2017": api.ES2017,
		"es2018": api.ES2018,
		"es2019": api.ES2019,
		"es2020": api.ES2020,
		"es2021": api.ES2021,
		"es2022": api.ES2022,
	}

	if target, ok := validTargets[d.Target.Value]; ok {
		return target, d.Target.Value, nil
	}

	return api.DefaultTarget, "default", fmt.Errorf(`unknown target "%s"`, d.Target.Value)
}

func (d *ESBuildBundleDataSourceModel) getPlatform() (api.Platform, string, error) {
	if d.Platform.IsNull() {
		return api.PlatformNode, "node", nil
	}

	validPlatforms := map[string]api.Platform{
		"browser": api.PlatformBrowser,
		"node":    api.PlatformNode,
		"neutral": api.PlatformNeutral,
	}

	if platform, ok := validPlatforms[d.Platform.Value]; ok {
		return platform, d.Platform.Value, nil
	}

	return api.PlatformDefault, "default", fmt.Errorf(`unknown platform "%s"`, d.Platform.Value)
}

func (d *ESBuildBundleDataSourceModel) getSourceMapMode() (api.SourceMap, string, error) {
	if d.SourceMapMode.IsNull() {
		return api.SourceMapNone, "none", nil
	}

	validModes := map[string]api.SourceMap{
		"none":     api.SourceMapNone,
		"inline":   api.SourceMapInline,
		"linked":   api.SourceMapLinked,
		"external": api.SourceMapExternal,
		"both":     api.SourceMapInlineAndExternal,
	}

	if mode, ok := validModes[d.SourceMapMode.Value]; ok {
		return mode, d.SourceMapMode.Value, nil
	}

	return api.SourceMapNone, "none", fmt.Errorf(`unknown sourcemap mode "%s"`, d.SourceMapMode.Value)
}

func (d *ESBuildBundleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *ESBuildBundleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ESBuildBundleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	platform, platformAsString, err := data.getPlatform()

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")

		return
	}

	data.Platform = types.String{Value: platformAsString}

	target, targetAsString, err := data.getTarget()

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")

		return
	}

	data.Target = types.String{Value: targetAsString}

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{data.Filename.Value},
		Bundle:      true,
		Platform:    platform,
		Target:      target,
		Sourcemap:   api.SourceMapInline,
	})

	if len(result.Warnings) > 0 {
		for _, err := range result.Warnings {
			detail := ""

			if err.Location != nil {
				detail = fmt.Sprintf("in %s#L%d", err.Location.File, err.Location.Line)
			}

			resp.Diagnostics.AddWarning(err.Text, detail)
		}
	}

	if len(result.Errors) > 0 {
		for _, err := range result.Errors {
			detail := ""

			if err.Location != nil {
				detail = fmt.Sprintf("in %s#L%d", err.Location.File, err.Location.Line)
			}

			resp.Diagnostics.AddError(err.Text, detail)
		}

		return
	}

	if len(result.OutputFiles) != 1 {
		resp.Diagnostics.AddError(
			"unexpected number of files outputted",
			fmt.Sprintf("only 1 output file is expected, but got %d - please rpeort this!", len(result.OutputFiles)),
		)
	}

	file := result.OutputFiles[0]

	data.Id = types.String{Value: fmt.Sprintf("%x", sha256.Sum256(file.Contents))}
	data.Content = types.String{Value: string(file.Contents)}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

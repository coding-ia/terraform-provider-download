package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &DownloadProvider{}
var _ provider.ProviderWithFunctions = &DownloadProvider{}

type DownloadProvider struct {
	version string
}

func (d *DownloadProvider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "download"
	response.Version = d.version
}

func (d *DownloadProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "The download Terraform provider allows you to download a file from an http website.",
	}
}

func (d *DownloadProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
}

func (d *DownloadProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDownloadFileDataSource,
	}
}

func (d *DownloadProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewDownloadFileFunction,
	}
}

func (d *DownloadProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DownloadProvider{
			version: version,
		}
	}
}

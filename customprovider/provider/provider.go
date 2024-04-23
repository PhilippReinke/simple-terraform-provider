package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type customProviderModel struct {
	Host types.String `tfsdk:"host"`
}

var _ provider.Provider = &customProvider{} // check interface compliance

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &customProvider{
			version: version,
		}
	}
}

type customProvider struct {
	version string
}

func (p *customProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "customprovider"
	resp.Version = p.version
}

func (p *customProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *customProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// retrieve provider data from configuration
	var config customProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if practitioner provided a configuration value for any of the attributes,
	// it must be a known value
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown CustomProvider API Host",
			"The provider cannot create the CustomProvider API client as there is an unknown configuration value for the CustomProvider API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CUSTOMPROVIDERAPI_HOST environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// default values to environment variables, but override with Terraform
	// configuration value if set
	host := os.Getenv("CUSTOMPROVIDERAPI_HOST")
	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	// if any of the expected configurations are missing, return errors with
	// provider-specific guidance
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing CustomProvider API Host",
			"The provider cannot create the CustomProvider API client as there is a missing or empty value for the CustomProvider API host. "+
				"Set the host value in the configuration or use the CUSTOMPROVIDERAPI_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// make the CustomProvider client available during DataSource and Resource
	// type Configure methods
	resp.DataSourceData = host
	resp.ResourceData = host
}

func (p *customProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewItemsDataSource,
	}
}

func (p *customProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewItemResource,
	}
}

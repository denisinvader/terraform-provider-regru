package provider

import (
	"context"
	"os"
	"terraform-provider-regru/internal/regru"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RegruProvider struct {
	version string
}

type RegruProviderModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *RegruProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "regru"
	resp.Version = p.version
}

func (p *RegruProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *RegruProvider) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	var config RegruProviderModel

	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.Username.IsUnknown() {
		res.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown reg.ru API username",
			"The provider cannot create the reg.ru API client as there is an unknown configuration value for the username. Either target apply the source of the value first, set the value statically in the configuration, or use the REGRU_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		res.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown reg.ru API password",
			"The provider cannot create the reg.ru API client as there is an unknown configuration value for the password. Either target apply the source of the value first, set the value statically in the configuration, or use the REGRU_PASSWORD environment variable.",
		)
	}

	if res.Diagnostics.HasError() {
		return
	}

	username := os.Getenv("REGRU_USERNAME")
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	if username == "" {
		res.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing reg.ru API username",
			"The provider cannot create the reg.ru API client as there is a missing or empty value for the reg.ru API username. "+
				"Set the username value in the configuration or use the REGRU_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	password := os.Getenv("REGRU_PASSWORD")
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}
	if password == "" {
		res.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing reg.ru API password",
			"The provider cannot create the reg.ru API client as there is a missing or empty value for the reg.ru API password. "+
				"Set the password value in the configuration or use the REGRU_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if res.Diagnostics.HasError() {
		return
	}

	client := regru.NewClient(username, password)

	res.DataSourceData = &client
	res.ResourceData = &client
}

func (p *RegruProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceRecordA,
	}
}

func (p *RegruProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataSourceZone,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RegruProvider{
			version: version,
		}
	}
}

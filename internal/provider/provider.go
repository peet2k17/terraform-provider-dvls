package provider

import (
	"context"
	"os"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure DvlsProvider satisfies various provider interfaces.
var _ provider.Provider = &DvlsProvider{}

// DvlsProvider defines the provider implementation.
type DvlsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DvlsProviderModel describes the provider data model.
type DvlsProviderModel struct {
	BaseUri   types.String `tfsdk:"base_uri"`
	AppId     types.String `tfsdk:"app_id"`
	AppSecret types.String `tfsdk:"app_secret"`
}

func (p *DvlsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dvls"
	resp.Version = p.version
}

func (p *DvlsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The provider can be configured using the environment variables DVLS_APP_ID and DVLS_APP_SECRET",
		Attributes: map[string]schema.Attribute{
			"base_uri": schema.StringAttribute{
				Description: "DVLS base URI",
				Required:    true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "DVLS App ID `$DVLS_APP_ID`",
				Optional:            true,
			},
			"app_secret": schema.StringAttribute{
				MarkdownDescription: "DVLS App Secret `$DVLS_APP_SECRET`",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *DvlsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data DvlsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	baseuri := data.BaseUri.ValueString()
	appId := os.Getenv("DVLS_APP_ID")
	appSecret := os.Getenv("DVLS_APP_SECRET")

	if !data.AppId.IsNull() {
		appId = data.AppId.ValueString()
	}

	if !data.AppSecret.IsNull() {
		appSecret = data.AppSecret.ValueString()
	}

	if appId == "" || appSecret == "" {
		resp.Diagnostics.AddError("unable to set up dvls client", "'app_id' and 'app_secret' cannot be empty")
		return
	}

	dvlsClient, err := dvls.NewClient(appId, appSecret, baseuri)
	if err != nil {
		resp.Diagnostics.AddError("unable to set up dvls client", err.Error())
		return
	}

	resp.DataSourceData = &dvlsClient
	resp.ResourceData = &dvlsClient
}

func (p *DvlsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEntryUserCredentialResource,
		NewEntryCertificateResource,
		NewVaultResource,
	}
}

func (p *DvlsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEntryUserCredentialDataSource,
		NewEntryCertificateDataSource,
		NewEntryWebsiteDataSource,
		NewVaultDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DvlsProvider{
			version: version,
		}
	}
}

package provider

import (
	"context"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EntryWebsiteDataSource{}

func NewEntryWebsiteDataSource() datasource.DataSource {
	return &EntryWebsiteDataSource{}
}

// EntryWebsiteDataSource defines the resource implementation.
type EntryWebsiteDataSource struct {
	client *dvls.Client
}

// EntryWebsiteDataSourceModel describes the resource data model.
type EntryWebsiteDataSourceModel struct {
	Id                    types.String   `tfsdk:"id"`
	VaultId               types.String   `tfsdk:"vault_id"`
	Name                  types.String   `tfsdk:"name"`
	Description           types.String   `tfsdk:"description"`
	Username              types.String   `tfsdk:"username"`
	Password              types.String   `tfsdk:"password"`
	Url                   types.String   `tfsdk:"url"`
	Folder                types.String   `tfsdk:"folder"`
	Tags                  []types.String `tfsdk:"tags"`
	WebBrowserApplication types.Int64    `tfsdk:"web_browser_application"`
}

func (d *EntryWebsiteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_website"
}

func (d *EntryWebsiteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Website data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "User Credential ID",
				Required:    true,
				Validators:  []validator.String{entryWebsiteIdValidator{}},
			},
			"vault_id": schema.StringAttribute{
				Description: "Vault ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Website name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Website description",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Website username",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Website password",
				Computed:    true,
				Sensitive:   true,
			},
			"folder": schema.StringAttribute{
				Description: "Website folder path",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Website tags",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "Website URL",
				Computed:    true,
			},
			"web_browser_application": schema.Int64Attribute{
				Description: "Web browser application ID",
				Computed:    true,
			},
		},
	}
}

func (d *EntryWebsiteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dvls.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *dvls.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *EntryWebsiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EntryWebsiteDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryWebsite, err := d.client.Entries.Website.Get(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Website Entry",
			err.Error(),
		)
		return
	}

	entryWebsiteSensitiveData, err := d.client.Entries.Website.GetWebsiteDetails(entryWebsite)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Website Entry Sensitive Data",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(entryWebsite.ID)
	data.VaultId = types.StringValue(entryWebsite.VaultId)
	data.Name = types.StringValue(entryWebsite.EntryName)
	data.Description = types.StringValue(entryWebsite.Description)
	data.Username = types.StringValue(entryWebsiteSensitiveData.WebsiteDetails.Username)
	data.Password = types.StringValue(*entryWebsiteSensitiveData.WebsiteDetails.Password)
	data.Url = types.StringValue(entryWebsiteSensitiveData.WebsiteDetails.URL)
	data.Folder = types.StringValue(entryWebsite.EntryFolderPath)
	tags := make([]types.String, len(entryWebsite.Tags))
	for i, tag := range entryWebsite.Tags {
		tags[i] = types.StringValue(tag)
	}
	data.Tags = tags
	data.WebBrowserApplication = types.Int64Value(int64(entryWebsiteSensitiveData.WebsiteDetails.WebBrowserApplication))

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

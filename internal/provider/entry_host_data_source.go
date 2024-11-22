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
var _ datasource.DataSource = &EntryHostDataSource{}

func NewEntryHostDataSource() datasource.DataSource {
	return &EntryHostDataSource{}
}

// EntryHostDataSource defines the resource implementation.
type EntryHostDataSource struct {
	client *dvls.Client
}

// EntryHostDataSourceModel describes the resource data model.
type EntryHostDataSourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Username    types.String   `tfsdk:"username"`
	Password    types.String   `tfsdk:"password"`
	Host        types.String   `tfsdk:"host"`
	Folder      types.String   `tfsdk:"folder"`
	Tags        []types.String `tfsdk:"tags"`
}

func (d *EntryHostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_host"
}

func (d *EntryHostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Host data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "User Credential ID",
				Required:    true,
				Validators:  []validator.String{entryHostIdValidator{}},
			},
			"vault_id": schema.StringAttribute{
				Description: "Vault ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Host name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Host description",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Host username",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Host password",
				Computed:    true,
				Sensitive:   true,
			},
			"host": schema.StringAttribute{
				Description: "Host",
				Computed:    true,
			},
			"folder": schema.StringAttribute{
				Description: "Host folder path",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Host tags",
				Computed:    true,
			},
		},
	}
}

func (d *EntryHostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EntryHostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EntryHostDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryHost, err := d.client.Entries.Host.Get(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host Entry",
			err.Error(),
		)
		return
	}

	entryHostSensitiveData, err := d.client.Entries.Host.GetHostDetails(entryHost)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host Entry Sensitive Data",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(entryHost.ID)
	data.VaultId = types.StringValue(entryHost.VaultId)
	data.Name = types.StringValue(entryHost.EntryName)
	data.Description = types.StringValue(entryHost.Description)
	data.Username = types.StringValue(entryHostSensitiveData.HostDetails.Username)
	data.Password = types.StringValue(*entryHostSensitiveData.HostDetails.Password)
	data.Host = types.StringValue(entryHostSensitiveData.HostDetails.Host)
	data.Folder = types.StringValue(entryHost.EntryFolderPath)
	tags := make([]types.String, len(entryHost.Tags))
	for i, tag := range entryHost.Tags {
		tags[i] = types.StringValue(tag)
	}
	data.Tags = tags

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

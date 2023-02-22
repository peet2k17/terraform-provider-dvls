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
var _ datasource.DataSource = &EntryDataSource{}

func NewEntryDataSource() datasource.DataSource {
	return &EntryDataSource{}
}

// EntryDataSource defines the data source implementation.
type EntryDataSource struct {
	client *dvls.Client
}

// EntryDataSourceModel describes the data source data model.
type EntryDataSourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Username    types.String   `tfsdk:"username"`
	Password    types.String   `tfsdk:"password"`
	Folder      types.String   `tfsdk:"folder"`
	Tags        []types.String `tfsdk:"tags"`
}

func (d *EntryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry"
}

func (d *EntryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Entry data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Entry ID",
				Required:    true,
				Validators:  []validator.String{entryIdValidator{}},
			},
			"vault_id": schema.StringAttribute{
				Description: "Vault ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Entry name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Entry description",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Entry username",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Entry password",
				Computed:    true,
				Sensitive:   true,
			},
			"folder": schema.StringAttribute{
				Description: "Entry folder path",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Entry tags",
				Computed:    true,
			},
		},
	}
}

func (d *EntryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EntryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EntryDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := d.client.GetEntry(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to read entry", err.Error())
		return
	}

	entry, err = d.client.GetEntryCredentialsPassword(entry)
	if err != nil {
		resp.Diagnostics.AddError("unable to read entry sensitive information", err.Error())
		return
	}

	setEntryDataModel(entry, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

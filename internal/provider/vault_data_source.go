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
var _ datasource.DataSource = &VaultDataSource{}

func NewVaultDataSource() datasource.DataSource {
	return &VaultDataSource{}
}

// VaultDataSource defines the data source implementation.
type VaultDataSource struct {
	client *dvls.Client
}

// VaultDataSourceModel describes the data source data model.
type VaultDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	SecurityLevel types.String `tfsdk:"security_level"`
	Visibility    types.String `tfsdk:"visibility"`
}

func (d *VaultDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault"
}

func (d *VaultDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Vault data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Vault ID",
				Required:    true,
				Validators:  []validator.String{vaultIdValidator{}},
			},
			"name": schema.StringAttribute{
				Description: "Vault name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Vault description",
				Computed:    true,
			},
			"security_level": schema.StringAttribute{
				Description: "Vault security level",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Vault visibility",
				Computed:    true,
			},
		},
	}
}

func (d *VaultDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VaultDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *VaultDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vault, err := d.client.Vaults.Get(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to read vault", err.Error())
		return
	}

	setVaultDataModel(vault, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

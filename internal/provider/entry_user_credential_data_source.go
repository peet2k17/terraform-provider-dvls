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
var _ datasource.DataSource = &EntryUserCredentialDataSource{}

func NewEntryUserCredentialDataSource() datasource.DataSource {
	return &EntryUserCredentialDataSource{}
}

// EntryUserCredentialDataSource defines the data source implementation.
type EntryUserCredentialDataSource struct {
	client *dvls.Client
}

// EntryUserCredentialDataSourceModel describes the data source data model.
type EntryUserCredentialDataSourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Username    types.String   `tfsdk:"username"`
	Password    types.String   `tfsdk:"password"`
	Folder      types.String   `tfsdk:"folder"`
	Tags        []types.String `tfsdk:"tags"`
}

func (d *EntryUserCredentialDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_user_credential"
}

func (d *EntryUserCredentialDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User Credential data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "User Credential ID",
				Required:    true,
				Validators:  []validator.String{entryusercredentialIdValidator{}},
			},
			"vault_id": schema.StringAttribute{
				Description: "Vault ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User Credential name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "User Credential description",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "User Credential username",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "User Credential password",
				Computed:    true,
				Sensitive:   true,
			},
			"folder": schema.StringAttribute{
				Description: "User Credential folder path",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "User Credential tags",
				Computed:    true,
			},
		},
	}
}

func (d *EntryUserCredentialDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EntryUserCredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EntryUserCredentialDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryusercredential, err := d.client.Entries.UserCredential.Get(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to read user credential entry", err.Error())
		return
	}

	entryusercredential, err = d.client.Entries.UserCredential.GetUserAuthDetails(entryusercredential)
	if err != nil {
		resp.Diagnostics.AddError("unable to read user credential entry sensitive information", err.Error())
		return
	}

	setEntryUserCredentialDataModel(entryusercredential, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

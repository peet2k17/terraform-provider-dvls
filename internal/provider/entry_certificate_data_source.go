package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EntryCertificateDataSource{}

func NewEntryCertificateDataSource() datasource.DataSource {
	return &EntryCertificateDataSource{}
}

// EntryCertificateDataSource defines the data source implementation.
type EntryCertificateDataSource struct {
	client *dvls.Client
}

// EntryCertificateDataSourceModel describes the data source data model.
type EntryCertificateDataSourceModel struct {
	Id          types.String      `tfsdk:"id"`
	VaultId     types.String      `tfsdk:"vault_id"`
	Name        types.String      `tfsdk:"name"`
	Description types.String      `tfsdk:"description"`
	Password    types.String      `tfsdk:"password"`
	Folder      types.String      `tfsdk:"folder"`
	Url         types.Object      `tfsdk:"url"`
	File        types.Object      `tfsdk:"file"`
	Expiration  timetypes.RFC3339 `tfsdk:"expiration"`
	Tags        []types.String    `tfsdk:"tags"`
}

type EntryCertificateDataSourceModelData struct {
	Data *EntryCertificateDataSourceModel
	Url  *EntryCertificateDataSourceModelUrl
	File *EntryCertificateDataSourceModelFile
}

type EntryCertificateDataSourceModelUrl struct {
	Url                   types.String `tfsdk:"url"`
	UseDefaultCredentials types.Bool   `tfsdk:"use_default_credentials"`
}

func (m EntryCertificateDataSourceModelUrl) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"url":                     types.StringType,
		"use_default_credentials": types.BoolType,
	}
}

type EntryCertificateDataSourceModelFile struct {
	ContentB64 types.String `tfsdk:"content_b64"`
	Name       types.String `tfsdk:"name"`
}

func (m EntryCertificateDataSourceModelFile) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"content_b64": types.StringType,
		"name":        types.StringType,
	}
}

func (d *EntryCertificateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_certificate"
}

func (d *EntryCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Certificate data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Certificate ID",
				Required:    true,
				Validators:  []validator.String{entryCertificateIdValidator{}},
			},
			"vault_id": schema.StringAttribute{
				Description: "Vault ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Certificate name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Certificate description",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Certificate password",
				Computed:    true,
				Sensitive:   true,
			},
			"folder": schema.StringAttribute{
				Description: "Certificate folder path",
				Computed:    true,
			},

			"url": schema.SingleNestedAttribute{
				Description: "Certificate url. Either file or url must be specified.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description: "Certificate url",
						Computed:    true,
					},
					"use_default_credentials": schema.BoolAttribute{
						Description: "Use default credentials",
						Computed:    true,
					},
				},
			},

			"file": schema.SingleNestedAttribute{
				Description: "Certificate file. Either file or url must be specified.",
				Computed:    true,
				Sensitive:   true,

				Attributes: map[string]schema.Attribute{
					"content_b64": schema.StringAttribute{
						Description: "Certificate base 64 encoded string",
						Computed:    true,
						Sensitive:   true,
					},
					"name": schema.StringAttribute{
						Description: "Certificate file name",
						Computed:    true,
					},
				},
			},

			"expiration": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Description: "Certificate expiration date, in RFC3339 format (e.g. 2022-12-31T23:59:59-05:00)",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Certificate tags",
				Computed:    true,
			},
		},
	}
}

func (d *EntryCertificateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EntryCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EntryCertificateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entrycertificateId := data.Id.ValueString()

	entrycertificate, err := d.client.Entries.Certificate.Get(entrycertificateId)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read certificate entry", err.Error())
		return
	}

	entrycertificate, err = d.client.Entries.Certificate.GetPassword(entrycertificate)
	if err != nil {
		resp.Diagnostics.AddError("unable to read certificate entry sensitive information", err.Error())
		return
	}

	entryBytes, err := d.client.Entries.Certificate.GetFileContent(entrycertificate.ID)
	if err != nil {
		resp.Diagnostics.AddError("unable to read certificate entry content", err.Error())
		return
	}

	diagsModel := setEntryCertificateDataModel(ctx, entrycertificate, data, entryBytes)
	resp.Diagnostics.Append(diagsModel...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

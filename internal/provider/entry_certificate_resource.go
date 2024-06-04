package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EntryCertificateResource{}
var _ resource.ResourceWithImportState = &EntryCertificateResource{}

func NewEntryCertificateResource() resource.Resource {
	return &EntryCertificateResource{}
}

// EntryCertificateResource defines the resource implementation.
type EntryCertificateResource struct {
	client *dvls.Client
}

// EntryCertificateResourceModel describes the resource data model.
type EntryCertificateResourceModel struct {
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

type EntryCertificateResourceModelData struct {
	Data *EntryCertificateResourceModel
	Url  *EntryCertificateResourceModelUrl
	File *EntryCertificateResourceModelFile
}

type EntryCertificateResourceModelUrl struct {
	Url                   types.String `tfsdk:"url"`
	UseDefaultCredentials types.Bool   `tfsdk:"use_default_credentials"`
}

func (m EntryCertificateResourceModelUrl) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"url":                     types.StringType,
		"use_default_credentials": types.BoolType,
	}
}

type EntryCertificateResourceModelFile struct {
	ContentB64 types.String `tfsdk:"content_b64"`
	Name       types.String `tfsdk:"name"`
}

func (m EntryCertificateResourceModelFile) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"content_b64": types.StringType,
		"name":        types.StringType,
	}
}

func (r *EntryCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_certificate"
}

func (r *EntryCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Certificate",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Certificate ID",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"vault_id": schema.StringAttribute{
				Description:   "Vault ID",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description: "Certificate name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Certificate description",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Certificate password",
				Optional:    true,
				Sensitive:   true,
			},
			"folder": schema.StringAttribute{
				Description: "Certificate folder path",
				Optional:    true,
			},

			"url": schema.SingleNestedAttribute{
				Description:   "Certificate url. Either file or url must be specified.",
				Optional:      true,
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplaceIfConfigured()},

				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description: "Certificate url",
						Required:    true,
					},
					"use_default_credentials": schema.BoolAttribute{
						Description: "Use default credentials",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
				Validators: []validator.Object{objectvalidator.ExactlyOneOf(path.MatchRoot("file"))},
			},

			"file": schema.SingleNestedAttribute{
				Description:   "Certificate file. Either file or url must be specified.",
				Optional:      true,
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplaceIfConfigured()},
				Sensitive:     true,

				Attributes: map[string]schema.Attribute{
					"content_b64": schema.StringAttribute{
						Description: "Certificate base 64 encoded string",
						Required:    true,
						Sensitive:   true,
					},
					"name": schema.StringAttribute{
						Description: "Certificate file name",
						Required:    true,
					},
				},
				Validators: []validator.Object{objectvalidator.ExactlyOneOf(path.MatchRoot("url"))},
			},

			"expiration": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Description: "Certificate expiration date, in RFC3339 format (e.g. 2022-12-31T23:59:59-05:00)",
				Required:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Certificate tags",
				Optional:    true,
			},
		},
	}
}

func (r *EntryCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *EntryCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plans, diags := getPlans(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entrycertificate := newEntryCertificateFromResourceModel(&plans)

	entrycertificate = updateCertificateContent(plans, r.client, entrycertificate, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	entrycertificate, err := r.client.Entries.Certificate.GetPassword(entrycertificate)
	if err != nil {
		resp.Diagnostics.AddError("unable to read certificate entry sensitive information", err.Error())
		return
	}

	entryBytes, err := r.client.Entries.Certificate.GetFileContent(entrycertificate.ID)
	if err != nil {
		resp.Diagnostics.AddError("unable to read certificate entry content", err.Error())
		return
	}

	diagsModel := setEntryCertificateResourceModel(ctx, entrycertificate, plans.Data, entryBytes)
	resp.Diagnostics.Append(diagsModel...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plans.Data)...)
}

func (r *EntryCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	states, diags := getPlans(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entrycertificate := newEntryCertificateFromResourceModel(&states)

	entrycertificate, err := r.client.Entries.Certificate.Get(entrycertificate.ID)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read certificate entry", err.Error())
		return
	}

	entrycertificate, err = r.client.Entries.Certificate.GetPassword(entrycertificate)
	if err != nil {
		resp.Diagnostics.AddError("unable to read certificate entry sensitive information", err.Error())
		return
	}

	entryBytes, err := r.client.Entries.Certificate.GetFileContent(entrycertificate.ID)
	if err != nil {
		resp.Diagnostics.AddError("unable to read certificate entry content", err.Error())
		return
	}

	diagsModel := setEntryCertificateResourceModel(ctx, entrycertificate, states.Data, entryBytes)
	resp.Diagnostics.Append(diagsModel...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &states.Data)...)
}

func (r *EntryCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plans, diags := getPlans(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entrycertificate := newEntryCertificateFromResourceModel(&plans)

	_, err := r.client.Entries.Certificate.Update(entrycertificate)
	if err != nil {
		resp.Diagnostics.AddError("unable to update certificate entry", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plans.Data)...)
}

func (r *EntryCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryCertificateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Entries.Certificate.Delete(state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete certificate entry", err.Error())
		return
	}
}

func (r *EntryCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

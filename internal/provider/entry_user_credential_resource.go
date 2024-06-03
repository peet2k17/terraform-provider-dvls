package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EntryUserCredentialResource{}
var _ resource.ResourceWithImportState = &EntryUserCredentialResource{}

func NewEntryUserCredentialResource() resource.Resource {
	return &EntryUserCredentialResource{}
}

// EntryUserCredentialResource defines the resource implementation.
type EntryUserCredentialResource struct {
	client *dvls.Client
}

// EntryUserCredentialResourceModel describes the resource data model.
type EntryUserCredentialResourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Username    types.String   `tfsdk:"username"`
	Password    types.String   `tfsdk:"password"`
	Folder      types.String   `tfsdk:"folder"`
	Tags        []types.String `tfsdk:"tags"`
}

func (r *EntryUserCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_user_credential"
}

func (r *EntryUserCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS User Credential",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "User Credential ID",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"vault_id": schema.StringAttribute{
				Description:   "Vault ID",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description: "User Credential name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User Credential description",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "User Credential username",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "User Credential password",
				Optional:    true,
				Sensitive:   true,
			},
			"folder": schema.StringAttribute{
				Description: "User Credential folder path",
				Optional:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "User Credential tags",
				Optional:    true,
			},
		},
	}
}

func (r *EntryUserCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EntryUserCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EntryUserCredentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userDetails := r.client.Entries.UserCredential.NewUserAuthDetails(plan.Username.ValueString(), plan.Password.ValueString())
	entryusercredential := newEntryUserCredentialFromResourceModel(plan, userDetails)

	entryusercredential, err := r.client.Entries.UserCredential.New(entryusercredential)
	if err != nil {
		resp.Diagnostics.AddError("unable to create entryusercredential", err.Error())
		return
	}

	setEntryUserCredentialResourceModel(entryusercredential, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryUserCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntryUserCredentialResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userDetails := r.client.Entries.UserCredential.NewUserAuthDetails(state.Username.ValueString(), state.Password.ValueString())
	entryusercredential := newEntryUserCredentialFromResourceModel(state, userDetails)

	entryusercredential, err := r.client.Entries.UserCredential.Get(entryusercredential.ID)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read entryusercredential", err.Error())
		return
	}

	entryusercredential, err = r.client.Entries.UserCredential.GetUserAuthDetails(entryusercredential)
	if err != nil {
		resp.Diagnostics.AddError("unable to read entryusercredential sensitive information", err.Error())
		return
	}

	setEntryUserCredentialResourceModel(entryusercredential, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryUserCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EntryUserCredentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userDetails := r.client.Entries.UserCredential.NewUserAuthDetails(plan.Username.ValueString(), plan.Password.ValueString())
	entryusercredential := newEntryUserCredentialFromResourceModel(plan, userDetails)

	_, err := r.client.Entries.UserCredential.Update(entryusercredential)
	if err != nil {
		resp.Diagnostics.AddError("unable to update entryusercredential", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryUserCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryUserCredentialResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Entries.UserCredential.Delete(state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete entryusercredential", err.Error())
		return
	}
}

func (r *EntryUserCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

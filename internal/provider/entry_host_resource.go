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
var _ resource.Resource = &EntryHostResource{}
var _ resource.ResourceWithImportState = &EntryHostResource{}

func NewEntryHostResource() resource.Resource {
	return &EntryHostResource{}
}

// EntryHostResource defines the resource implementation.
type EntryHostResource struct {
	client *dvls.Client
}

// EntryHostResourceModel describes the resource data model.
type EntryHostResourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Hostname    types.String   `tfsdk:"hostname"`
	Username    types.String   `tfsdk:"username"`
	Password    types.String   `tfsdk:"password"`
	Folder      types.String   `tfsdk:"folder"`
	Tags        []types.String `tfsdk:"tags"`
}

func (r *EntryHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_host"
}

func (r *EntryHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Host Entry",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Host Entry ID",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"vault_id": schema.StringAttribute{
				Description:   "Vault ID",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description: "Host Entry Name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Host Entry Description",
			},
			"hostname": schema.StringAttribute{
				Description: "Host Entry Hostname",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Host Entry Username",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Host Entry Password",
				Optional:    true,
				Sensitive:   true,
			},
			"folder": schema.StringAttribute{
				Description: "Host Entry Folder",
			},
			"tags": schema.ListAttribute{
				Description: "Host Entry Tags",
				Elem:        schema.StringAttribute{},
			},
		},
	}
}

func (r *EntryHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get the resource model.
	model := &EntryHostResourceModel{}
	if err := req.Stash.Get(model); err != nil {
		resp.Error = err
		return
	}

	// Get the entry.
	entry, err := r.client.GetEntryHost(ctx, model.VaultId.String(), model.Id.String())
	if err != nil {
		resp.Error = err
		return
	}

	// Update the model.
	model.Name = types.String(entry.Name)
	model.Description = types.String(entry.Description)
	model.Hostname = types.String(entry.Hostname)
	model.Username = types.String(entry.Username)
	model.Password = types.String(entry.Password)
	model.Folder = types.String(entry.Folder)
	model.Tags = types.StringSlice(entry.Tags)

	// Set the model.
	resp.State = model
}

func (r *EntryHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get the resource model.
	model := &EntryHostResourceModel{}
	if err := req.Plan.Get(ctx, model); err != nil {
		resp.Error = err
		return
	}

	// Create the entry.
	entry, err := r.client.CreateEntryHost(ctx, model.VaultId.String(), model.Name.String(), model.Description.String(), model.Hostname.String(), model.Username.String(), model.Password.String(), model.Folder.String(), model.Tags.Strings())
	if err != nil {
		resp.Error = err
		return
	}

	// Update the model.
	model.Id = types.String(entry.Id)

	// Set the model.
	resp.State = model
}

func (r *EntryHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get the resource model.
	model := &EntryHostResourceModel{}
	if err := req.State.Get(ctx, model); err != nil {
		resp.Error = err
		return
	}

	// Update the entry.
	entry, err := r.client.UpdateEntryHost(ctx, model.VaultId.String(), model.Id.String(), model.Name.String(), model.Description.String(), model.Hostname.String(), model.Username.String(), model.Password.String(), model.Folder.String(), model.Tags.Strings())
	if err != nil {
		resp.Error = err
		return
	}

	// Update the model.
	model.Name = types.String(entry.Name)
	model.Description = types.String(entry.Description)
	model.Hostname = types.String(entry.Hostname)
	model.Username = types.String(entry.Username)
	model.Password = types.String(entry.Password)
	model.Folder = types.String(entry.Folder)
	model.Tags = types.StringSlice(entry.Tags)

	// Set the model.
	resp.State = model
}

func (r *EntryHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get the resource model.
	model := &EntryHostResourceModel{}
	if err := req.State.Get(ctx, model); err != nil {
		resp.Error = err
		return
	}

	// Delete the entry.
	if err := r.client.DeleteEntryHost(ctx, model.VaultId.String(), model.Id.String()); err != nil {
		resp.Error = err
		return
	}
}

func (r *EntryHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

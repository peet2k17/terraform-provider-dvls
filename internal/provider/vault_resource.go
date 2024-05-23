package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VaultResource{}
var _ resource.ResourceWithImportState = &VaultResource{}

func NewVaultResource() resource.Resource {
	return &VaultResource{}
}

// VaultResource defines the resource implementation.
type VaultResource struct {
	client *dvls.Client
}

// VaultResourceModel describes the resource data model.
type VaultResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	SecurityLevel  types.String `tfsdk:"security_level"`
	Visibility     types.String `tfsdk:"visibility"`
	MasterPassword types.String `tfsdk:"master_password"`
}

func (r *VaultResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault"
}

func (r *VaultResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Vault",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Vault ID",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "Vault name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Vault description",
				Optional:    true,
			},
			"security_level": schema.StringAttribute{
				Description: fmt.Sprintf("Vault security level. Must be one of the following: %s", listMapValues(vaultSecurityLevels)),
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("standard"),
				Validators:  []validator.String{vaultSecurityLevelValidator{}},
			},
			"visibility": schema.StringAttribute{
				Description: fmt.Sprintf("Vault visibility. Must be one of the following: %s", listMapValues(vaultVisibilities)),
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
				Validators:  []validator.String{vaultVisibilityValidator{}},
			},
			"master_password": schema.StringAttribute{
				Description: "Vault master password",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *VaultResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *VaultResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vault, err := newVaultFromResourceModel(plan)
	if err != nil {
		resp.Diagnostics.AddError("unable to create vault", err.Error())
		return
	}
	vault.ID = uuid.NewString()

	var options dvls.VaultOptions
	if !plan.MasterPassword.IsNull() {
		options.Password = plan.MasterPassword.ValueStringPointer()
	}

	err = r.client.NewVault(vault, &options)
	if err != nil {
		resp.Diagnostics.AddError("unable to create vault", err.Error())
		return
	}

	setVaultResourceModel(vault, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *VaultResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vault, err := newVaultFromResourceModel(state)
	if err != nil {
		resp.Diagnostics.AddError("unable to read vault", err.Error())
		return
	}

	vault, err = r.client.GetVault(vault.ID)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read vault", err.Error())
		return
	}

	setVaultResourceModel(vault, state)

	valid, err := r.client.ValidateVaultPassword(vault.ID, state.MasterPassword.ValueString())
	if err != nil && strings.Contains(err.Error(), "unexpected result code 0 (Error)") {
		state.MasterPassword = basetypes.NewStringNull()
	} else if err != nil {
		resp.Diagnostics.AddError("unable validate vault password", err.Error())
		return
	}

	if !valid {
		state.MasterPassword = basetypes.NewStringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *VaultResourceModel
	var state *VaultResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vault, err := newVaultFromResourceModel(plan)
	if err != nil {
		resp.Diagnostics.AddError("unable to update vault", err.Error())
		return
	}

	var options dvls.VaultOptions
	if !plan.MasterPassword.IsNull() {
		options.Password = plan.MasterPassword.ValueStringPointer()
	}

	err = r.client.UpdateVault(vault, &options)
	if err != nil {
		resp.Diagnostics.AddError("unable to update vault", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VaultResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *VaultResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVault(state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete vault", err.Error())
		return
	}
}

func (r *VaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

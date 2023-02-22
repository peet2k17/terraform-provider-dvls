package provider

import (
	"context"

	"github.com/Devolutions/go-dvls"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type entryIdValidator struct{}

func (validator entryIdValidator) Description(_ context.Context) string {
	return "entry must be a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)"
}

func (validator entryIdValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (d entryIdValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	id := request.ConfigValue.ValueString()

	_, err := uuid.Parse(id)
	if err != nil {
		response.Diagnostics.AddError("entry id is not a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)", err.Error())
		return
	}
}

func newEntryFromResourceModel(data *EntryResourceModel) dvls.Entry {
	var tags []string

	for _, v := range data.Tags {
		tags = append(tags, v.ValueString())
	}

	entry := dvls.Entry{
		ID:                data.Id.ValueString(),
		VaultId:           data.VaultId.ValueString(),
		EntryName:         data.Name.ValueString(),
		Description:       data.Description.ValueString(),
		Credentials:       dvls.NewEntryCredentials(data.Username.ValueString(), data.Password.ValueString()),
		EntryFolderPath:   data.Folder.ValueString(),
		ConnectionType:    dvls.ServerConnectionCredential,
		ConnectionSubType: dvls.ServerConnectionSubTypeDefault,
		Tags:              tags,
	}
	return entry
}

func setEntryResourceModel(entry dvls.Entry, data *EntryResourceModel) {
	var model EntryResourceModel

	model.Id = basetypes.NewStringValue(entry.ID)
	model.VaultId = basetypes.NewStringValue(entry.VaultId)
	model.Name = basetypes.NewStringValue(entry.EntryName)

	if entry.Credentials.Password != nil && *entry.Credentials.Password != "" {
		model.Password = basetypes.NewStringValue(*entry.Credentials.Password)
	}

	if entry.Description != "" {
		model.Description = basetypes.NewStringValue(entry.Description)
	}

	if entry.Credentials.Username != "" {
		model.Username = basetypes.NewStringValue(entry.Credentials.Username)
	}

	if entry.EntryFolderPath != "" {
		model.Folder = basetypes.NewStringValue(entry.EntryFolderPath)
	}

	if entry.Tags != nil {
		var tagsBase []types.String

		for _, v := range entry.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	*data = model
}

func setEntryDataModel(entry dvls.Entry, data *EntryDataSourceModel) {
	var model EntryDataSourceModel

	model.Id = basetypes.NewStringValue(entry.ID)
	model.VaultId = basetypes.NewStringValue(entry.VaultId)
	model.Name = basetypes.NewStringValue(entry.EntryName)

	if entry.Credentials.Password != nil && *entry.Credentials.Password != "" {
		model.Password = basetypes.NewStringValue(*entry.Credentials.Password)
	}

	if entry.Description != "" {
		model.Description = basetypes.NewStringValue(entry.Description)
	}

	if entry.Credentials.Username != "" {
		model.Username = basetypes.NewStringValue(entry.Credentials.Username)
	}

	if entry.EntryFolderPath != "" {
		model.Folder = basetypes.NewStringValue(entry.EntryFolderPath)
	}

	if entry.Tags != nil {
		var tagsBase []types.String

		for _, v := range entry.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	*data = model
}

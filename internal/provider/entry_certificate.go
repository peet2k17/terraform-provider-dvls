package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCertificateFromResourceModel(plans *EntryCertificateResourceModelData) dvls.EntryCertificate {
	var tags []string

	for _, v := range plans.Data.Tags {
		tags = append(tags, v.ValueString())
	}

	expiration, _ := plans.Data.Expiration.ValueRFC3339Time()

	entrycertificate := dvls.EntryCertificate{
		ID:              plans.Data.Id.ValueString(),
		VaultId:         plans.Data.VaultId.ValueString(),
		Name:            plans.Data.Name.ValueString(),
		Description:     plans.Data.Description.ValueString(),
		EntryFolderPath: plans.Data.Folder.ValueString(),
		Password:        plans.Data.Password.ValueString(),
		Expiration:      expiration,
		Tags:            tags,
	}

	if !plans.Data.File.IsNull() {
		entrycertificate.CertificateIdentifier = plans.File.Name.ValueString()
	} else if !plans.Data.Url.IsNull() {
		entrycertificate.CertificateIdentifier = plans.Url.Url.ValueString()
		entrycertificate.UseDefaultCredentials = plans.Url.UseDefaultCredentials.ValueBool()
	}

	return entrycertificate
}

func setEntryCertificateResourceModel(ctx context.Context, entrycertificate dvls.EntryCertificate, data *EntryCertificateResourceModel, content []byte) diag.Diagnostics {
	var diags diag.Diagnostics
	timeVal, timeDiags := timetypes.NewRFC3339Value(entrycertificate.Expiration.Format(time.RFC3339))
	diags.Append(timeDiags...)
	if diags.HasError() {
		return diags
	}

	model := EntryCertificateResourceModel{
		Id:         basetypes.NewStringValue(entrycertificate.ID),
		VaultId:    basetypes.NewStringValue(entrycertificate.VaultId),
		Name:       basetypes.NewStringValue(entrycertificate.Name),
		Expiration: timeVal,
		Url:        basetypes.NewObjectNull(EntryCertificateResourceModelUrl{}.AttributeTypes()),
		File:       basetypes.NewObjectNull(EntryCertificateResourceModelFile{}.AttributeTypes()),
	}

	switch entrycertificate.GetDataMode() {
	case dvls.EntryCertificateDataModeFile:
		fileObject := EntryCertificateResourceModelFile{
			ContentB64: basetypes.NewStringValue(base64.StdEncoding.EncodeToString(content)),
			Name:       basetypes.NewStringValue(entrycertificate.CertificateIdentifier),
		}

		objectValue, objDiags := types.ObjectValueFrom(ctx, fileObject.AttributeTypes(), fileObject)
		diags.Append(objDiags...)
		if diags.HasError() {
			return diags
		}

		model.File = objectValue
	case dvls.EntryCertificateDataModeURL:
		urlObject := EntryCertificateResourceModelUrl{
			Url:                   basetypes.NewStringValue(entrycertificate.CertificateIdentifier),
			UseDefaultCredentials: basetypes.NewBoolValue(entrycertificate.UseDefaultCredentials),
		}

		objectValue, objDiags := types.ObjectValueFrom(ctx, urlObject.AttributeTypes(), urlObject)
		diags.Append(objDiags...)
		if diags.HasError() {
			return diags
		}

		model.Url = objectValue
	default:
		diags.AddError("unable to set certificate entry", fmt.Sprintf("unknown data mode %d. Should be 2 for files or 3 for url", entrycertificate.GetDataMode()))
	}

	if entrycertificate.Password != "" {
		model.Password = basetypes.NewStringValue(entrycertificate.Password)
	}

	if entrycertificate.Description != "" {
		model.Description = basetypes.NewStringValue(entrycertificate.Description)
	}

	if entrycertificate.EntryFolderPath != "" {
		model.Folder = basetypes.NewStringValue(entrycertificate.EntryFolderPath)
	}

	if entrycertificate.Tags != nil {
		var tagsBase []types.String

		for _, v := range entrycertificate.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	*data = model

	return diags
}

func setEntryCertificateDataModel(ctx context.Context, entrycertificate dvls.EntryCertificate, data *EntryCertificateDataSourceModel, content []byte) diag.Diagnostics {
	var diags diag.Diagnostics
	timeVal, timeDiags := timetypes.NewRFC3339Value(entrycertificate.Expiration.Format(time.RFC3339))
	diags.Append(timeDiags...)
	if diags.HasError() {
		return diags
	}

	model := EntryCertificateDataSourceModel{
		Id:         basetypes.NewStringValue(entrycertificate.ID),
		VaultId:    basetypes.NewStringValue(entrycertificate.VaultId),
		Name:       basetypes.NewStringValue(entrycertificate.Name),
		Expiration: timeVal,
		Url:        basetypes.NewObjectNull(EntryCertificateResourceModelUrl{}.AttributeTypes()),
		File:       basetypes.NewObjectNull(EntryCertificateResourceModelFile{}.AttributeTypes()),
	}

	switch entrycertificate.GetDataMode() {
	case dvls.EntryCertificateDataModeFile:
		fileObject := EntryCertificateResourceModelFile{
			ContentB64: basetypes.NewStringValue(base64.StdEncoding.EncodeToString(content)),
			Name:       basetypes.NewStringValue(entrycertificate.CertificateIdentifier),
		}

		objectValue, objDiags := types.ObjectValueFrom(ctx, fileObject.AttributeTypes(), fileObject)
		diags.Append(objDiags...)
		if diags.HasError() {
			return diags
		}

		model.File = objectValue
	case dvls.EntryCertificateDataModeURL:
		urlObject := EntryCertificateResourceModelUrl{
			Url:                   basetypes.NewStringValue(entrycertificate.CertificateIdentifier),
			UseDefaultCredentials: basetypes.NewBoolValue(entrycertificate.UseDefaultCredentials),
		}

		objectValue, objDiags := types.ObjectValueFrom(ctx, urlObject.AttributeTypes(), urlObject)
		diags.Append(objDiags...)
		if diags.HasError() {
			return diags
		}

		model.Url = objectValue
	default:
		diags.AddError("unable to set certificate entry", fmt.Sprintf("unknown data mode %d. Should be 2 for files or 3 for url", entrycertificate.GetDataMode()))
	}

	if entrycertificate.Password != "" {
		model.Password = basetypes.NewStringValue(entrycertificate.Password)
	}

	if entrycertificate.Description != "" {
		model.Description = basetypes.NewStringValue(entrycertificate.Description)
	}

	if entrycertificate.EntryFolderPath != "" {
		model.Folder = basetypes.NewStringValue(entrycertificate.EntryFolderPath)
	}

	if entrycertificate.Tags != nil {
		var tagsBase []types.String

		for _, v := range entrycertificate.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	*data = model

	return diags
}

type planInterface interface {
	Get(ctx context.Context, target interface{}) diag.Diagnostics
}

func getPlans(ctx context.Context, plan planInterface) (EntryCertificateResourceModelData, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model *EntryCertificateResourceModel
	var urlPlan *EntryCertificateResourceModelUrl
	var filePlan *EntryCertificateResourceModelFile

	diags.Append(plan.Get(ctx, &model)...)
	if diags.HasError() {
		return EntryCertificateResourceModelData{}, diags
	}

	diags.Append(model.File.As(ctx, &filePlan, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return EntryCertificateResourceModelData{}, diags
	}

	diags.Append(model.Url.As(ctx, &urlPlan, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return EntryCertificateResourceModelData{}, diags
	}

	return EntryCertificateResourceModelData{
		Data: model,
		File: filePlan,
		Url:  urlPlan,
	}, diags
}

func updateCertificateContent(plans EntryCertificateResourceModelData, client *dvls.Client, entrycertificate dvls.EntryCertificate, diags *diag.Diagnostics) dvls.EntryCertificate {
	var err error

	if !plans.Data.File.IsNull() {
		content, err := base64.StdEncoding.DecodeString(plans.File.ContentB64.ValueString())
		if err != nil {
			diags.AddError("unable to update certificate entry", err.Error())
			return dvls.EntryCertificate{}
		}

		entrycertificate, err = client.Entries.Certificate.NewFile(entrycertificate, content)
		if err != nil {
			diags.AddError("unable to update certificate entry", err.Error())
			return dvls.EntryCertificate{}
		}
	} else {
		entrycertificate, err = client.Entries.Certificate.NewURL(entrycertificate)
		if err != nil {
			diags.AddError("unable to update certificate entry", err.Error())
			return dvls.EntryCertificate{}
		}
	}

	return entrycertificate
}

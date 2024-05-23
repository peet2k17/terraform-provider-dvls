package provider

import (
	"fmt"
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var vaultSecurityLevels map[dvls.VaultSecurityLevel]string = map[dvls.VaultSecurityLevel]string{
	dvls.VaultSecurityLevelStandard: "standard",
	dvls.VaultSecurityLevelHigh:     "high",
}

var vaultVisibilities map[dvls.VaultVisibility]string = map[dvls.VaultVisibility]string{
	dvls.VaultVisibilityDefault: "default",
	dvls.VaultVisibilityPublic:  "public",
	dvls.VaultVisibilityPrivate: "private",
}

func newVaultFromResourceModel(data *VaultResourceModel) (dvls.Vault, error) {
	securityLevel, err := lookupMapValue(vaultSecurityLevels, data.SecurityLevel.ValueString())
	if err != nil {
		return dvls.Vault{}, err
	}

	visibility, err := lookupMapValue(vaultVisibilities, data.Visibility.ValueString())
	if err != nil {
		return dvls.Vault{}, err
	}

	vault := dvls.Vault{
		ID:            data.Id.ValueString(),
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueString(),
		Visibility:    dvls.VaultVisibility(visibility),
		SecurityLevel: dvls.VaultSecurityLevel(securityLevel),
	}

	return vault, nil
}

func setVaultResourceModel(vault dvls.Vault, data *VaultResourceModel) {
	model := VaultResourceModel{
		Id:             basetypes.NewStringValue(vault.ID),
		Name:           basetypes.NewStringValue(vault.Name),
		Visibility:     basetypes.NewStringValue(vaultVisibilities[vault.Visibility]),
		SecurityLevel:  basetypes.NewStringValue(vaultSecurityLevels[vault.SecurityLevel]),
		MasterPassword: data.MasterPassword,
	}

	if vault.Description != "" {
		model.Description = basetypes.NewStringValue(vault.Description)
	}

	*data = model
}

func setVaultDataModel(vault dvls.Vault, data *VaultDataSourceModel) {
	model := VaultDataSourceModel{
		Id:            basetypes.NewStringValue(vault.ID),
		Name:          basetypes.NewStringValue(vault.Name),
		Visibility:    basetypes.NewStringValue(vaultVisibilities[vault.Visibility]),
		SecurityLevel: basetypes.NewStringValue(vaultSecurityLevels[vault.SecurityLevel]),
	}

	if vault.Description != "" {
		model.Description = basetypes.NewStringValue(vault.Description)
	}

	*data = model
}

func lookupMapValue[K interface {
	int | dvls.VaultSecurityLevel | dvls.VaultVisibility
}](lookup map[K]string, value string) (int, error) {
	for k, v := range lookup {
		if v == value {
			return int(k), nil
		}
	}

	return 0, fmt.Errorf("value %s not found in lookup", value)
}

func listMapValues[K interface {
	int | dvls.VaultSecurityLevel | dvls.VaultVisibility
}](lookup map[K]string) string {
	var values []string
	for _, v := range lookup {
		values = append(values, v)
	}

	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type vaultIdValidator struct{}

func (validator vaultIdValidator) Description(_ context.Context) string {
	return "vault must be a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)"
}

func (validator vaultIdValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (d vaultIdValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	id := request.ConfigValue.ValueString()

	_, err := uuid.Parse(id)
	if err != nil {
		response.Diagnostics.AddError("vault id is not a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)", err.Error())
		return
	}
}

type vaultSecurityLevelValidator struct{}

func (validator vaultSecurityLevelValidator) Description(_ context.Context) string {
	values := listMapValues(vaultSecurityLevels)
	return fmt.Sprintf("valid values are: %v", values)
}

func (validator vaultSecurityLevelValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (d vaultSecurityLevelValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	securityLevel := request.ConfigValue.ValueString()

	_, err := lookupMapValue(vaultSecurityLevels, securityLevel)
	if err != nil {
		values := listMapValues(vaultSecurityLevels)
		response.Diagnostics.AddError("vault security level is invalid", fmt.Sprintf("valid values are: %s", values))
		return
	}
}

type vaultVisibilityValidator struct{}

func (validator vaultVisibilityValidator) Description(_ context.Context) string {
	values := listMapValues(vaultVisibilities)
	return fmt.Sprintf("valid values are: %v", values)
}

func (validator vaultVisibilityValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (d vaultVisibilityValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	visibility := request.ConfigValue.ValueString()

	_, err := lookupMapValue(vaultVisibilities, visibility)
	if err != nil {
		values := listMapValues(vaultVisibilities)
		response.Diagnostics.AddError("vault visibility is invalid", fmt.Sprintf("valid values are: %s", values))
		return
	}
}

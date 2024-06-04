package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type entryusercredentialIdValidator struct{}

func (validator entryusercredentialIdValidator) Description(_ context.Context) string {
	return "user credential entry must be a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)"
}

func (validator entryusercredentialIdValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (d entryusercredentialIdValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	id := request.ConfigValue.ValueString()

	_, err := uuid.Parse(id)
	if err != nil {
		response.Diagnostics.AddError("user credential entry id is not a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)", err.Error())
		return
	}
}

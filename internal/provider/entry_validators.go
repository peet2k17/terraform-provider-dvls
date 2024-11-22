package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type entryusercredentialIdValidator struct{}
type entryCertificateIdValidator struct{}
type entryHostIdValidator struct{}
type entryWebsiteIdValidator struct{}

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

func (validator entryCertificateIdValidator) Description(_ context.Context) string {
	return "certificate entry must be a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)"
}

func (validator entryCertificateIdValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (d entryCertificateIdValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	id := request.ConfigValue.ValueString()

	_, err := uuid.Parse(id)
	if err != nil {
		response.Diagnostics.AddError("certificate entry id is not a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)", err.Error())
		return
	}
}

func (validator entryHostIdValidator) Description(_ context.Context) string {
	return "host entry must be a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)"
}

func (validator entryHostIdValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (d entryHostIdValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	id := request.ConfigValue.ValueString()

	_, err := uuid.Parse(id)
	if err != nil {
		response.Diagnostics.AddError("host entry id is not a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)", err.Error())
		return
	}
}

func (validator entryWebsiteIdValidator) Description(_ context.Context) string {
	return "website entry must be a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)"
}

func (validator entryWebsiteIdValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (d entryWebsiteIdValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	id := request.ConfigValue.ValueString()

	_, err := uuid.Parse(id)
	if err != nil {
		response.Diagnostics.AddError("website entry id is not a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)", err.Error())
		return
	}
}

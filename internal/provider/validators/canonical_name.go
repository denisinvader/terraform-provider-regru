package validators

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ validator.String = canonicalNameValidator{}
)

func CanonicalName() canonicalNameValidator {
	return canonicalNameValidator{}
}

type canonicalNameValidator struct{}

func (v canonicalNameValidator) Description(_ context.Context) string {
	return "must be a subdomain alias or FQDN and not target TLD"
}

func (v canonicalNameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v canonicalNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	firstDotIndex := strings.Index(value, ".")

	if firstDotIndex == -1 {
		return
	}

	if firstDotIndex < len(value)-1 {
		return
	}

	res.Diagnostics.Append(createDiagnosticError(req.Path, "Invalid canonical name value", v.Description(ctx), value))
}

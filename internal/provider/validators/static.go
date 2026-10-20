package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func StringWithoutWhitespace() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^[^\s]+$`),
		"must not contain any whitespace",
	)
}

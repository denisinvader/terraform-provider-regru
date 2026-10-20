package validators

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func createDiagnosticError(path path.Path, title, validatorDescription, value string) diag.DiagnosticWithPath {
	return diag.NewAttributeErrorDiagnostic(
		path,
		title,
		fmt.Sprintf("Attribute %s %s, got: %s", path, validatorDescription, value),
	)
}

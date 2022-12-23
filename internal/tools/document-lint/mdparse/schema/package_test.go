package schema_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

func TestPackage(t *testing.T) {
	pack := schema.NewPackage("github.com/hashicorp/terraform-provider-azurerm/tools")
	//pack := schema.NewPackage("github.com/hashicorp/terraform-provider-azurerm/internal/services/automation")
	t.Logf("got pack: %v", util.Stringify(pack, true))
}

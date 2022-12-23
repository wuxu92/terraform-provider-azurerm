package schema_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

func TestResourceFile(t *testing.T) {
	p := automation.SoftwareUpdateConfigurationResource{}
	//w := sdk.NewResourceWrapper(p)
	//res, _ := w.Resource()
	file := schema.FileForResource(p.Read().Func)
	t.Logf(file)

	// inspect schema
	r := schema.NewResourceByTyped(p)
	//name := r.ValidateFor("allocation_method")

	//packs := r.LoadImports()
	//for name, pack := range packs {
	//	t.Logf("%s: %v", name, pack)
	//}
	//t.Logf("%s", name)

	r.FindAllInSliceProp()
	//t.Logf("%s", util.Stringify(r, true))
	t.Logf("%s", util.Stringify(r.PossibleValues, true))
}

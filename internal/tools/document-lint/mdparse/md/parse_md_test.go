package md

import (
	"path/filepath"
	"testing"
)

func Test_unmarshalFile(t *testing.T) {
	dir := "/home/wuxu/go/src/github.com/terraform-provider-azurerm/website/docs/r"
	file := filepath.Join(dir, "automation_account.html.markdown")
	m := mustNewMarkFromFile(file)
	if len(m.Items) != 49 {
		t.Fatal(len(m.Items))
	}
}

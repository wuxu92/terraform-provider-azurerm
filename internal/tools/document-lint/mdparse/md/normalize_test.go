package md

import (
	"os"
	"path"
	"testing"
)

func Test_fixFileNormalize(t *testing.T) {
	dir, err := os.ReadDir(ResourceDir())
	_ = err
	for _, en := range dir {
		if en.IsDir() {
			continue
		}
		fullPath := path.Join(ResourceDir(), en.Name())
		fixFileNormalize(fullPath)
	}
}

func TestMDFile(t *testing.T) {
	file := "automation_watcher.html.markdown"
	fixFileNormalize("/home/wuxu/go/src/github.com/terraform-provider-azurerm/website/docs/r/" + file)
}

func TestRegSubMatch(t *testing.T) {
	idx := oldBlockHeadReg.FindStringSubmatchIndex("`traffic_analytics` supports the following:")
	t.Logf("%v", idx)

	for _, val := range []string{
		"  * `abc`  def",
		"* `abc` -  something  here.  ",
	} {
		res := removeRedundantSpace(val)
		t.Logf("from `%s` => `%s`", val, res)
	}
}

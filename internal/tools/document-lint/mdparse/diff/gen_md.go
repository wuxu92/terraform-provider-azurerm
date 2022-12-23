package diff

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

func (d *DiffResult) CrossCheckIssues() string {
	var bs strings.Builder
	var sf = fmt.Sprintf

	bs.WriteString(`# All missed properties

## Examples

1. ` + "`" + `a.b.c missed in doc` + "`" + `, means the property defined in schema code but not presents in documents.
2. ` + "`" + `a.b.c missed in doc` + "`" + ` with another ` + "`" + `a.b.c1 missed in code` + "`" + `. when two similar property missed in both doc and code, there should be typo in document's property name.
3. ` + "`" + `a.b not block missed in doc` + "`" + ` a block missed in the document, or there are malformed syntax the tool cannot parse, just skip it if so.
4. ` + "`" + `:L23 missed in doc` + "`" + ` there is a syntax error in line 23 of the document cause the tool can not parse it.
5.` + "`" + `a.b deprecated miss in code` + "`" + `property not exists in code schema but still in document and mark as deprecated, may delete from document?.

`)
	// group by owner
	byOwner := map[string]*strings.Builder{}
	for cat, resources := range d.byCate {
		var subBs strings.Builder
		owner := util.GetRPOwner(cat)
		if byOwner[owner] == nil {
			byOwner[owner] = &strings.Builder{}
		}
		var hasMiss bool
		subBs.WriteString(sf("## %s\n\n", cat))
		for _, resource := range resources {
			var sub3 strings.Builder
			sub3.WriteString(sf("### %s\n\n", resource.tf.ResourceType))
			var hasSub3 bool
			sort.Slice(resource.Diff, func(i, j int) bool {
				return resource.Diff[i].Key < resource.Diff[j].Key
			})
			for _, item := range resource.Diffs() {
				if item.MissType > 0 {
					hasSub3 = true
					sub3.WriteString(sf("- [ ] %s missed in %s\n", item.Key, item.MissType.String()))
				}
			}
			if hasSub3 {
				hasMiss = true
				subBs.WriteString(sub3.String())
				subBs.WriteString("\n")
			}
		}
		if hasMiss {
			bs.WriteString(subBs.String())
			//bs.WriteString("\n")
			byOwner[owner].WriteString(subBs.String())
		}
	}

	writeByOwner(byOwner)
	return bs.String()
}

func writeByOwner(data map[string]*strings.Builder) {
	f, err := os.OpenFile("diff_by_owner.md", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("open file by owner err: %v", err)
		return
	}
	sf := fmt.Sprintf
	var keys []string
	for own, _ := range data {
		keys = append(keys, own)
	}
	sort.Strings(keys)
	for _, own := range keys {
		f.WriteString(sf("# %s\n\n", own))
		f.WriteString(data[own].String())
	}
	f.Sync()
	f.Close()
}

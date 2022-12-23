package check

import (
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/md"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/schema"
)

// logic to load schema and markdown to print the diff

type ResourceDiff struct {
	tf *schema.Resource
	md *model.ResourceDoc

	SchemaFile string
	MDFile     string

	Diff []Checker // diff of exist in both md and code
}

func (d *ResourceDiff) Diffs() []Checker {
	return d.Diff
}

func (d *ResourceDiff) ToString() string {
	var bs strings.Builder

	bs.WriteString(
		fmt.Sprintf("%s : diff size: %d, document file: %s\n",
			d.tf.ResourceType,
			len(d.Diff),
			path.Base(d.MDFile),
		),
	)
	file := d.MDFile + ":"
	if idx := strings.Index(file, "website"); idx > 0 {
		file = "./" + file[idx:]
	}
	for _, item := range d.Diffs() {
		bs.WriteString(file + item.String() + "\n")
	}
	return bs.String()
}

// NewResourceDiff tf is required,
// mdPath is optional, it can be detected by resource name
func NewResourceDiff(tf *schema.Resource) *ResourceDiff {
	r := &ResourceDiff{
		tf: tf,
	}
	// try to detect Markdown path from resource
	// can set it if not a regular MD path
	r.MDFile = md.MDPathFor(tf.ResourceType)
	return r
}

func (r *ResourceDiff) DiffAll() {
	if r.md == nil {
		r.md, _ = md.UnmarshalResourceFromFile(r.MDFile)
	}
	r.Diff = checkPossibleValues(r.tf, r.md)

	missDiff := crossCheckProperty(r.tf, r.md)
	r.Diff = append(r.Diff, missDiff...)

	timeouts := diffTimeout(r.tf, r.md)
	r.Diff = append(r.Diff, timeouts...)
}

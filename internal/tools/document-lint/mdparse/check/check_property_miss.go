package check

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
)

type MissType int

func (m MissType) String() string {
	return []string{"ok", "doc", "code", "doc attribute"}[m]
}

const (
	NotMiss MissType = iota
	MissInDoc
	MissInCode
	MissInDocAttr
)

type propertyMissDiff struct {
	checkBase
	MissType MissType
}

func newPropertyMiss(checkBase checkBase, missType MissType) *propertyMissDiff {
	return &propertyMissDiff{checkBase: checkBase, MissType: missType}
}

func (c propertyMissDiff) String() string {
	return fmt.Sprintf("%s missed in %s", c.checkBase.Str(), c.MissType)
}

func (c propertyMissDiff) Fix(line string) (result string, err error) {
	return line, nil
}

var _ Checker = (*propertyMissDiff)(nil)

func newMissItem(path string, f *model.Field, typ MissType) Checker {
	base := newCheckBase(0, path, f)
	if f != nil {
		base.line = f.Line
	}
	return newPropertyMiss(base, typ)
}

func newMissInCode(path string, f *model.Field) Checker {
	return newMissItem(path, f, MissInCode)
}

func newMissInDoc(path string, f *model.Field) Checker {
	return newMissItem(path, f, MissInDoc)
}

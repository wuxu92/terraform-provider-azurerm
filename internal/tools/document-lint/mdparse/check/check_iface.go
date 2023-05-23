package check

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
)

type Checker interface {
	Line() int
	Key() string // property key path
	ShouldSkip() bool
	MDField() *model.Field

	// String display diff item information for check
	String() string

	// Fix try to fix this issue with line. return the updated line
	Fix(line string) (result string, err error)
}

type checkBase struct {
	line    int
	key     string
	mdField *model.Field
}

func (c checkBase) ShouldSkip() bool {
	if c.MDField() == nil || c.MDField().Skip {
		return true
	}
	return false
}

func (i checkBase) Str() string {
	return fmt.Sprintf("%d %s", i.Line(), i.Key())
}

func (i checkBase) Line() int {
	return i.line
}

func (i checkBase) Key() string {
	return i.key
}

func (i checkBase) MDField() *model.Field {
	return i.mdField
}

func newCheckBase(line int, key string, mdField *model.Field) checkBase {
	return checkBase{
		line:    line,
		key:     key,
		mdField: mdField,
	}
}

package check

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

type formatErr struct {
	Origin string
	checkBase
}

func newFormatErr(origin string, checkBase checkBase) *formatErr {
	return &formatErr{Origin: origin, checkBase: checkBase}
}

func (f formatErr) String() string {
	return fmt.Sprintf("%s should be formatted as: %s.",
		f.checkBase.Str(),
		util.FormatCode("* `field` - (Required/Optional) Xxx..."),
	)
}

func (f formatErr) Fix(line string) (result string, err error) {
	// some Note lines with a misleading star mark, try to remove it
	if strings.HasPrefix(line, "* ~>") {
		line = strings.TrimPrefix(line, "* ")
	}
	return line, nil // no fix for format error
}

var _ Checker = (*formatErr)(nil)

package check

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

type formatErr struct {
	Origin string
	msg    string
	checkBase
}

func newFormatErr(origin, msg string, checkBase checkBase) *formatErr {
	return &formatErr{
		Origin:    origin,
		msg:       msg,
		checkBase: checkBase,
	}
}

func (f formatErr) String() string {
	if strings.Contains(f.msg, "missing block for") {
		return fmt.Sprintf("%s %s", f.checkBase.Str(), util.IssueLine(f.msg))
	}
	if strings.Contains(f.msg, "duplicate") {
		return fmt.Sprintf("%s %s", f.checkBase.Str(), util.IssueLine(f.msg))
	}
	return fmt.Sprintf("%s should be formatted as: %s or '%s'",
		f.checkBase.Str(),
		util.FormatCode("* `field` - (Required/Optional) Xxx..."),
		f.msg,
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

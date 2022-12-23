package check

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/md"
)

type defaultDiff struct {
	checkBase
	Current string // the value in current document
	Default string // the right default value
}

func newDefaultDiff(checkBase checkBase, current string, defaultVaulue string) *defaultDiff {
	return &defaultDiff{checkBase: checkBase, Current: current, Default: defaultVaulue}
}

func (c defaultDiff) String() string {
	return fmt.Sprintf("%s default value should be `%s`", c.checkBase.Str(), c.Default)
}

// Fix if value is "", then we should remove the default value part from the document
func (c defaultDiff) Fix(line string) (result string, err error) {
	if idxs := md.DefaultsReg.FindStringSubmatchIndex(line); len(idxs) > 2 {
		if c.Default == "" {
			// remove default part from line
			line = line[:idxs[0]+1] + line[idxs[1]:]
		} else {
			line = line[:idxs[2]] + c.Default + line[idxs[3]:]
		}
	} else {
		line = strings.TrimSpace(line) + " Defaults to `" + c.Default + "`."
	}
	return line, nil
}

var _ Checker = (*defaultDiff)(nil)

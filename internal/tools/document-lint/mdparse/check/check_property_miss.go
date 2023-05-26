package check

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

type MissType int

func (m MissType) String() string {
	return []string{"ok", "doc", "code", "doc attribute", "block declare"}[m]
}

const (
	NotMiss MissType = iota
	MissInDoc
	MissInCode
	MissInDocAttr
	MissBlockDeclare // document block declare wrong-formatted
	Misspelling
)

type propertyMissDiff struct {
	checkBase
	MissType    MissType
	correctName string // for misspelling diff only
}

func newPropertyMiss(checkBase checkBase, missType MissType) *propertyMissDiff {
	return &propertyMissDiff{checkBase: checkBase, MissType: missType}
}

func (c propertyMissDiff) String() string {
	if c.MissType == MissBlockDeclare {
		return fmt.Sprintf("%s block should be declared like '%s'", c.checkBase.Str(), util.ItalicCode("A `xxx` block as defined below."))
	}
	if c.MissType == Misspelling {
		return fmt.Sprintf("%s with field name misspelling, should be %s?", c.checkBase.Str(), util.FixedCode(c.correctName))
	}
	return fmt.Sprintf("%s not exist in %s", c.checkBase.Str(), c.MissType)
}

func (c propertyMissDiff) Fix(line string) (result string, err error) {
	if c.MissType == Misspelling && c.correctName != "" {
		return strings.ReplaceAll(line, c.mdField.Name, c.correctName), nil
	}
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

// miss in doc will fill a mock `f`
func newMissInDoc(path string, f *model.Field) Checker {
	return newMissItem(path, f, MissInDoc)
}

func newMissBlockDeclare(path string, f *model.Field) Checker {
	return newMissItem(path, f, MissBlockDeclare)
}

func newMisspelling(c *propertyMissDiff, d *propertyMissDiff) Checker {
	item := newPropertyMiss(c.checkBase, Misspelling)
	item.correctName = d.MDField().Name
	return item
}

// missing in doc/code can be a misspelling in document. do have a check

func mergeMisspelling(checks []Checker) (res []Checker) {
	var missInDoc, missInCode []*propertyMissDiff
	for _, c := range checks {
		if p, ok := c.(*propertyMissDiff); ok {
			if p.MissType == MissInDoc {
				missInDoc = append(missInDoc, p)
			} else if p.MissType == MissInCode {
				missInCode = append(missInCode, p)
			}
		}
	}
	// check if missed name be like
	filterOut := map[*propertyMissDiff]struct{}{}
	for _, c := range missInCode {
		for _, d := range missInDoc {
			if dist := levenshtein(c.MDField().Name, d.mdField.Name); dist <= 3 {
				// if the edit distances less than 3, we think it's a misspelling
				filterOut[c] = struct{}{}
				filterOut[d] = struct{}{}
				res = append(res, newMisspelling(c, d))
			}
		}
	}
	for _, c := range checks {
		if _, ok := c.(*propertyMissDiff); !ok {
			res = append(res, c)
		}
	}
	return res
}

func levenshtein(str1, str2 string) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastKey := x - 1
		for y := 1; y <= s1len; y++ {
			oldKey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimumOf3(column[y]+1, column[y-1]+1, lastKey+incr)
			lastKey = oldKey
		}
	}
	return column[s1len]
}

func minimumOf3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else { // b < a
		if b < c {
			return b
		}
	}
	return c
}

package md

import (
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

// try to Unmarshal a markdown file to `model.ResourceDoc`

var posRegs map[model.PosType]*regexp.Regexp

func init() {
	InitRegExp()
}

func InitRegExp() {
	posRegs = map[model.PosType]*regexp.Regexp{
		model.PosExample: regexp.MustCompile("## Example"),
		model.PosArgs:    regexp.MustCompile(`## Arguments? Reference`),
		model.PosAttr:    regexp.MustCompile(`## Attributes? Reference`),
		model.PosTimeout: regexp.MustCompile(`## Timeouts?`),
		model.PosImport:  regexp.MustCompile(`## Imports?`),
	}
}

// `store_name` - (Required) The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.
var fieldReg = regexp.MustCompile("^[*-] *`(.*?)`" + ` +\- +(\(Required\)|\(Optional\))? ?(.*)`)

// var codeReg = regexp.MustCompile("`([a-zA-Z0-9-_ ,./~]+)`")
var codeReg = regexp.MustCompile("`([^`]+)`")

var blockPropReg = regexp.MustCompile("blocks?.*(as.*(below|above))")

var blockHeadReg = regexp.MustCompile("^(an?|An?|The)[^`]+(`[a-zA-Z0-9_]+`[, and]*)+.*blocks?.*$")

var DefaultsReg = regexp.MustCompile("[. ]+[Dd]efaults to `([^`]+)`[ .,]?")

func getDefaultValue(line string) string {
	if vals := DefaultsReg.FindStringSubmatch(line); len(vals) > 0 {
		return vals[1]
	}
	return ""
}

// ForceNewReg have to stop at dot or end of line to remove this part from the document when needed
var ForceNewReg = regexp.MustCompile(` ?Changing.*forces? a [^.]*(\.|$)`)

func isForceNew(line string) bool {
	return ForceNewReg.MatchString(line)
}

func extractFieldFromLine(line string) (field *model.Field) {
	field = &model.Field{
		Content: line,
	}
	// if defautl exists
	field.Default = getDefaultValue(line)
	field.ForceNew = isForceNew(line)

	res := fieldReg.FindStringSubmatch(line)
	if len(res) <= 1 || res[1] == "" {
		field.Name = util.FirstCodeValue(line) // try to use the first code as name
		field.FormatErr = "no field name found"
		return
	}
	field.Name = res[1]
	if field.Name == "" {
		log.Printf("field name is empty")
	}
	if len(res) > 2 {
		// may not exist
		switch {
		case strings.Contains(line, "(Required)"):
			field.Required = model.Required
		case strings.Contains(line, "(Optional)"):
			field.Required = model.Optional
		case strings.Contains(line, "Required"):
			field.Required = model.Required
		case strings.Contains(line, "Optional"):
			field.Required = model.Optional
		}
	}

	possibleValueSep := func(line string) int {
		line = strings.ToLower(line)
		for _, sep := range []string{"possible value", "be one of", "allowed value", "valid value",
			"supported value", "valid option", "accepted value"} {
			if sepIdx := strings.Index(line, sep); sepIdx >= 0 {
				return sepIdx
			}
		}
		return -1
	}

	var enums []string
	if len(res) > 3 {
		// extract enums from code part
		// from possible value to first '.'
		// skip if there are more than one sep exists
		// do not check the possible part
		if sepIdx := possibleValueSep(line); sepIdx > 0 {
			subStr := line[sepIdx:]
			field.EnumStart = sepIdx
			// end with dot may not work in values like `7.2` ....
			// should be . not in ` mark
			// Possible values are `a`, `b`, `a.b` and `def`.
			pointEnd := strings.Index(subStr, ".")
			if pointEnd < 0 {
				pointEnd = len(subStr)
			}
			enumIndex := codeReg.FindAllStringIndex(subStr, -1)
			for idx, val := range enumIndex {
				_ = idx
				start, end := val[0], val[1]
				if pointEnd > start && pointEnd < end {
					// point inside the code block
					if pointEnd = strings.Index(subStr[end:], "."); pointEnd < 0 {
						pointEnd = len(subStr)
					} else {
						pointEnd += end
					}
				}
				// search end to a dot
				if pointEnd < start {
					break
				}
				enums = append(enums, strings.Trim(subStr[start:end], "`'\""))
				field.EnumEnd = sepIdx + end
			}
			// breaks if there are more than 1 possible value
			if sepIdx = possibleValueSep(line[sepIdx+1:]); sepIdx >= 0 {
				field.Skip = true
			}
		}
		if len(enums) == 0 && strings.Index(res[3], "`") > 0 {
			guessValues := codeReg.FindAllString(res[3], -1)
			field.SetGuessEnums(guessValues)
		}
	}
	field.AddEnum(enums...)
	return field
}

func extractBlockNames(line string) (res []string) {
	if blockHeadReg.MatchString(line) {
		idx := strings.Index(line, "block")
		names := codeReg.FindAllString(line[:idx], -1)
		for idx, val := range names {
			names[idx] = strings.Trim(val, "`'")
		}
		return names
	}
	return
}

var blockPropGuessReg = regexp.MustCompile(`(defined|documented).*(below|above)`)
var blockPropGuessReg2 = regexp.MustCompile("(one or more) `.*` block")

func guessBlockProperty(line string) bool {
	if blockPropReg.MatchString(line) {
		return true
	}

	if blockPropGuessReg.MatchString(line) {
		return true
	}

	if blockPropGuessReg2.MatchString(line) {
		return true
	}
	if strings.Contains(line, "A block to") {
		return true
	}
	return false
}

func newFieldFromLine(line string) *model.Field {
	f := extractFieldFromLine(line)
	if guessBlockProperty(line) {
		// extract real block type
		f.BlockTypeName = f.Name
		// use the first code block value as block type name todo this may not right
		start := strings.Index(line, ")")
		end := strings.Index(line, "block")
		if start > 0 && end > 0 && start < end {
			if names := util.ExtractCodeValue(line[start:end]); len(names) > 0 {
				f.BlockTypeName = names[0]
			}
		}
		f.Typ = model.FieldTypeBlock
	}
	return f
}

func headPos(line string) (pos model.PosType) {
	if !strings.HasPrefix(line, "#") {
		return 0
	}
	for pos, reg := range posRegs {
		if reg.MatchString(line) {
			return pos
		}
	}
	// only head2
	if strings.HasPrefix(line, "##") && !strings.HasPrefix(line, "###") {
		return model.PosOther
	}
	return 0
}

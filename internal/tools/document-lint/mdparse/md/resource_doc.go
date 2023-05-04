package md

import (
	"bytes"
	"log"
	"os"
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

// ForceNewReg have to stop at dot or end of line to remove this part from document when needed
var ForceNewReg = regexp.MustCompile(` ?Changing.*forces? a [^.]*(\.|$)`)

func isForceNew(line string) bool {
	return ForceNewReg.MatchString(line)
}

func ExtractListItem(line string) (field *model.Field) {
	field = &model.Field{
		Content: line,
	}
	// if defautl exists
	field.Default = getDefaultValue(line)
	field.ForceNew = isForceNew(line)

	res := fieldReg.FindStringSubmatch(line)
	if len(res) <= 1 || res[1] == "" {
		field.Name = util.FirstCodeValue(line) // try to use the first code as name
		field.FormatErr = true
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
			// end with . not work in values like `7.2` ....
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
			// breaks if  there are more than 1 possible value
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

func ExtractBlockNames(line string) (res []string) {
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

func NewFieldFromLine(line string) *model.Field {
	f := ExtractListItem(line)
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

func UnmarshalResourceFromFile(filePath string) (res *model.ResourceDoc, err error) {
	content, _ := os.ReadFile(filePath)
	return UnmarshalResource(content)
}

// UnmarshalResource read line by line and unmarshal to a structure
func UnmarshalResource(content []byte) (res *model.ResourceDoc, err error) {
	if len(content) == 0 {
		return
	}
	res = model.NewResourceDoc()
	content = bytes.TrimSuffix(content, []byte{'\n'})
	lines := strings.Split(string(content), "\n")

	var curProp model.Properties
	_ = curProp
	var curPos = model.PosDefault
	var sameBlockNames []string
	var curXPath string

	var missSubBlocks []string
	removeMissSub := func(name string) {
		for idx, val := range missSubBlocks {
			if val == name {
				missSubBlocks = append(missSubBlocks[:idx], missSubBlocks[idx+1:]...)
				return
			}
		}
	}

	for lineNum := 1; lineNum <= len(lines); lineNum++ {
		txt := strings.TrimSpace(lines[lineNum-1])
		if strings.HasPrefix(txt, "* ") {
			if curPos == model.PosTimeout {
				res.SetTimeout(lineNum, txt)
				continue
			}
			if curProp == nil {
				// skip this property for now, only process arguments and attribute
				continue
			}
			// multiline doc for property
			startLine := lineNum
			for lineNum+1 <= len(lines) {
				nextTxt := lines[lineNum]
				var shouldBreak bool
				for _, sep := range []string{"*", "-", "#", "---"} {
					if strings.HasPrefix(nextTxt, sep) {
						shouldBreak = true
						break
					}
				}
				// contains specific prefix or is a block definition line
				if shouldBreak || blockHeadReg.MatchString(nextTxt) {
					break
				}
				txt += "\n" + nextTxt
				lineNum++
			}
			field := NewFieldFromLine(txt)
			field.Line = startLine
			field.Path = field.Name
			field.Pos = curPos
			if curXPath != "" {
				field.Path = curXPath + "." + field.Path
			}
			if field.Typ == model.FieldTypeBlock {
				if sub, ok := res.Blocks[field.BlockTypeName]; ok {
					field.Subs = sub
					removeMissSub(field.Name)
				}
			}
			curProp.AddField(field)
			res.PossibleValues[field.Path] = model.NewPossibleValue(field.PossibleValues(), field)
		} else if blocks := ExtractBlockNames(txt); len(blocks) > 0 {
			// only process args or attr for current version
			// there can be multiple property reference to this block
			if !curPos.IsArgOrAttr() {
				continue
			}
			sameBlockNames = blocks
			curProp = model.Properties{}
			curXPath = "."
		} else if pos := headPos(txt); txt == "---" || pos > 0 {
			for _, blockName := range sameBlockNames {
				if exists, ok := res.Blocks[blockName]; ok {
					exists.Merge(curProp)
				} else {
					res.Blocks[blockName] = curProp
				}
				top := res.CurProp(curPos)
				subBLocks := top.FindAllSubBlock(blockName)
				if len(subBLocks) == 0 {
					missSubBlocks = append(missSubBlocks, blockName)
				}
				for _, subBlock := range subBLocks {
					if len(subBlock.Subs) == 0 {
						subBlock.Subs = curProp
					}
				}
			}
			sameBlockNames = nil
			if pos > 0 {
				curPos = pos
			}
			curProp = res.CurProp(curPos)
			curXPath = ""
		} else if strings.HasPrefix(txt, "page_title:") {
			parts := strings.Split(txt, ":")
			if len(parts) == 3 {
				res.ResourceName = strings.Trim(parts[2], " \"")
			} else {
				if idx := strings.Index(txt, "auzrerm"); idx > 0 {
					res.ResourceName = strings.Trim(txt[idx:], " \"")
				}
			}
		}
	}
	// last try to find the block to link
	top := res.CurProp(model.PosArgs)
	for _, blockName := range missSubBlocks {
		subBlocks := top.FindAllSubBlock(blockName)
		if len(subBlocks) > 0 {
			removeMissSub(blockName)
		}
		for _, subBlock := range subBlocks {
			if len(subBlock.Subs) == 0 {
				subBlock.Subs = res.Blocks[blockName]
			}
		}

	}
	fixed := res.TuneSubBlocks()
	for _, name := range fixed {
		removeMissSub(name)
	}
	if len(missSubBlocks) > 0 {
		log.Printf("[doc] %s not block for names %v", res.ResourceName, missSubBlocks)
	}
	return res, nil
}

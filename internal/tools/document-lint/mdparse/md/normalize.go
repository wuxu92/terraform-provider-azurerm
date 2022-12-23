package md

import (
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
)

var orderFixMap = map[string]string{
	"(Optional) -": "- (Optional)",
	"(Required) -": "- (Required)",
	"`-":           "` -",
}

var (
	tryBlockHeadReg = regexp.MustCompile("^([*] )?(An?|The) +(`[a-zA-Z0-9_]+`[, and]*)+blocks?.*$")
	oldBlockHeadReg = regexp.MustCompile("^(`.*`) supports.*:$")
)

func tryBlockHeadDetect(line string) bool {
	if tryBlockHeadReg.MatchString(line) {
		return true
	}

	// sometimes a line with supports
	if strings.Contains(line, "supports") && strings.HasSuffix(line, ":") {
		return true
	}
	return false
}

func tryFixBlockHead(line string) string {
	if strings.HasPrefix(line, "*") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "*"))
	}

	if strings.HasPrefix(line, "`") {
		line = "The " + line
	}

	// add a block string after the second `
	if !strings.Contains(line, "block") {
		if idx := strings.LastIndexByte(line, '`'); idx > 0 {
			idx += 1
			line = line[:idx] + " block" + line[idx:]
		}
	}
	return line
}

func replaceNBSP(line string) string {
	isNormalSpace := func(r rune) bool {
		switch r {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			return true
		}
		return false
	}
	var res []rune
	for _, ch := range []rune(line) {
		if unicode.IsSpace(ch) && !isNormalSpace(ch) {
			res = append(res, rune(' '))
		} else {
			res = append(res, ch)
		}
	}
	return string(res)
}

func fixFileNormalize(file string) {
	contentBs, _ := os.ReadFile(file)
	content := string(contentBs)
	var hasModity bool
	_ = hasModity

	//content = bytes.TrimSuffix(content, []byte{'\n'})
	lines := strings.Split(string(content), "\n")
	var curScope model.PosType
	var newContent []string
	for idx, line := range lines {
		line = replaceNBSP(line)
		if pos := headPos(line); pos > 0 {
			curScope = pos
			//if pos > model.PosAttr {
			//	newContent = append(newContent, lines[idx:]...)
			//	break
			//}
		}
		if line == "--" {
			hasModity = true
			line = "---"
			lines[idx] = line
		}
		// some doc use `-` as list mark
		if strings.HasPrefix(line, "- `") {
			hasModity = true
			line = "*" + line[1:]
		}
		if !curScope.IsArgOrAttr() {
			newContent = append(newContent, line)
			continue
		}

		if tryBlockHeadDetect(line) {
			hasModity = true
			line = tryFixBlockHead(line)
		}

		if strings.Contains(line, "(Optional)") || strings.Contains(line, "(Required)") {
			if strings.HasPrefix(line, "`") {
				line = "* " + line
			}
		}

		isSep := func(line string) bool {
			return line == "---" || strings.HasPrefix(line, "#")
		}
		if blockHeadReg.MatchString(line) {
			if !isSep(lines[idx-1]) && !isSep(lines[idx-2]) {
				hasModity = true
				newContent = append(newContent, "---", "")
			}
		} else if strings.HasPrefix(line, "*") {
			// need a dash(-) after property name
			line2 := tryFixProp(line)
			hasModity = line != line2
			line = line2
			for k, v := range orderFixMap {
				if strings.Contains(line, k) && !strings.Contains(line, v) {
					hasModity = true
					line = strings.Replace(line, k, v, 1)
				}
			}
		}
		newContent = append(newContent, line)
	}
	//if hasModity {
	newBs := strings.Join(newContent, "\n")
	if newBs != content {
		f, _ := os.OpenFile(file, os.O_TRUNC|os.O_WRONLY, 0666)
		_, _ = f.WriteString(newBs)
		_ = f.Sync()
	}
	//}
}

func requireIndex(line string) int {
	if idx := strings.Index(line, "(Optional)"); idx > 0 {
		return idx
	}

	if idx := strings.Index(line, "(Required)"); idx > 0 {
		return idx
	}
	return -1
}

var tryBlockPropReg = regexp.MustCompile("[*] `.*` .*(One or more |A |A list of ).* blocks?.*")

var defaultValueReg = regexp.MustCompile(`[. ]+[d|D]efaults to (\w+)[,.]`)

func tryFixProp(line string) string {
	if strings.HasPrefix(line, "**") {
		return "~> " + line
	}
	if reqIdx := requireIndex(line); reqIdx > 0 {
		if !strings.HasSuffix(strings.TrimSpace(line[:reqIdx]), "-") {
			line = line[:reqIdx] + "- " + line[reqIdx:]
			reqIdx += 2
		}

		// a blank character before requiredness
		if divIdx := strings.Index(line[:reqIdx], "-"); divIdx > 0 {
			if ch := line[divIdx-1]; ch != ' ' {
				if unicode.IsSpace(rune(ch)) {
					line = line[:divIdx-1] + " " + line[divIdx:]
				} else {
					line = line[:divIdx] + " " + line[divIdx:]
				}
				divIdx += 1
			}
			if ch := line[divIdx+1]; ch != ' ' {
				if unicode.IsSpace(rune(ch)) {
					line = line[:divIdx+1] + " " + line[divIdx+2:]
				} else {
					line = line[:divIdx+1] + " " + line[divIdx+1:]
				}
			}
		}
		// one or more blocks as defined below
	} else {
		// add an optional/required to line
	}
	// need a dash after property name
	idx := strings.Index(line, "`")
	if idx += strings.Index(line[idx+1:], "`") + 1; idx > 0 {
		for idx2 := idx + 1; idx2 < len(line); idx2++ {
			if line[idx2] == ' ' {
				continue
			}
			if line[idx2] != '-' {
				line2 := line[:idx+2]
				if line[idx2-1] != ' ' {
					line2 += " "
				}
				line2 += "- "
				line2 += line[idx2:]
				line = line2
				break
			} else {
				// if line[idx2] == '-'  exists
				if line[idx2-1] != ' ' {
					line = line[:idx2] + " " + line[idx2:]
				}
				if line[idx2+1] != ' ' {
					line = line[:idx2+1] + " " + line[idx2+1:]
				}
				break
			}
		}
	}

	// add ` to default value
	if vals := defaultValueReg.FindStringSubmatchIndex(line); len(vals) > 0 {
		// add a backquote to value
		valStr := line[vals[2]:vals[3]]
		line = line[:vals[2]] + "`" + valStr + "`" + line[vals[3]:]
	}

	// need a blank
	if tryBlockPropReg.MatchString(line) {
		if !strings.Contains(line, "below") && !strings.Contains(line, "above") {
			if strings.HasSuffix(line, ".") {
				line = line[:len(line)-1]
			}
			line += " as defined below."
		}
	}

	return line
}

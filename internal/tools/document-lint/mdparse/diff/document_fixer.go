package diff

import (
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/md"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

// try to fix document with diff items

type Diffs []DiffItem

func (d Diffs) Len() int {
	return len(d)
}

func (d Diffs) Less(i, j int) bool {
	return d[i].Line < d[j].Line
}

func (d Diffs) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

var _ sort.Interface = &Diffs{}

// Fixer fix document with diff
type Fixer struct {
	MDFile       string
	SchemaFile   string
	ResourceType string

	Diff         Diffs            // diff of exist in both md and code
	DanglingInMD model.Properties // properties exists in markdown not in schema code
	MissingInMD  []string         // properties not exists in markdown

	FixedContent string
}

func NewFixer(d *ResourceDiff) *Fixer {
	f := &Fixer{
		MDFile:       d.MDFile,
		SchemaFile:   d.SchemaFile,
		ResourceType: d.tf.ResourceType,
		Diff:         d.Diff,
		FixedContent: "",
	}
	return f
}

func patchWantEnums(want []string) string {
	for idx, val := range want {
		want[idx] = "`" + val + "`"
	}
	if len(want) == 1 {
		return want[0]
	}

	res := want[0]
	if len(want) >= 3 {
		res = strings.Join(want[:len(want)-1], ", ")
	}
	res += " and " + want[len(want)-1]
	return res
}

func addRequiredness(line, req string) string {
	// add after the first -
	if idx := strings.Index(line, " - "); idx > 0 {
		return line[:idx+3] + req + " " + line[idx+3:]
	}
	// no dash add after second `
	idx := strings.Index(line, "`")
	idx += strings.Index(line[idx+1:], "`") + 1
	line = line[:idx+1] + " " + req + line[idx+1:]
	return line
}

func fixRequiredness(line string, required RequiredMiss) string {
	from, to := "(Required)", "(Optional)"

	switch required {
	case ShouldBeComputed:
		// remove from both. may cause two whitespace, have to remove one
		var idx, size int
		if idx = strings.Index(line, from); idx > 0 {
			size = len(from)
		} else if idx = strings.Index(line, to); idx > 0 {
			size = len(to)
		}
		if idx > 0 {
			if line[idx-1] == ' ' && line[idx+size] == ' ' {
				idx -= 1
				size += 1
			}
			line = line[:idx] + line[idx+size:]
		}
	case ShouldBeOptional, ShouldBeRequired:
		if required == ShouldBeRequired {
			from, to = to, from
		}
		if strings.Contains(line, from) {
			line = strings.Replace(line, from, to, 1)
		} else {
			line = addRequiredness(line, to)
		}
	}
	return line
}

// rt: resource type
func tryFixTimeouts(rt string, lines []string, diffs []TimeoutDiffItem) []string {
	var suf []string
	addSuf := func(line string) {
		suf = append(suf, line)
	}
	if diffs[0].Type == TimeoutMissed {
		// no such timeout block, add to the end of lines
		addSuf("## Timeouts")
		addSuf("")
		addSuf("The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:")
		addSuf("")
		diffs = diffs[1:]
	}
	// find timeout block
	var toLine, importIdx int
	for idx, line := range lines {
		if line == "## Timeouts" {
			toLine = idx + 4
			for toLine < len(lines) && lines[toLine] != "" {
				toLine++ // insert to an empty line
			}
		}
		if strings.HasPrefix(line, "## Import") {
			importIdx = idx
		}
	}
	rt = util.NormalizeResourceName(rt)

	for _, diff := range diffs {
		if diff.Line == 0 {
			// append line
			gen := diff.GenLine(rt)
			if len(suf) > 0 {
				addSuf(gen)
			} else {
				lines = append(lines[:toLine+1], lines[toLine:]...)
				lines[toLine] = gen
				//lines = append(append(lines[:toLine], gen), lines[toLine:]...)
			}
		} else {
			lines[diff.Line-1] = diff.FixLine(lines[diff.Line-1])
		}
	}
	if len(suf) > 0 {
		addSuf("")
		// insert before import
		if importIdx > 0 {
			end := make([]string, len(lines)-importIdx)
			copy(end, lines[importIdx:])
			lines = append(lines[:importIdx], suf...)
			lines = append(lines, end...)
		} else {
			lines = append(lines, suf...)
		}
	}
	return lines
}

func (f *Fixer) TryFix() {
	// read file as bytes
	if len(f.Diff) == 0 {
		return
	}
	//fd, err := os.Open(f.MDFile)
	content, err := os.ReadFile(f.MDFile)
	if err != nil {
		log.Printf("open %s: %v", f.MDFile, err)
	}
	// sort Diff by line
	sort.Sort(f.Diff)
	//sc := bufio.NewScanner(fd)
	lines := strings.Split(string(content), "\n")
	// find the resource name

	for idx, item := range f.Diff {
		// fix timeout!
		if len(item.TimeoutDiff) > 0 {
			lines = tryFixTimeouts(f.ResourceType, lines, item.TimeoutDiff)
			continue
		}

		// mdField is nil for no document exists or page title mismatch
		if item.MDFiled == nil || item.MDFiled.Skip {
			if item.MDFiled == nil && item.MissType == NotMiss {
				log.Printf("page title may mismatch, skip it...")
			}
			continue
		}
		// there maybe multiple reference to sub-block property, skip it
		if idx > 0 && item.Equals(f.Diff[idx-1]) {
			continue
		}

		lineIdx := item.Line - 1
		line := lines[lineIdx]

		// fix requiredness
		if item.RequiredMiss > 0 {
			lines[lineIdx] = fixRequiredness(line, item.RequiredMiss)
			continue
		}

		if item.DefaultDiff != "" || item.ShouldRemoveDefault {
			lines[lineIdx] = fixDefaultValue(line, item.DefaultDiff)
			continue
		}

		if item.ForceNewDiff > 0 {
			lines[lineIdx] = fixForceNewDiff(line, item.ForceNewDiff)
			continue
		}

		if len(item.Want) == 0 {
			continue
		}

		// replace from field.EnumStart to field.EnumEnd
		var bs strings.Builder
		if len(item.Got) == 0 || (item.MDFiled.EnumStart == 0 && len(item.Missed) > 0) {
			// skip this kind of field. may submit in a separate run
			// find default index
			idx := strings.Index(line, "Defaults to")
			if idx < 0 {
				idx = strings.Index(line, "Changing this forces")
			}
			if idx > 0 {
				bs.WriteString(line[:idx])
			} else {
				bs.WriteString(line)
				bs.WriteByte(' ')
			}
			if len(item.Want) == 1 {
				bs.WriteString("The only possible value is ")
			} else {
				bs.WriteString("Possible values are ")
			}
			bs.WriteString(patchWantEnums(item.Want))
			bs.WriteByte('.')
			if idx > 0 {
				bs.WriteByte(' ')
				bs.WriteString(line[idx:])
			}
			lines[lineIdx] = bs.String()
		} else if len(item.Missed) > 0 {
			// only replace missed values
			bs.WriteString(line[:item.MDFiled.EnumStart])
			if len(item.Want) == 1 {
				bs.WriteString("The only possible value is ")
			} else {
				bs.WriteString("Possible values are ")
			}
			bs.WriteString(patchWantEnums(item.Want))
			if item.MDFiled.EnumEnd < len(line) {
				bs.WriteString(line[item.MDFiled.EnumEnd:])
			} else {
				log.Printf("warning enum end %s:L%d len %dvs%d; %s", path.Base(f.MDFile), item.MDFiled.Line, item.MDFiled.EnumEnd, len(line), line)
			}
			lines[lineIdx] = bs.String()
		}
		if !strings.HasSuffix(strings.TrimSpace(lines[lineIdx]), ".") {
			lines[lineIdx] += "."
		}
	}
	f.FixedContent = strings.Join(lines, "\n")
	return
}

func fixForceNewDiff(line string, diff int) string {
	switch diff {
	case ShouldBeForceNew:
		line = strings.TrimRight(line, " ")
		if strings.HasSuffix(line, ",") {
			line = line[:len(line)-1] + "."
		} else if !strings.HasSuffix(line, ".") {
			line += "."
		}
		line += " Changing this forces a new resource to be created."
	case ShouldBeNotForceNew:
		line = md.ForceNewReg.ReplaceAllString(line, "")
	}
	return line
}

// if value if "", then we should remove the default value part from the document
func fixDefaultValue(line string, value string) string {
	if idxs := md.DefaultsReg.FindStringSubmatchIndex(line); len(idxs) > 2 {
		if value == "" {
			// remoev default part from line
			line = line[:idxs[0]+1] + line[idxs[1]:]
		} else {
			line = line[:idxs[2]] + value + line[idxs[3]:]
		}
	} else {
		line = strings.TrimSpace(line) + " Defaults to `" + value + "`."
	}
	return line
}

func (f *Fixer) WriteBack() {
	if len(f.Diff) == 0 {
		return
	}
	if f.FixedContent == "" {
		log.Printf("%s no content to write back, skip", f.MDFile)
		return
	}
	fd, err := os.OpenFile(f.MDFile, os.O_TRUNC|os.O_RDWR, 066)
	if err != nil {
		log.Printf("open %s: %v", f.MDFile, err)
	}
	defer func() {
		_ = fd.Sync()
		_ = fd.Close()
	}()
	_, err = fd.WriteString(f.FixedContent)
	if err != nil {
		log.Printf("write %s back: %v", f.MDFile, err)
	}
}

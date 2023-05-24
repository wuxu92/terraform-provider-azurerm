package md

import (
	"os"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

type ItemType int

const (
	Default ItemType = iota
	ItemMeteInfo
	ItemHeader1
	ItemHeader2
	ItemHeader3
	ItemExample
	ItemArgument
	ItemField
	ItemBlockHead // this is a xxx block supports
	ItemNote
	ItemSeparator
	ItemAttribute
	ItemPlainText
	//ItemBlankBlock
)

type MarkItem struct {
	FromLine int
	ToLine   int      // not set if this item has only one line
	lines    []string // by lines
	Type     ItemType
	Field    *model.Field
}

func (m *MarkItem) content() string {
	return strings.Join(m.lines, "\n")
}

func (m *MarkItem) addLine(num int, line string) {
	m.lines = append(m.lines, line)
	m.ToLine = num
}

func NewMarkItem(fromLine int, content string, typ ItemType) *MarkItem {
	m := &MarkItem{
		FromLine: fromLine,
		lines:    []string{content},
		Type:     typ,
	}
	return m
}

type Block struct {
	Names    []string // at least one name
	Of       string   // this is a block of xx Field, only some special blocks have it
	Name     string
	HeadLine int
	Fields   []*model.Field
}

func (b *Block) addField(f *model.Field) {
	b.Fields = append(b.Fields, f)
}

type Mark struct {
	Items        []*MarkItem
	content      *string
	FilePath     string
	ResourceType string // azurerm_xxx
	Blocks       []Block
	Fields       map[string]*model.Field
}

func (m *Mark) lastItem() *MarkItem {
	if len(m.Items) > 0 {
		return m.Items[len(m.Items)-1]
	}
	return nil
}

func (m *Mark) addItem(item *MarkItem) {
	m.Items = append(m.Items, item)
}

func (m *Mark) addItemWith(num int, line string, typ ItemType) {
	m.addItem(NewMarkItem(num, line, typ))
}

func (m *Mark) addLineOrItem(num int, line string, typ ItemType) {
	last := m.lastItem()
	if last.Type == typ {
		last.addLine(num, line)
	} else {
		m.addItem(NewMarkItem(num, line, typ))
	}
}

func (m *Mark) addLine(num int, line string) {
	m.lastItem().addLine(num, line)
}

func mustNewMarkFromFile(file string) *Mark {
	bs, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	m := newMarkFromString(string(bs))
	m.FilePath = file
	m.build()
	return m
}

func newMarkFromString(content string) *Mark {
	lines := strings.Split(content, "\n")
	result := &Mark{
		content: &content,
		Fields:  map[string]*model.Field{},
	}
	for idx, line := range lines {
		// if line starts with #, * or A Block supports, it is a special item
		switch {
		case strings.HasPrefix(line, "###"):
			result.addItem(NewMarkItem(idx, line, ItemHeader3))
			continue
		case strings.HasPrefix(line, "##"):
			result.addItem(NewMarkItem(idx, line, ItemHeader2))
			continue
		case strings.HasPrefix(line, "#"):
			result.addItem(NewMarkItem(idx, line, ItemHeader1))
			continue
		case strings.HasPrefix(line, "*"):
			result.addItem(NewMarkItem(idx, line, ItemField))
		case strings.HasPrefix(line, "---"):
			if idx == 0 {
				result.addItem(NewMarkItem(idx, line, ItemMeteInfo))
			} else {
				last := result.lastItem()
				if last.Type == ItemMeteInfo {
					last.addLine(idx, line)
				} else {
					result.addItem(NewMarkItem(idx, line, ItemSeparator))
				}
			}
		case strings.HasPrefix(line, "```"):
			result.addLineOrItem(idx, line, ItemExample)
		case strings.HasPrefix(line, "->"), strings.HasPrefix(line, "~>"):
			result.addItemWith(idx, line, ItemNote)
		case isBlockHead(line):
			result.addItemWith(idx, line, ItemBlockHead)
		default:
			// plain text
			last := result.lastItem()
			switch last.Type {
			case ItemField, ItemMeteInfo, ItemExample, ItemPlainText:
				last.addLine(idx, line)
			default:
				if strings.TrimSpace(line) == "" {
					// for empty lines, append to last item
					last.addLine(idx, line)
				} else {
					result.addItem(NewMarkItem(idx, line, ItemPlainText))
				}
			}
		}
	}
	return result
}

func isBlockHead(line string) bool {
	return blockHeadReg.MatchString(line)
}

func (m *Mark) addBlock(b Block) {
	m.Blocks = append(m.Blocks, b)
}

func (m *Mark) build() {
	var inBlock bool
	var block Block
	var pos model.PosType

	for _, item := range m.Items {
		content := item.content()
		switch item.Type {
		case ItemHeader1:
			trimmed := strings.TrimFunc(content, func(r rune) bool {
				if unicode.IsSpace(r) || r == '#' {
					return true
				}
				return false
			})
			if !strings.Contains(trimmed, " ") {
				m.ResourceType = trimmed
			}
		case ItemField:
			if pos > model.PosAttr {
				continue
			}

			f := NewFieldFromLine(content)
			f.Line = item.FromLine
			item.Field = f
			if inBlock {
				block.addField(f)
			} else {
				// field exists in both Argument and Attribute
				if arg, ok := m.Fields[f.Name]; ok {
					arg.SameNameAttr = f
				} else {
					m.Fields[f.Name] = f
				}
			}
		case ItemBlockHead:
			if pos > model.PosAttr {
				continue
			}
			if inBlock {
				m.addBlock(block)
			}
			names := ExtractBlockNames(item.lines[0])
			block = Block{
				Names:    names,
				Name:     names[0],
				HeadLine: item.FromLine,
			}
			// an of exists
			if ofIdx := strings.Index(content, "of"); ofIdx > 0 {
				block.Of = util.FirstCodeValue(content[ofIdx+2:])
			}
			inBlock = true
		case ItemSeparator:
			if inBlock {
				m.addBlock(block)
			}
			inBlock = false
		case ItemHeader2, ItemHeader3:
			if strings.Contains(content, "Argument") {
				pos = model.PosArgs
			} else if strings.Contains(content, "Attributes") {
				pos = model.PosAttr
			} else if strings.Contains(content, "Timeout") {
				pos = model.PosTimeout
			} else if strings.Contains(content, "Import") {
				pos = model.PosImport
			}

			if inBlock {
				m.addBlock(block)
			}
			inBlock = false
		}
	}
}

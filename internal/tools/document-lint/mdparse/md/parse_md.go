package md

import (
	"log"
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
	ItemTimeout
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
	asProp   model.Properties
}

func (b *Block) asProperties() model.Properties {
	if b.asProp == nil {
		res := model.Properties{}
		for _, f := range b.Fields {
			if _, ok := res[f.Name]; ok {
				log.Printf("duplicate field in block %s:%s", b.Name, f.Name)
			}
			res[f.Name] = f
		}
		b.asProp = res
	}
	return b.asProp
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
	result.buildField()
	result.buildStruct()
	return result
}

func isBlockHead(line string) bool {
	return blockHeadReg.MatchString(line)
}

func (m *Mark) addBlock(b Block) {
	m.Blocks = append(m.Blocks, b)
}

func (m *Mark) buildField() {
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
			if pos == model.PosTimeout {
				item.Type = ItemTimeout
				continue
			}
			if pos > model.PosAttr {
				continue
			}

			f := NewFieldFromLine(content)
			f.Line = item.FromLine
			f.Pos = pos
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

func (m *Mark) blockOfName(name string) *Block {
	var res []Block
	for _, b := range m.Blocks {
		for _, n2 := range b.Names {
			if n2 == name {
				res = append(res, b)
			}
		}
	}
	if len(res) == 0 {
		return nil
	}
	if len(res) > 1 {
		log.Printf("duplicate block exists for %s.%s", m.ResourceType, name)
	}
	return &res[0]
}

// buildStruct build struct of blocks
func (m *Mark) buildStruct() {
	fillField := func(f *model.Field) {
		if f.Typ == model.FieldTypeBlock {
			// find block
			if b := m.blockOfName(f.Name); b != nil {
				f.Subs = b.asProperties()
			} else {
				log.Printf("missing block for field %s.%s", m.ResourceType, f.Name)
			}
		}
	}

	for _, f := range m.Fields {
		fillField(f)
	}

	// build for block fields
	for _, b := range m.Blocks {
		for _, f := range b.Fields {
			fillField(f)
		}
	}
}

func (m *Mark) buildResourceDoc() *model.ResourceDoc {
	var doc = model.NewResourceDoc()
	for _, f := range m.Fields {
		if f.Pos == model.PosArgs {
			doc.Args.AddField(f)
		} else if f.Pos == model.PosAttr {
			doc.Attr.AddField(f)
		}
	}
	doc.ResourceName = m.ResourceType
	for _, item := range m.Items {
		if item.Type == ItemExample {
			doc.ExampleHCL = item.content()
		}
	}

	return doc
}

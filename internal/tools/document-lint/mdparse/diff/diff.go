package diff

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/md"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
)

// logic to load schema and markdown to print the diff

// diff for resource first find the resource schema, then find the Markdown document file for it(or specify by string)

type MissType int

const (
	NotMiss MissType = iota
	MissInDoc
	MissInCode
	MissInDocAttr
)

type RequiredMiss int

func (r RequiredMiss) String() string {
	return []string{"ok", "should be required", "should be optional"}[r]
}

const (
	RequiredOK       RequiredMiss = iota // no need to fix requiredness
	ShouldBeRequired                     // code is required, but doc be optional or not specify
	ShouldBeOptional                     // code is optional, but doc be required or not specify
	ShouldBeComputed
)

func (m MissType) String() string {
	return []string{"not miss", "doc", "code", "doc attribute"}[m]
}

type TimeoutType int

const (
	TimeoutMissed TimeoutType = iota // no tieout part in document
	TimeoutCreate
	TimeoutRead
	TimeoutUpdate
	TimeoutDelete
)

const (
	ForceNewDefault = iota
	ShouldBeForceNew
	ShouldBeNotForceNew
)

func (t TimeoutType) String() string {
	return []string{"", "create", "read", "update", "delete"}[t]
}

func (t TimeoutType) IngString() string {
	return []string{"", "creating", "retrieving", "updating", "deleting"}[t]
}

// GenLine generate a line for timeout if not exists
func (t *TimeoutDiffItem) GenLine(rt string) string {
	return fmt.Sprintf("* `%s` - (Defaults to %s) Used when %s the %s.",
		t.Type.String(),
		t.ValueString(),
		t.Type.IngString(),
		rt,
	)
}

type TimeoutDiffItem struct {
	Line int
	Type TimeoutType
	Want int64
}

func (t *TimeoutDiffItem) ValueString() string {
	val, suf := t.Want/60, "minute"
	if t.Want > 60*60 && (t.Want%(60*60)) == 0 {
		val, suf = t.Want/(60*60), "hour"
	}
	if val > 1 {
		suf += "s" // add a 's' suffix
	}
	return fmt.Sprintf("%d %s", val, suf)
}

func (t *TimeoutDiffItem) FixLine(line string) string {
	if t.Want <= 0 {
		return line
	}
	// find place to replace
	start, end := util.TimeoutValueIdx(line)
	if end <= start {
		return line
	}
	res := fmt.Sprintf("%s%s%s", line[:start], t.ValueString(), line[end:])
	return res
}

func NewTimeoutDiffItem(line int, typ TimeoutType, want int64) TimeoutDiffItem {
	return TimeoutDiffItem{
		Line: line,
		Type: typ,
		Want: want,
	}
}

type DiffItem struct {
	Key  string
	Line int
	Want []string
	Got  []string

	Missed []string // value not exists in doc
	Odd    []string // value not exists in code
	Msg    string

	MDFiled *model.Field // keep origin field pointer

	MissType MissType // default as NotMiss

	// required/optional miss
	RequiredMiss RequiredMiss

	// save timeout diffs to a list
	TimeoutDiff []TimeoutDiffItem

	// Default value mismatch
	DefaultDiff         string // the right default value
	ShouldRemoveDefault bool

	ForceNewDiff int
}

func NewFoceNewDiff(key string, f *model.Field, forceNew int) DiffItem {
	ins := DiffItem{
		Key:          key,
		Line:         f.Line,
		MDFiled:      f,
		ForceNewDiff: forceNew,
	}
	return ins
}

// NewDefaultDiff if defaultValue is "", then should remove default from doc
// for really "" default value, it should be `""` string
func NewDefaultDiff(key string, f *model.Field, defaultValue string) DiffItem {
	ins := DiffItem{
		Key:         key,
		Line:        f.Line,
		MDFiled:     f,
		DefaultDiff: defaultValue,
	}
	if defaultValue == "" {
		ins.ShouldRemoveDefault = true
	}
	return ins
}

func (d DiffItem) Equals(dest DiffItem) bool {
	if d.Line != dest.Line {
		return false
	}

	if d.MissType != dest.MissType || d.RequiredMiss != dest.RequiredMiss || d.ForceNewDiff != dest.ForceNewDiff ||
		d.DefaultDiff != d.DefaultDiff || len(d.TimeoutDiff) != len(dest.TimeoutDiff) {
		return false
	}

	if len(d.Want) != len(dest.Want) || len(d.Got) != len(dest.Got) {
		return false
	}
	return true
}

func NewMissDiffItem(key string, typ MissType) DiffItem {
	return DiffItem{
		Key:      key,
		MissType: typ,
	}
}

func NewRequiredDiffItem(key string, mdField *model.Field, typ RequiredMiss) DiffItem {
	return DiffItem{
		Key:          key,
		Line:         mdField.Line,
		MDFiled:      mdField,
		RequiredMiss: typ,
	}
}

func NewDiffItem(key string, want []string, mdField *model.Field, missed, odd []string) DiffItem {
	return DiffItem{
		Key:     key,
		Line:    mdField.Line,
		Want:    want,
		Got:     mdField.PossibleValues(),
		MDFiled: mdField,
		Missed:  missed,
		Odd:     odd,
	}
}

type ResourceDiff struct {
	tf *schema.Resource
	md *model.ResourceDoc

	SchemaFile string
	MDFile     string

	Diff []DiffItem // diff of exist in both md and code
}

func (d *ResourceDiff) Diffs() []DiffItem {
	return d.Diff
}

func (d *ResourceDiff) ToString() string {
	var bs strings.Builder

	bs.WriteString(
		fmt.Sprintf("%s: diff size: %d, document file: %s\n",
			d.tf.ResourceType,
			len(d.Diff),
			path.Base(d.MDFile),
		),
	)
	bs.WriteString(DiffResultToString(d))
	return bs.String()
}

// NewResourceDiff tf is required,
// mdPath is optional, it can be detected by resource name
func NewResourceDiff(tf *schema.Resource) *ResourceDiff {
	r := &ResourceDiff{
		tf: tf,
	}
	// try to detect Markdown path from resource
	// can set it if not a regular MD path
	r.MDFile = md.MDPathFor(tf.ResourceType)
	return r
}

func (r *ResourceDiff) DiffAll() {
	if r.md == nil {
		r.md, _ = md.UnmarshalResourceFromFile(r.MDFile)
	}
	r.Diff = diffEunms(r.tf, r.md)

	missDiff := crossCheckProperty(r.tf, r.md)
	r.Diff = append(r.Diff, missDiff...)

	timeouts := diffTimeout(r.tf, r.md)
	r.Diff = append(r.Diff, timeouts...)
}

func diffTimeout(r *schema.Resource, md *model.ResourceDoc) (res []DiffItem) {
	to := r.Schema.Timeouts
	if to == nil {
		return
	}
	var items []TimeoutDiffItem
	if md.Timeouts == nil {
		items = append(items, NewTimeoutDiffItem(0, TimeoutMissed, 0))
		md.Timeouts = &model.Timeouts{} // use an empty timeouts object
	}
	if ptr := to.Read; ptr != nil {
		val := int64((*ptr) / time.Second)
		if mdVal := md.Timeouts.Read; val != mdVal.Value {
			items = append(items, NewTimeoutDiffItem(mdVal.Line, TimeoutRead, val))
		}
	}
	if ptr := to.Create; ptr != nil {
		val := int64((*ptr) / time.Second)
		if mdVal := md.Timeouts.Create; val != mdVal.Value {
			items = append(items, NewTimeoutDiffItem(mdVal.Line, TimeoutCreate, val))
		}
	}
	if ptr := to.Update; ptr != nil {
		val := int64((*ptr) / time.Second)
		if mdVal := md.Timeouts.Update; val != mdVal.Value {
			items = append(items, NewTimeoutDiffItem(mdVal.Line, TimeoutUpdate, val))
		}
	}
	if ptr := to.Delete; ptr != nil {
		val := int64((*ptr) / time.Second)
		if mdVal := md.Timeouts.Delete; val != mdVal.Value {
			items = append(items, NewTimeoutDiffItem(mdVal.Line, TimeoutDelete, val))
		}
	}
	if len(items) > 0 {
		res = append(res, DiffItem{TimeoutDiff: items})
	}
	return
}

func diffEunms(r *schema.Resource, md *model.ResourceDoc) (res []DiffItem) {
	schemModel := r.Schema.Schema
	_ = schemModel
	if md == nil {
		res = append(res, DiffItem{
			Key:  "no match document exists",
			Line: 0,
			Want: nil,
			Got:  nil,
			Msg:  fmt.Sprintf("document file name of page title did not match the resource type: %s", r.ResourceType),
		})
		return
	}
	// loop over document model
	for name, field := range md.Args {
		partRes := diffField(r, field, []string{name})
		res = append(res, partRes...)
	}
	return
}

// xPath property name for parent nodes
func diffField(r *schema.Resource, mdField *model.Field, xPath []string) (res []DiffItem) {

	fullPath := strings.Join(xPath, ".")
	if isSkipProp(r.ResourceType, fullPath) {
		return
	}

	// if end property
	if mdField.Subs == nil {
		want := r.PossibleValues[fullPath]
		docVal := mdField.PossibleValues()
		if missed, odd := SliceDiff(want, docVal, true); len(missed)+len(odd) > 0 {
			if !mayExistsInDoc(mdField.Content, want) {
				res = append(res, NewDiffItem(fullPath, want, mdField, missed, odd))
			}
		}
		return
	}
	for _, sub := range mdField.Subs {
		subRes := diffField(r, sub, append(xPath, sub.Name))
		res = append(res, subRes...)
	}
	return
}

func SliceDiff(want, got []string, caseInSensitive bool) (missed, odd []string) {
	// if `want` is nil then it may only write in doc, skip this
	if len(want) == 0 {
		return
	}
	// cross-check
	wantCpy, gotCpy := want, got
	if caseInSensitive {
		wantCpy = make([]string, len(want))
		gotCpy = make([]string, len(got))
		for idx := range want {
			wantCpy[idx] = strings.ToLower(want[idx])
		}
		for idx := range got {
			gotCpy[idx] = strings.ToLower(got[idx])
		}
	}
	wantMap := util.Slice2Map(wantCpy)
	gotMap := util.Slice2Map(gotCpy)

	for idx, k := range wantCpy {
		if _, ok := gotMap[k]; !ok {
			missed = append(missed, want[idx])
		}
	}

	for idx, k := range gotCpy {
		if _, ok := wantMap[k]; !ok {
			odd = append(odd, got[idx])
		}
	}

	return
}

// return true values exists in doc but may not with code quote
func mayExistsInDoc(docLine string, want []string) bool {
	for _, val := range want {
		if !strings.Contains(docLine, val) {
			return false
		}
	}
	return true
}

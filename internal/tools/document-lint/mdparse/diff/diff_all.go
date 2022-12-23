package diff

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/schema"
)

// function to call diff all resource provided by expose_schema.go

type IDiffResult interface {
	DiffAll()
	Diffs() []DiffItem
	ToString() string
}

func DiffResultToString(res IDiffResult) string {
	var bs strings.Builder
	for _, item := range res.Diffs() {
		if item.MDFiled != nil && item.RequiredMiss == RequiredOK {
			bs.WriteString(fmt.Sprintf("\t%s@%d skip: %v:\n\t\tmiss in doc: %v\n\t\tmay miss in code: %v\n", item.Key, item.Line, item.MDFiled.Skip, item.Missed, item.Odd))
		}
		if item.MissType != NotMiss {
			bs.WriteString(fmt.Sprintf("\t miss %s in %s\n", item.Key, item.MissType))
		}
		if item.RequiredMiss != RequiredOK {
			bs.WriteString(fmt.Sprintf("\t required miss: %s", item.RequiredMiss))
		}
	}
	return bs.String()
}

type DiffResult struct {
	start  time.Time
	end    time.Time
	result []*ResourceDiff
	byCate map[string][]*ResourceDiff
	mux    sync.Mutex
}

func NewDiffResult() *DiffResult {
	return &DiffResult{
		byCate: map[string][]*ResourceDiff{},
		start:  time.Now(),
	}
}

func (d *DiffResult) FixDocuments() {
	for _, r := range d.result {
		//if rd, ok := r.(*ResourceDiff); ok {
		fix := NewFixer(r)
		fix.TryFix()
		fix.WriteBack()
		//}
	}
}

func (d *DiffResult) GetResult() []*ResourceDiff {
	return d.result
}

func (d *DiffResult) HasDiff() bool {
	result := d.GetResult()
	return len(result) > 0
}

func (d *DiffResult) ToString() string {
	var bs strings.Builder
	var count int
	var missCount int
	var crossCount, resourceCount int
	for _, diff := range d.result {
		if len(diff.Diffs()) > 0 {
			resourceCount++
			count += len(diff.Diffs())
			bs.WriteString(diff.ToString())
			bs.WriteString("\n\n")

			for _, df := range diff.Diffs() {
				if len(df.Missed) > 0 {
					missCount += 1
				}
				if df.MissType > 0 {
					crossCount++
				}
			}
		}
	}
	bs.WriteString(
		fmt.Sprintf("total diff find: %d, missed in doc count: %d; crosscheck miss: %d. resource count: %d, cost: %s\n",
			count, missCount, crossCount, resourceCount, d.end.Sub(d.start)))
	return bs.String()
}

func (d *DiffResult) CostTime() time.Duration {
	return d.end.Sub(d.start)
}

const (
	_ = iota
	DiffV1
)

func DiffAll(regs Registers) *DiffResult {
	return doDiffAll(regs) // v2
}

func doDiffAll(regs Registers) *DiffResult {
	var dr = NewDiffResult()
	//runtime.GOMAXPROCS(1)
	// loop over
	var wg sync.WaitGroup
	// for debug to run only specific resource
	skipResource := func(name string) bool {
		if env := os.Getenv("ONLY_RESOURCE"); len(env) > 0 && env != "azurerm_" {
			return name != env
		}
		target := ""
		if target == "" {
			return false
		}
		return name != target
	}
	// can not split to package in different goroutine which may cause data-race and mix shared pointer up
	// register may repeat in typed and untyped, so use a map to remove the repeat entry
	var regMap = map[interface{}]bool{}
	for _, reg := range regs.Registers {
		if _, ok := regMap[reg]; ok {
			continue
		}
		regMap[reg] = true

		reg := reg
		wg.Add(1)
		go func() {
			defer wg.Done()
			// one register exists in both typed register and untyped register
			var rds []*ResourceDiff
			var catName string
			if typed, ok := reg.(sdk.TypedServiceRegistration); ok {
				for _, res := range typed.Resources() {
					if skipResource(res.ResourceType()) {
						continue
					}
					sch := schema.NewResourceByTyped(res)
					rd := NewResourceDiff(sch)
					rd.DiffAll()

					rds = append(rds, rd)
					catName = typed.Name()
				}
			}
			if untyped, ok := reg.(sdk.UntypedServiceRegistration); ok {
				for name, res := range untyped.SupportedResources() {
					if skipResource(name) {
						continue
					}

					sch := schema.NewResourceByUntyped(res, name)
					rd := NewResourceDiff(sch)
					rd.DiffAll()
					rds = append(rds, rd)
					catName = untyped.Name()
				}
			}
			if len(rds) > 0 {
				dr.mux.Lock()
				dr.result = append(dr.result, rds...)
				dr.byCate[catName] = rds
				dr.mux.Unlock()
			}
		}()
	}
	wg.Wait()
	dr.end = time.Now()
	return dr
}

package diff

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/util"
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
		if item.MDFiled != nil && len(item.Missed) > 0 {
			bs.WriteString(fmt.Sprintf("\t%s@%d skip: %v:\n\t\tmiss in doc: %v\n\t\tmay miss in code: %v\n", item.Key, item.Line, item.MDFiled.Skip, item.Missed, item.Odd))
		}
		if item.MissType != NotMiss {
			bs.WriteString(fmt.Sprintf("\t miss %s in %s\n", item.Key, item.MissType))
		}
		if item.RequiredMiss != RequiredOK {
			bs.WriteString(fmt.Sprintf("\t %s required miss: %s", item.Key, item.RequiredMiss))
		}
		if item.DefaultDiff != "" || item.ShouldRemoveDefault {
			bs.WriteString(fmt.Sprintf("\t %s default vlaue: %s", item.Key, item.DefaultDiff))
		}
		if item.ForceNewDiff > 0 {
			bs.WriteString(fmt.Sprintf("\t %s ForceNew: %d", item.Key, item.ForceNewDiff))
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
	for _, r := range result {
		for _, f := range r.Diff {
			if f.MDFiled == nil || !f.MDFiled.Skip {
				return true
			}
		}
	}
	return false
}

func (d *DiffResult) ToString() string {
	var bs strings.Builder
	var count int
	var possiblevalueMiss int
	var crossCount, resourceCount int
	var reqCount, defaultCount, timeoutCount, forceNewCount int
	var skipCount int
	for _, diff := range d.result {
		if len(diff.Diffs()) > 0 {
			resourceCount++
			count += len(diff.Diffs())
			bs.WriteString(diff.ToString())
			bs.WriteString("\n\n")

			for _, df := range diff.Diffs() {
				if df.MDFiled != nil && df.MDFiled.Skip {
					skipCount++
					continue
				}
				if df.MissType > 0 {
					crossCount++
					continue
				}
				if len(df.Missed) > 0 {
					possiblevalueMiss += 1
				}
				if df.RequiredMiss > 0 {
					reqCount++
				}
				if df.DefaultDiff != "" || df.ShouldRemoveDefault {
					defaultCount++
				}
				if len(df.TimeoutDiff) > 0 {
					timeoutCount += len(df.TimeoutDiff)
				}
				if df.ForceNewDiff > 0 {
					forceNewCount++
				}
			}
		}
	}
	bs.WriteString(
		fmt.Sprintf(
			`------
total issues find:    %d
possible value count: %d
requiredness count:   %d
default value count:  %d
timeout value count:  %d
skip property count:  %d
force new count: %d
crosscheck miss: %d
resource count:  %d
time costs: %s
------`,
			count, possiblevalueMiss, reqCount, defaultCount, timeoutCount, forceNewCount,
			skipCount, crossCount, resourceCount, d.end.Sub(d.start)))
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
	var wg sync.WaitGroup
	// for debug to run only specific resource
	skipByResource := func(name string) bool {
		if env := os.Getenv("ONLY_RESOURCE"); len(env) > 0 && env != "azurerm_" {
			return name != env
		}

		return isSkipResource(name)
	}

	skipByRP := func(name string) bool {
		own := util.GetRPOwner(name)
		return own == "xiaxin.yi"
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
				catName = typed.Name()
				if skipByRP(catName) {
					return
				}
				for _, res := range typed.Resources() {
					if skipByResource(res.ResourceType()) {
						continue
					}
					sch := schema.NewResourceByTyped(res)
					rd := NewResourceDiff(sch)
					rd.DiffAll()

					if len(rd.Diffs()) > 0 {
						rds = append(rds, rd)
					}
				}
			}
			if untyped, ok := reg.(sdk.UntypedServiceRegistration); ok {
				catName = untyped.Name()
				if skipByRP(catName) {
					return
				}
				for name, res := range untyped.SupportedResources() {
					if skipByResource(name) {
						continue
					}

					sch := schema.NewResourceByUntyped(res, name)
					rd := NewResourceDiff(sch)
					rd.DiffAll()
					if len(rd.Diffs()) > 0 {
						rds = append(rds, rd)
					}
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

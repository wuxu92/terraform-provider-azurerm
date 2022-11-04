package validation

import (
	"reflect"
	"sync"
)

var stringInSliceMap = map[uintptr][]string{}

var mux sync.RWMutex

func addSliceToMap(f interface{}, values []string) {
	// get call stack and get function ptr
	vd := reflect.ValueOf(f)
	ptr := vd.Pointer()
	mux.Lock()
	stringInSliceMap[ptr] = values
	mux.Unlock()
}

func GetFunctionValues(ptr uintptr) []string {
	mux.RLock()
	res := stringInSliceMap[ptr]
	mux.RUnlock()
	return res
}

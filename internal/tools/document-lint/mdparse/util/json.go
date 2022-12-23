package util

import (
	"encoding/json"
	"log"
)

func Stringify(v interface{}, indent bool) string {
	var bs []byte
	var err error
	if indent {
		bs, err = json.MarshalIndent(v, "", "  ")
	} else {
		bs, err = json.Marshal(v)
	}
	if err != nil {
		log.Printf("json marshal err: %v", err)
	}
	return string(bs)
}

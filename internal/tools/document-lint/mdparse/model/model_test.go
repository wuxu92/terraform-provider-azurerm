package model

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestJSON(t *testing.T) {
	t.Skip("skip lang test")
	content := []byte(`abcdef`)
	var ptr *string

	err := json.Unmarshal(content, ptr)
	if err != nil {
		t.Errorf("json.Unmarshal: %v", err)
	}

	dec := json.NewDecoder(bytes.NewBuffer(content))
	if err = dec.Decode(ptr); err != nil {
		t.Errorf("decode: %v", err)
	}
}

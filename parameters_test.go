package govaluate

import (
	"testing"
)

func TestMapParameters_Get(t *testing.T) {
	data := map[string]interface{}{
		"obj": map[string]interface{}{
			"item": true,
		},
	}
	parameters := MapParameters(data)
	val, err := parameters.Get("obj.item")
	if err != nil {
		t.Logf("Unexpected error while retrieving parameters: %v", err)
		t.Fail()
	}
	if val != true {
		t.Log("Unexpected value from parameter")
		t.Fail()
	}
}

package json_test

import (
	"os"
	"testing"

	"github.com/goccy/go-json"
)

var jsonFile, err = os.Open("json/small.json")

func TestValidWithComplexData(t *testing.T) {
	if err := json.Unmarshal(got, &s2); err != nil {
		t.Fatalf("Decode: %v", err)
	}
}

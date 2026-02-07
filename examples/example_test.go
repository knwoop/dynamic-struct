package examples

import (
	"fmt"
	"testing"
)

func TestExample(t *testing.T) {
	out, err := getJSONExample()
	if err != nil {
		t.Fatalf("getJSONExample() error: %v", err)
	}
	fmt.Printf("--- output:\n%s\n", out)
}

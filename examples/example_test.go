package _examples

import (
	"fmt"
	"testing"
)

func TestExample(t *testing.T) {
	yml := getYamlExample()
	fmt.Printf("--- t:\n%v\n\n", yml)
}

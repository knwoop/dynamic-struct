package examples

import (
	"encoding/json"
	"fmt"
	"reflect"

	dynamicstruct "github.com/knwoop/dynamic-struct"
)

func getJSONExample() (string, error) {
	ds, err := dynamicstruct.NewBuilder().
		AddField("Name", "", `json:"name"`).
		AddField("Age", 0, `json:"age"`).
		AddField("Email", "", `json:"email"`).
		AddField("Scores", []int{}, `json:"scores"`).
		AddField("Active", false, `json:"active"`).
		Build()
	if err != nil {
		return "", fmt.Errorf("build: %w", err)
	}

	instance := ds.New()

	data := []byte(`{
	"name": "Alice",
	"age": 30,
	"email": "alice@example.com",
	"scores": [95, 87, 92],
	"active": true
}`)

	if err := json.Unmarshal(data, instance); err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	fmt.Printf("--- type:\n%v\n\n", reflect.TypeOf(instance))
	fmt.Printf("--- instance:\n%v\n\n", instance)

	out, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	return string(out), nil
}

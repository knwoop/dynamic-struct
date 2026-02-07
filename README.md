# dynamic-struct

A Go library for building structs dynamically at runtime. Create struct types on the fly, set fields via reflection, and seamlessly integrate with `encoding/json` and other encoding packages.

## Installation

```
go get github.com/knwoop/dynamic-struct
```

Requires Go 1.24 or later. Zero external dependencies.

## Usage

### Basic example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	dynamicstruct "github.com/knwoop/dynamic-struct"
)

func main() {
	// Define the struct
	ds, err := dynamicstruct.NewBuilder().
		AddField("Name", "", `json:"name"`).
		AddField("Age", 0, `json:"age"`).
		AddField("Active", false, `json:"active"`).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	// Create an instance
	instance := ds.New()

	// Unmarshal JSON into the dynamic struct
	data := []byte(`{"name":"Alice","age":30,"active":true}`)
	if err := json.Unmarshal(data, instance); err != nil {
		log.Fatal(err)
	}

	// Read fields via reflection
	v := reflect.ValueOf(instance).Elem()
	fmt.Println(v.FieldByName("Name").String())  // Alice
	fmt.Println(v.FieldByName("Age").Int())       // 30
	fmt.Println(v.FieldByName("Active").Bool())   // true

	// Marshal back to JSON
	out, _ := json.Marshal(instance)
	fmt.Println(string(out)) // {"active":true,"age":30,"name":"Alice"}
}
```

### Field types

The type of each field is determined by the zero-value you pass to `AddField`. Pass any Go value to define the field's type:

```go
builder.AddField("Text", "", "")           // string
builder.AddField("Number", 0, "")          // int
builder.AddField("Price", 0.0, "")         // float64
builder.AddField("Flag", false, "")        // bool
builder.AddField("Items", []string{}, "")  // []string
```

## API Reference

### `NewBuilder() *Builder`

Creates a new `Builder` instance.

### `(*Builder) AddField(name string, typ any, tag string) *Builder`

Adds a field to the struct definition. Returns the builder for method chaining.

- `name` must be non-empty and exported (start with an uppercase letter).
- `typ` determines the field's type by its concrete Go type. Must not be `nil`.
- `tag` is the struct tag (e.g., `` `json:"name"` ``).

Validation errors are accumulated and returned when `Build()` is called.

### `(*Builder) RemoveField(name string) *Builder`

Removes a previously added field by name.

### `(*Builder) HasField(name string) bool`

Returns `true` if a field with the given name exists.

### `(*Builder) GetField(name string) *Field`

Returns the `Field` with the given name, or `nil` if not found.

### `(*Builder) Build() (*DynamicStruct, error)`

Builds the struct type definition. Fields are sorted by name for deterministic ordering. Returns an error if any validation errors were accumulated during `AddField` calls.

### `(*DynamicStruct) New() any`

Creates a new instance of the dynamic struct (returns a pointer).

## License

[MIT](LICENSE)

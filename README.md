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
	// Define the struct using type-specific methods
	ds, err := dynamicstruct.NewBuilder().
		AddStringField("Name", `json:"name"`).
		AddIntField("Age", `json:"age"`).
		AddBoolField("Active", `json:"active"`).
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

### Type-specific methods

For common Go types, dedicated methods are provided so you don't need to pass zero values:

```go
builder.AddStringField("Name", `json:"name"`)     // string
builder.AddIntField("Count", `json:"count"`)       // int
builder.AddInt64Field("ID", `json:"id"`)           // int64
builder.AddFloat64Field("Price", `json:"price"`)   // float64
builder.AddBoolField("Active", `json:"active"`)    // bool
builder.AddTimeField("CreatedAt", `json:"created_at"`)  // time.Time
builder.AddSliceField("Tags", "", `json:"tags"`)   // []string (element type inferred from 2nd arg)
builder.AddSliceField("Scores", 0, `json:"scores"`) // []int
```

Full list: `AddStringField`, `AddBoolField`, `AddIntField`, `AddInt8Field`, `AddInt16Field`, `AddInt32Field`, `AddInt64Field`, `AddUintField`, `AddUint8Field`, `AddUint16Field`, `AddUint32Field`, `AddUint64Field`, `AddFloat32Field`, `AddFloat64Field`, `AddTimeField`, `AddSliceField`.

### Generic helper

For any type not covered by the built-in methods, use the generic `AddTypedField` function:

```go
type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

b := dynamicstruct.NewBuilder()
dynamicstruct.AddTypedField[string](b, "Name", `json:"name"`)
dynamicstruct.AddTypedField[[]int](b, "Scores", `json:"scores"`)
dynamicstruct.AddTypedField[time.Time](b, "CreatedAt", `json:"created_at"`)
dynamicstruct.AddTypedField[Address](b, "Addr", `json:"addr"`)
ds, err := b.Build()
```

Since Go does not support type parameters on methods, `AddTypedField` is a package-level function. It accepts any type via its type parameter, making it possible to add fields of custom structs, maps, or any other Go type.

### Low-level AddField

You can also use the general-purpose `AddField` method directly, passing a zero value of the desired type:

```go
builder.AddField("Name", "", `json:"name"`)           // string
builder.AddField("Count", 0, `json:"count"`)          // int
builder.AddField("Price", 0.0, `json:"price"`)        // float64
builder.AddField("Items", []string{}, `json:"items"`) // []string
```

## API Reference

### `NewBuilder() *Builder`

Creates a new `Builder` instance.

### `(*Builder) AddField(name string, typ any, tag string) *Builder`

Adds a field to the struct definition. The field's type is inferred from the concrete type of `typ`. Returns the builder for method chaining.

- `name` must be non-empty and exported (start with an uppercase letter).
- `typ` determines the field's type by its concrete Go type. Must not be `nil`.
- `tag` is the struct tag (e.g., `` `json:"name"` ``).

Validation errors are accumulated and returned when `Build()` is called.

### Type-specific methods

All return `*Builder` for chaining. Signature: `(name string, tag string) *Builder`

| Method | Field type |
|--------|-----------|
| `AddStringField` | `string` |
| `AddBoolField` | `bool` |
| `AddIntField` | `int` |
| `AddInt8Field` ~ `AddInt64Field` | `int8` ~ `int64` |
| `AddUintField` ~ `AddUint64Field` | `uint` ~ `uint64` |
| `AddFloat32Field` | `float32` |
| `AddFloat64Field` | `float64` |
| `AddTimeField` | `time.Time` |

## License

[MIT](LICENSE)

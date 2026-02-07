package dynamicstruct

import (
	"errors"
	"fmt"
	"go/token"
	"reflect"
	"sort"
	"time"
)

// Builder builds a struct dynamically.
type Builder struct {
	fields map[string]*Field
	errs   []error
}

// NewBuilder returns a new clean instance of Builder.
// builder := dynamicstruct.NewBuilder()
func NewBuilder() *Builder {
	return &Builder{
		fields: map[string]*Field{},
	}
}

// AddField adds a field to the struct being built.
// It validates the field name and type, accumulating errors
// that are reported when Build is called.
func (b *Builder) AddField(name string, typ any, tag string) *Builder {
	if name == "" {
		b.errs = append(b.errs, errors.New("field name must not be empty"))
		return b
	}
	if !token.IsExported(name) {
		b.errs = append(b.errs, fmt.Errorf("field name %q must be exported (start with uppercase)", name))
		return b
	}
	if typ == nil {
		b.errs = append(b.errs, fmt.Errorf("field %q: type must not be nil", name))
		return b
	}
	b.fields[name] = &Field{Name: name, Type: typ, Tag: tag}
	return b
}

// AddStringField adds a string field.
func (b *Builder) AddStringField(name string, tag string) *Builder {
	return b.AddField(name, "", tag)
}

// AddBoolField adds a bool field.
func (b *Builder) AddBoolField(name string, tag string) *Builder {
	return b.AddField(name, false, tag)
}

// AddIntField adds an int field.
func (b *Builder) AddIntField(name string, tag string) *Builder {
	return b.AddField(name, 0, tag)
}

// AddInt8Field adds an int8 field.
func (b *Builder) AddInt8Field(name string, tag string) *Builder {
	return b.AddField(name, int8(0), tag)
}

// AddInt16Field adds an int16 field.
func (b *Builder) AddInt16Field(name string, tag string) *Builder {
	return b.AddField(name, int16(0), tag)
}

// AddInt32Field adds an int32 field.
func (b *Builder) AddInt32Field(name string, tag string) *Builder {
	return b.AddField(name, int32(0), tag)
}

// AddInt64Field adds an int64 field.
func (b *Builder) AddInt64Field(name string, tag string) *Builder {
	return b.AddField(name, int64(0), tag)
}

// AddUintField adds a uint field.
func (b *Builder) AddUintField(name string, tag string) *Builder {
	return b.AddField(name, uint(0), tag)
}

// AddUint8Field adds a uint8 field.
func (b *Builder) AddUint8Field(name string, tag string) *Builder {
	return b.AddField(name, uint8(0), tag)
}

// AddUint16Field adds a uint16 field.
func (b *Builder) AddUint16Field(name string, tag string) *Builder {
	return b.AddField(name, uint16(0), tag)
}

// AddUint32Field adds a uint32 field.
func (b *Builder) AddUint32Field(name string, tag string) *Builder {
	return b.AddField(name, uint32(0), tag)
}

// AddUint64Field adds a uint64 field.
func (b *Builder) AddUint64Field(name string, tag string) *Builder {
	return b.AddField(name, uint64(0), tag)
}

// AddFloat32Field adds a float32 field.
func (b *Builder) AddFloat32Field(name string, tag string) *Builder {
	return b.AddField(name, float32(0), tag)
}

// AddFloat64Field adds a float64 field.
func (b *Builder) AddFloat64Field(name string, tag string) *Builder {
	return b.AddField(name, float64(0), tag)
}

// AddTimeField adds a time.Time field.
func (b *Builder) AddTimeField(name string, tag string) *Builder {
	return b.AddField(name, time.Time{}, tag)
}

// AddSliceField adds a slice field of the given element type.
//
//	builder.AddSliceField("Tags", "", `json:"tags"`) // []string
//	builder.AddSliceField("Scores", 0, `json:"scores"`) // []int
func (b *Builder) AddSliceField(name string, elemTyp any, tag string) *Builder {
	if elemTyp == nil {
		b.errs = append(b.errs, fmt.Errorf("field %q: slice element type must not be nil", name))
		return b
	}
	sliceTyp := reflect.SliceOf(reflect.TypeOf(elemTyp))
	return b.AddField(name, reflect.New(sliceTyp).Elem().Interface(), tag)
}

// AddTypedField adds a field whose type is determined by the type parameter T.
// This is a generic helper that allows adding fields of any type without
// explicitly passing a zero value.
//
//	dynamicstruct.AddTypedField[string](builder, "Name", `json:"name"`)
//	dynamicstruct.AddTypedField[[]int](builder, "Scores", `json:"scores"`)
//	dynamicstruct.AddTypedField[time.Time](builder, "CreatedAt", `json:"created_at"`)
func AddTypedField[T any](b *Builder, name string, tag string) *Builder {
	var zero T
	return b.AddField(name, zero, tag)
}

// RemoveField deletes a field held in Builder.
func (b *Builder) RemoveField(name string) *Builder {
	delete(b.fields, name)
	return b
}

// HasField checks if the field already exists.
func (b *Builder) HasField(name string) bool {
	_, ok := b.fields[name]
	return ok
}

// GetField gets a field that already exists.
// If it does not exist, it returns nil.
func (b *Builder) GetField(name string) *Field {
	if !b.HasField(name) {
		return nil
	}
	return b.fields[name]
}

// Build returns a dynamically created struct definition.
// It returns an error if any validation errors were accumulated during AddField calls.
func (b *Builder) Build() (*DynamicStruct, error) {
	if len(b.errs) > 0 {
		return nil, errors.Join(b.errs...)
	}

	structFields := make([]reflect.StructField, 0, len(b.fields))
	for name, field := range b.fields {
		structFields = append(structFields, reflect.StructField{
			Name: name,
			Type: reflect.TypeOf(field.Type),
			Tag:  reflect.StructTag(field.Tag),
		})
	}

	// Sort fields by name for deterministic ordering.
	sort.Slice(structFields, func(i, j int) bool {
		return structFields[i].Name < structFields[j].Name
	})

	return &DynamicStruct{
		definition: reflect.StructOf(structFields),
	}, nil
}

// DynamicStruct represents a struct created dynamically by Builder.
type DynamicStruct struct {
	definition reflect.Type
}

// New creates a new instance of the dynamic struct.
func (ds *DynamicStruct) New() any {
	return reflect.New(ds.definition).Interface()
}

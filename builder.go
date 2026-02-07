package dynamicstruct

import (
	"errors"
	"fmt"
	"go/token"
	"reflect"
	"sort"
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

package dynamicstruct

import "reflect"

// Builder builds struct dynamically
type Builder struct {
	fields map[string]*Field
}

// DynamicStruct represents a struct created dynamically by Builder
type DynamicStruct struct {
	definition reflect.Type
}

// NewBuilder returns new clean instance of Builder interface
// builder := dynamicstruct.NewStruct()
func NewBuilder() *Builder {
	return &Builder{
		fields: map[string]*Field{},
	}
}

// AddField adds a field to a struct
func (b *Builder) AddField(name string, typ interface{}, tag string) *Builder {
	b.fields[name] = &Field{
		Name: name,
		Type: typ,
		Tag:  tag,
	}

	return b
}

// RemoveField deletes a field held in Builder
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

// Build returns a dynamically created struct
func (b *Builder) Build() *DynamicStruct {
	structFields := make([]reflect.StructField, len(b.fields))
	for name, field := range b.fields {
		structFields = append(structFields, reflect.StructField{
			Name: name,
			Type: reflect.TypeOf(field.Type),
			Tag:  reflect.StructTag(field.Tag),
		})
	}

	return &DynamicStruct{
		definition: reflect.StructOf(structFields),
	}
}

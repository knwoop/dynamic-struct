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

// NewStruct returns new clean instance of Builder interface
// builder := dynamicstruct.NewStruct()
func NewStruct() *Builder {
	return &Builder{
		fields: map[string]*Field{},
	}
}

// AddField adds
func (b *Builder) AddField(name string, typ interface{}, tag string) *Builder {
	b.fields[name] = &Field{
		typ: typ,
		tag: tag,
	}

	return b
}

func (b *Builder) RemoveField(name string) *Builder {
	delete(b.fields, name)

	return b
}

func (b *Builder) HasField(name string) bool {
	_, ok := b.fields[name]
	return ok
}

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
			Type: reflect.TypeOf(field.typ),
			Tag:  reflect.StructTag(field.tag),
		})
	}

	return &DynamicStruct{
		definition: reflect.StructOf(structFields),
	}
}

func (f *Field) SetType(typ interface{}) *Field {
	f.typ = typ
	return f
}

func (f *Field) SetTag(tag string) *Field {
	f.tag = tag
	return f
}

func (ds *dynamicStructImpl) New() interface{} {
	return reflect.New(ds.definition).Interface()
}

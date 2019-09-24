package dynamicstruct

import "reflect"

type (
	// Builder build dynamic struct
	Builder interface {
		// AddField adds new struct's field.
		AddField(name string, typ interface{}, tag string) Builder

		// RemoveField removes existing struct's field.
		// builder.RemoveField("SomeFloatField")
		RemoveField(name string) Builder

		// HasField checks if struct has a field with a given name.
		HasField(name string) bool

		// GetField returns struct's field definition.
		// field := builder.GetField("SomeFloatField")
		GetField(name string) Field

		// Build returns dynamic struct.
		Build() DynamicStruct
	}

	// FieldConfig holds single field's definition.
	// It provides possibility to edit field's type and tag.
	Field interface {
		// SetType changes field's type.
		// Expected value is an instance of golang type.
		//
		// field.SetType([]int{})
		//
		SetType(typ interface{}) Field
		// SetTag changes fields's tag.
		// Expected value is an string which represents classical
		// golang tag.
		//
		// field.SetTag(`yaml:"slice"`)
		//
		SetTag(tag string) Field
	}

	// DynamicStruct contains defined dynamic struct.
	DynamicStruct interface {
		// New provides new instance of defined dynamic struct.
		New() interface{}
	}

	builderImpl struct {
		fields map[string]*fieldConfigImpl
	}

	fieldConfigImpl struct {
		typ interface{}
		tag string
	}

	dynamicStructImpl struct {
		definition reflect.Type
	}
)

// NewStruct returns new clean instance of Builder interface
//
// builder := dynamicstruct.NewStruct()
//
func NewStruct() Builder {
	return &builderImpl{
		fields: map[string]*fieldConfigImpl{},
	}
}

func (b *builderImpl) AddField(name string, typ interface{}, tag string) Builder {
	b.fields[name] = &fieldConfigImpl{
		typ: typ,
		tag: tag,
	}

	return b
}

func (b *builderImpl) RemoveField(name string) Builder {
	delete(b.fields, name)

	return b
}

func (b *builderImpl) HasField(name string) bool {
	_, ok := b.fields[name]
	return ok
}

func (b *builderImpl) GetField(name string) Field {
	if !b.HasField(name) {
		return nil
	}
	return b.fields[name]
}

func (b *builderImpl) Build() DynamicStruct {
	var structFields []reflect.StructField

	for name, field := range b.fields {
		structFields = append(structFields, reflect.StructField{
			Name: name,
			Type: reflect.TypeOf(field.typ),
			Tag:  reflect.StructTag(field.tag),
		})
	}

	return &dynamicStructImpl{
		definition: reflect.StructOf(structFields),
	}
}

func (f *fieldConfigImpl) SetType(typ interface{}) Field {
	f.typ = typ
	return f
}

func (f *fieldConfigImpl) SetTag(tag string) Field {
	f.tag = tag
	return f
}

func (ds *dynamicStructImpl) New() interface{}  {
	return reflect.New(ds.definition).Interface()
}


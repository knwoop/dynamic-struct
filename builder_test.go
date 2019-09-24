package dynamicstruct

import (
	"reflect"
	"testing"
)

func TestNewStruct(t *testing.T) {
	value := NewStruct()

	builder, ok := value.(*builderImpl)
	if !ok {
		t.Errorf(`TestNewStruct - expected instance of *builder got %#v`, value)
	}

	if builder.fields == nil {
		t.Error(`TestNewStruct - expected instance of *map[string]*fieldConfig got nil`)
	}

	if len(builder.fields) > 0 {
		t.Errorf(`TestNewStruct - expected length of fields map to be 0 got %d`, len(builder.fields))
	}
}

func TestBuilderImpl_AddField(t *testing.T) {
	builder := &builderImpl{
		fields: map[string]*fieldConfigImpl{},
	}

	builder.AddField("Field", 1, `key:"value"`)

	field, ok := builder.fields["Field"]
	if !ok {
		t.Error(`TestBuilder_AddField - expected to have field "Field"`)
	}

	expected := &fieldConfigImpl{
		typ: 1,
		tag: `key:"value"`,
	}

	if !reflect.DeepEqual(field, expected) {
		t.Errorf(`TestExtendStruct - expected field to be %#v got %#v`, expected, field)
	}
}

func TestBuilderImpl_RemoveField(t *testing.T) {
	builder := &builderImpl{
		fields: map[string]*fieldConfigImpl{
			"Field": {
				tag: `key:"value"`,
				typ: 1,
			},
		},
	}

	builder.RemoveField("Field")

	if _, ok := builder.fields["Field"]; ok {
		t.Error(`TestBuilder_RemoveField - expected not to have field "Field"`)
	}
}

func TestBuilderImpl_HasField(t *testing.T) {
	builder := &builderImpl{
		fields: map[string]*fieldConfigImpl{},
	}

	if builder.HasField("Field") {
		t.Error(`TestBuilder_HasField - expected not to have field "Field"`)
	}

	builder = &builderImpl{
		fields: map[string]*fieldConfigImpl{
			"Field": {
				tag: `key:"value"`,
				typ: 1,
			},
		},
	}

	if !builder.HasField("Field") {
		t.Error(`TestBuilder_HasField - expected to have field "Field"`)
	}
}

func TestBuilderImpl_GetField(t *testing.T) {
	builder := &builderImpl{
		fields: map[string]*fieldConfigImpl{
			"Field": {
				tag: `key:"value"`,
				typ: 1,
			},
		},
	}

	value := builder.GetField("Field")

	field, ok := value.(*fieldConfigImpl)
	if !ok {
		t.Errorf(`TestBuilder_GetField - expected instance of *fieldConfig got %#v`, value)
	}

	expected := &fieldConfigImpl{
		typ: 1,
		tag: `key:"value"`,
	}

	if !reflect.DeepEqual(field, expected) {
		t.Errorf(`TestExtendStruct - expected field to be %#v got %#v`, expected, field)
	}

	undefined := builder.GetField("Undefined")
	if undefined != nil {
		t.Errorf(`TestBuilder_GetField - expected nil got %#v`, value)
	}
}


package dynamicstruct

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewBuilder(t *testing.T) {
	got := NewBuilder()

	if got.fields == nil {
		t.Error(`TestNewBuilder - expected instance of *map[string]*fieldConfig got nil`)
	}

	if len(got.fields) > 0 {
		t.Errorf(`TestNewGot - expected length of fields map to be 0 got %d`, len(got.fields))
	}
}

func TestBuilder_AddField(t *testing.T) {
	builder := NewBuilder()

	builder.AddField("Field", 1, `key:"value"`)

	got, ok := builder.fields["Field"]
	if !ok {
		t.Fatal(`TestBuilder_AddField - expected to have field "Field"`)
	}

	expected := &Field{
		Name: "Field",
		Type: 1,
		Tag:  `key:"value"`,
	}

	if diff := cmp.Diff(got, expected); diff != "" {
		t.Errorf(`TestBuilder_AddField - (-got +want)\n%v`, diff)
	}
}

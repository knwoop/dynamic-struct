package dynamicstruct

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewBuilder(t *testing.T) {
	got := NewBuilder()

	if got.fields == nil {
		t.Error("expected instance of map[string]*Field got nil")
	}

	if len(got.fields) > 0 {
		t.Errorf("expected length of fields map to be 0 got %d", len(got.fields))
	}
}

func TestBuilder_AddField(t *testing.T) {
	builder := NewBuilder()

	builder.AddField("Field", 1, `key:"value"`)

	got, ok := builder.fields["Field"]
	if !ok {
		t.Fatal(`expected to have field "Field"`)
	}

	expected := &Field{
		Name: "Field",
		Type: 1,
		Tag:  `key:"value"`,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("AddField mismatch: got %+v, want %+v", got, expected)
	}
}

func TestBuilder_RemoveField(t *testing.T) {
	builder := NewBuilder()
	builder.AddField("Field", 1, `key:"value"`)
	builder.RemoveField("Field")

	if _, ok := builder.fields["Field"]; ok {
		t.Error(`expected not to have field "Field"`)
	}
}

func TestBuilder_HasField(t *testing.T) {
	builder := NewBuilder()
	if builder.HasField("Field") {
		t.Error(`expected not to have field "Field"`)
	}

	builder.AddField("Field", 1, `key:"value"`)

	if !builder.HasField("Field") {
		t.Error(`expected to have field "Field"`)
	}
}

func TestBuilder_GetField(t *testing.T) {
	builder := NewBuilder()
	field := builder.GetField("Field")
	if field != nil {
		t.Error("expected nil for non-existent field")
	}

	builder.AddField("Field", 1, `key:"value"`)
	expected := &Field{
		Name: "Field",
		Type: 1,
		Tag:  `key:"value"`,
	}
	got := builder.GetField("Field")
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("GetField mismatch: got %+v, want %+v", got, expected)
	}
}

func TestBuilder_Build(t *testing.T) {
	ds, err := NewBuilder().
		AddField("Name", "", `json:"name"`).
		AddField("Age", 0, `json:"age"`).
		Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	typ := ds.definition

	if typ.Kind() != reflect.Struct {
		t.Fatalf("expected struct kind, got %v", typ.Kind())
	}

	if typ.NumField() != 2 {
		t.Fatalf("expected 2 fields, got %d", typ.NumField())
	}

	// Fields should be sorted by name: Age, Name
	f0 := typ.Field(0)
	if f0.Name != "Age" {
		t.Errorf("expected first field to be Age, got %s", f0.Name)
	}
	if f0.Tag != `json:"age"` {
		t.Errorf("expected tag json:\"age\", got %s", f0.Tag)
	}

	f1 := typ.Field(1)
	if f1.Name != "Name" {
		t.Errorf("expected second field to be Name, got %s", f1.Name)
	}
}

func TestBuilder_New(t *testing.T) {
	ds, err := NewBuilder().
		AddField("Name", "", `json:"name"`).
		AddField("Age", 0, `json:"age"`).
		Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	instance := ds.New()

	// Should be a pointer to a struct
	rv := reflect.ValueOf(instance)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		t.Fatalf("expected pointer to struct, got %v", rv.Type())
	}

	elem := rv.Elem()

	// Set values via reflection
	elem.FieldByName("Name").SetString("Alice")
	elem.FieldByName("Age").SetInt(30)

	// Read them back
	if got := elem.FieldByName("Name").String(); got != "Alice" {
		t.Errorf("expected Name=Alice, got %s", got)
	}
	if got := elem.FieldByName("Age").Int(); got != 30 {
		t.Errorf("expected Age=30, got %d", got)
	}
}

func TestBuilder_New_JSON(t *testing.T) {
	ds, err := NewBuilder().
		AddField("Name", "", `json:"name"`).
		AddField("Age", 0, `json:"age"`).
		Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	instance := ds.New()

	data := `{"name":"Bob","age":25}`
	if err := json.Unmarshal([]byte(data), instance); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	elem := reflect.ValueOf(instance).Elem()
	if got := elem.FieldByName("Name").String(); got != "Bob" {
		t.Errorf("expected Name=Bob, got %s", got)
	}
	if got := elem.FieldByName("Age").Int(); got != 25 {
		t.Errorf("expected Age=25, got %d", got)
	}

	// Marshal back
	out, err := json.Marshal(instance)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(out) != `{"age":25,"name":"Bob"}` {
		t.Errorf("unexpected JSON output: %s", string(out))
	}
}

func TestBuilder_Validation_EmptyName(t *testing.T) {
	_, err := NewBuilder().
		AddField("", "", "").
		Build()
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
	if !strings.Contains(err.Error(), "must not be empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestBuilder_Validation_UnexportedName(t *testing.T) {
	_, err := NewBuilder().
		AddField("lowercase", "", "").
		Build()
	if err == nil {
		t.Fatal("expected error for unexported field name")
	}
	if !strings.Contains(err.Error(), "must be exported") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestBuilder_Validation_NilType(t *testing.T) {
	_, err := NewBuilder().
		AddField("Field", nil, "").
		Build()
	if err == nil {
		t.Fatal("expected error for nil type")
	}
	if !strings.Contains(err.Error(), "must not be nil") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestBuilder_Validation_MultipleErrors(t *testing.T) {
	_, err := NewBuilder().
		AddField("", "", "").
		AddField("lower", "", "").
		AddField("Good", nil, "").
		Build()
	if err == nil {
		t.Fatal("expected errors")
	}
	msg := err.Error()
	if !strings.Contains(msg, "must not be empty") {
		t.Errorf("expected empty name error in: %v", msg)
	}
	if !strings.Contains(msg, "must be exported") {
		t.Errorf("expected unexported error in: %v", msg)
	}
	if !strings.Contains(msg, "must not be nil") {
		t.Errorf("expected nil type error in: %v", msg)
	}
}

func TestBuilder_TypedFieldMethods(t *testing.T) {
	ds, err := NewBuilder().
		AddStringField("Name", `json:"name"`).
		AddBoolField("Active", `json:"active"`).
		AddIntField("Count", "").
		AddInt8Field("I8", "").
		AddInt16Field("I16", "").
		AddInt32Field("I32", "").
		AddInt64Field("I64", `json:"i64"`).
		AddUintField("U", "").
		AddUint8Field("U8", "").
		AddUint16Field("U16", "").
		AddUint32Field("U32", "").
		AddUint64Field("U64", "").
		AddFloat32Field("F32", "").
		AddFloat64Field("F64", `json:"f64"`).
		Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	typ := ds.definition
	if typ.NumField() != 14 {
		t.Fatalf("expected 14 fields, got %d", typ.NumField())
	}

	// Verify types via JSON round-trip for a subset
	instance := ds.New()
	data := `{"name":"Test","active":true,"i64":42,"f64":3.14}`
	if err := json.Unmarshal([]byte(data), instance); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	elem := reflect.ValueOf(instance).Elem()
	if got := elem.FieldByName("Name").String(); got != "Test" {
		t.Errorf("expected Name=Test, got %s", got)
	}
	if got := elem.FieldByName("Active").Bool(); got != true {
		t.Errorf("expected Active=true, got %v", got)
	}
	if got := elem.FieldByName("I64").Int(); got != 42 {
		t.Errorf("expected I64=42, got %d", got)
	}
	if got := elem.FieldByName("F64").Float(); got != 3.14 {
		t.Errorf("expected F64=3.14, got %f", got)
	}
}

func TestBuilder_AddTimeField(t *testing.T) {
	ds, err := NewBuilder().
		AddTimeField("CreatedAt", `json:"created_at"`).
		Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	instance := ds.New()
	data := `{"created_at":"2024-01-15T10:30:00Z"}`
	if err := json.Unmarshal([]byte(data), instance); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	elem := reflect.ValueOf(instance).Elem()
	got := elem.FieldByName("CreatedAt").Interface().(time.Time)
	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("expected CreatedAt=%v, got %v", expected, got)
	}
}

func TestBuilder_AddSliceField(t *testing.T) {
	ds, err := NewBuilder().
		AddSliceField("Tags", "", `json:"tags"`).
		AddSliceField("Scores", 0, `json:"scores"`).
		Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	instance := ds.New()
	data := `{"tags":["go","test"],"scores":[100,200]}`
	if err := json.Unmarshal([]byte(data), instance); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	elem := reflect.ValueOf(instance).Elem()

	tags := elem.FieldByName("Tags")
	if tags.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", tags.Len())
	}
	if got := tags.Index(0).String(); got != "go" {
		t.Errorf("expected tags[0]=go, got %s", got)
	}

	scores := elem.FieldByName("Scores")
	if scores.Len() != 2 {
		t.Fatalf("expected 2 scores, got %d", scores.Len())
	}
	if got := scores.Index(1).Int(); got != 200 {
		t.Errorf("expected scores[1]=200, got %d", got)
	}
}

func TestBuilder_AddSliceField_NilElem(t *testing.T) {
	_, err := NewBuilder().
		AddSliceField("Bad", nil, "").
		Build()
	if err == nil {
		t.Fatal("expected error for nil slice element type")
	}
	if !strings.Contains(err.Error(), "slice element type must not be nil") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAddTypedField(t *testing.T) {
	b := NewBuilder()
	AddTypedField[string](b, "Name", `json:"name"`)
	AddTypedField[int64](b, "Age", `json:"age"`)
	AddTypedField[[]string](b, "Tags", `json:"tags"`)
	AddTypedField[time.Time](b, "CreatedAt", `json:"created_at"`)

	ds, err := b.Build()
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	instance := ds.New()
	data := `{"name":"Alice","age":30,"tags":["a","b"],"created_at":"2024-01-15T10:30:00Z"}`
	if err := json.Unmarshal([]byte(data), instance); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	elem := reflect.ValueOf(instance).Elem()
	if got := elem.FieldByName("Name").String(); got != "Alice" {
		t.Errorf("expected Name=Alice, got %s", got)
	}
	if got := elem.FieldByName("Age").Int(); got != 30 {
		t.Errorf("expected Age=30, got %d", got)
	}
	tags := elem.FieldByName("Tags")
	if tags.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", tags.Len())
	}
	createdAt := elem.FieldByName("CreatedAt").Interface().(time.Time)
	if createdAt.Year() != 2024 {
		t.Errorf("expected year 2024, got %d", createdAt.Year())
	}
}

func TestBuilder_Build_DeterministicOrder(t *testing.T) {
	// Build multiple times and verify the field order is always the same.
	for i := 0; i < 100; i++ {
		ds, err := NewBuilder().
			AddField("Zebra", "", "").
			AddField("Alpha", 0, "").
			AddField("Middle", false, "").
			Build()
		if err != nil {
			t.Fatalf("Build() unexpected error: %v", err)
		}

		typ := ds.definition
		names := make([]string, typ.NumField())
		for j := 0; j < typ.NumField(); j++ {
			names[j] = typ.Field(j).Name
		}
		expected := []string{"Alpha", "Middle", "Zebra"}
		if !reflect.DeepEqual(names, expected) {
			t.Fatalf("iteration %d: expected field order %v, got %v", i, expected, names)
		}
	}
}

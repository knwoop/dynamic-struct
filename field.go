package dynamicstruct

// Field represents a single field in a dynamic struct.
type Field struct {
	Name string
	Type any
	Tag  string
}

package annotation

// EntryHeader represents the metadata for an entry.
type EntryHeader struct {
	Title       string // Title of the entry
	Description string // Description of the entry
}

// EntryFuncType represents a parameter or a result type of a function.
type EntryFuncType struct {
	Name string // Name of the parameter/result
	Type string // Type of the parameter/result
}

// EntryFunc represents a function and its details.
type EntryFunc struct {
	Name       string          // Name of the function
	Parameters []EntryFuncType // Parameters of the function
	Results    []EntryFuncType // Results of the function
}

// Entry represents a single entry parsed from the *ast.File.
type Entry struct {
	Header      EntryHeader // Metadata for the entry
	Comments    []string
	Module      string       // Name of the module where the entry is located
	File        string       // Name of the file where the entry is located
	Path        string       // Path to the file where the entry is located
	Package     string       // Name of the package where the entry is located
	Func        EntryFunc    // Details about the function in the entry
	Struct      string       // Name of the struct in the entry
	Annotations []Annotation // Annotations for the entry
}

func (b *Entry) IsStruct() bool {
	return b.Struct != "" && b.Func.Name == ""
}

func (b *Entry) IsFunc() bool {
	return b.Struct == "" && b.Func.Name != ""
}

func (b *Entry) IsMethod() bool {
	return b.Struct != "" && b.Func.Name != ""
}

package tdconv

// TableSet is a struct of a set of tables.
type TableSet struct {
	Name   string
	Tables []*Table
}

// Table is a struct of table.
type Table struct {
	Name        string
	Columns     []Column
	PKeyColumns []string
	UniqueKeys  []Key
	IndexKeys   []Key
}

// Column is a struct of Column.
type Column struct {
	Name    string
	Type    string
	PKey    bool
	NotNull bool
	Unique  bool
	Index   bool
	Option  string
	Comment string
}

// Key is a struct of Key like Unique Key and Index Key.
type Key struct {
	Name    string
	Columns []string
}

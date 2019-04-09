// Package tdconv converts table definitions to SQL and Go struct etc.
// Currently this package supports SQL and Go format.
// You can change the header and footer text as you want.
// If you want to output with new format, you can do it with creating a new Formatter.
// And more, if the parsed `Table` data are not enough for you, you can modify them as you want.
// For usage of this package, see the standard tool below which uses this package.
// - github.com/takuoki/tdconv/tools/tdconverter
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
	Name     string
	Type     string
	PKey     bool
	NotNull  bool
	Unique   bool
	Index    bool
	Option   string
	Comment  string
	IsCommon bool
}

// Key is a struct of Key like Unique Key and Index Key.
type Key struct {
	Name    string
	Columns []string
}

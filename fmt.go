package tdconv

import "io"

// Formatter is an interface for formatting.
type Formatter interface {
	Fprint(w io.Writer, ts *Table)
}

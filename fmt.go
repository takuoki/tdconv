package tdconv

import (
	"errors"
	"io"
	"os"
)

// Formatter is an interface for formatting.
type Formatter interface {
	Fprint(w io.Writer, t *Table)
}

// Output outputs file(s) using Formatter.
func Output(f Formatter, tables []*Table, multi bool, outdir, extension string) error {

	if !multi {
		return errors.New("I'm sorry. Currently unsupported")
	}

	for _, t := range tables {
		file, err := os.Create(outdir + "/" + t.Name + "." + extension)
		if err != nil {
			return err
		}
		defer file.Close()
		f.Fprint(file, t)
	}

	return nil
}

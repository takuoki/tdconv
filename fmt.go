package tdconv

import (
	"fmt"
	"io"
	"os"

	"github.com/iancoleman/strcase"
)

// Formatter is an interface for formatting.
type Formatter interface {
	Extension() string
	Header() string
	Fprint(w io.Writer, t *Table)
}

// Output outputs file(s) using Formatter.
func Output(f Formatter, tableSet TableSet, multi bool, outdir string) error {

	if !multi {
		err := output(f, tableSet.Tables, fmt.Sprintf("%s/%s.%s", outdir, strcase.ToSnake(tableSet.Name), f.Extension()))
		if err != nil {
			return err
		}
	} else {
		for _, t := range tableSet.Tables {
			err := output(f, []*Table{t}, fmt.Sprintf("%s/%s.%s", outdir, strcase.ToSnake(t.Name), f.Extension()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func output(f Formatter, tables []*Table, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprint(file, f.Header())
	for i, t := range tables {
		f.Fprint(file, t)
		if i < len(tables)-1 {
			fmt.Fprintln(file)
		}
	}
	return nil
}

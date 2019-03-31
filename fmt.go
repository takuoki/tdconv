package tdconv

import (
	"fmt"
	"io"
	"os"
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
		output(f, tableSet.Tables, outdir+"/"+tableSet.Name+"."+f.Extension())
	} else {
		for i := 0; i < len(tableSet.Tables); i++ {
			output(f, tableSet.Tables[i:i+1], outdir+"/"+tableSet.Tables[i].Name+"."+f.Extension())
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

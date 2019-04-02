package tdconv

import (
	"fmt"
	"io"
	"os"

	"github.com/iancoleman/strcase"
)

// Formatter is an interface for formatting.
type Formatter interface {
	Header() func(w io.Writer, tableSet TableSet)
	TableHeader() func(w io.Writer, table *Table)
	TableFooter() func(w io.Writer, table *Table)
	Footer() func(w io.Writer, tableSet TableSet)
	Extension() string
	Fprint(w io.Writer, t *Table)
}

type formatter struct {
	header      func(w io.Writer, tableSet TableSet)
	tableHeader func(w io.Writer, table *Table)
	tableFooter func(w io.Writer, table *Table)
	footer      func(w io.Writer, tableSet TableSet)
}

func (f *formatter) Header() func(w io.Writer, tableSet TableSet) {
	return f.header
}

func (f *formatter) TableHeader() func(w io.Writer, table *Table) {
	return f.tableHeader
}

func (f *formatter) TableFooter() func(w io.Writer, table *Table) {
	return f.tableFooter
}

func (f *formatter) Footer() func(w io.Writer, tableSet TableSet) {
	return f.footer
}

func (f *formatter) setHeader(fc func(w io.Writer, tableSet TableSet)) {
	f.header = fc
}

func (f *formatter) setTableHeader(fc func(w io.Writer, table *Table)) {
	f.tableHeader = fc
}

func (f *formatter) setTableFooter(fc func(w io.Writer, table *Table)) {
	f.tableFooter = fc
}

func (f *formatter) setFooter(fc func(w io.Writer, tableSet TableSet)) {
	f.footer = fc
}

// Output outputs file(s) using Formatter.
func Output(f Formatter, tableSet TableSet, multi bool, outdir string) error {

	if !multi {
		err := output(f, tableSet, 0, len(tableSet.Tables), fmt.Sprintf("%s/%s.%s", outdir, strcase.ToSnake(tableSet.Name), f.Extension()))
		if err != nil {
			return err
		}
	} else {
		for i := 0; i < len(tableSet.Tables); i++ {
			err := output(f, tableSet, i, i+1, fmt.Sprintf("%s/%s.%s", outdir, strcase.ToSnake(tableSet.Tables[i].Name), f.Extension()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func output(f Formatter, tableSet TableSet, from, to int, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if fc := f.Header(); fc != nil {
		fc(file, tableSet)
	}
	for i := from; i < to; i++ {
		if fc := f.TableHeader(); fc != nil {
			fc(file, tableSet.Tables[i])
		}
		f.Fprint(file, tableSet.Tables[i])
		if fc := f.TableFooter(); fc != nil {
			fc(file, tableSet.Tables[i])
		}
		if i < to-1 {
			fmt.Fprintln(file)
		}
	}
	if fc := f.Footer(); fc != nil {
		fc(file, tableSet)
	}
	return nil
}

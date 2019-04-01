package tdconv

import (
	"fmt"
	"io"
	"os"

	"github.com/iancoleman/strcase"
)

// Formatter is an interface for formatting.
type Formatter interface {
	Header() func(w io.Writer, tables []*Table)
	TableHeader() func(w io.Writer, table *Table)
	TableFooter() func(w io.Writer, table *Table)
	Footer() func(w io.Writer, tables []*Table)
	Extension() string
	Fprint(w io.Writer, t *Table)
}

type formatter struct {
	header      func(w io.Writer, tables []*Table)
	tableHeader func(w io.Writer, table *Table)
	tableFooter func(w io.Writer, table *Table)
	footer      func(w io.Writer, tables []*Table)
}

func (f *formatter) Header() func(w io.Writer, tables []*Table) {
	return f.header
}

func (f *formatter) TableHeader() func(w io.Writer, table *Table) {
	return f.tableHeader
}

func (f *formatter) TableFooter() func(w io.Writer, table *Table) {
	return f.tableFooter
}

func (f *formatter) Footer() func(w io.Writer, tables []*Table) {
	return f.footer
}

func (f *formatter) setHeader(fc func(w io.Writer, tables []*Table)) {
	f.header = fc
}

func (f *formatter) setTableHeader(fc func(w io.Writer, table *Table)) {
	f.tableHeader = fc
}

func (f *formatter) setTableFooter(fc func(w io.Writer, table *Table)) {
	f.tableFooter = fc
}

func (f *formatter) setFooter(fc func(w io.Writer, tables []*Table)) {
	f.footer = fc
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

	if fc := f.Header(); fc != nil {
		fc(file, tables)
	}
	for i, t := range tables {
		if fc := f.TableHeader(); fc != nil {
			fc(file, t)
		}
		f.Fprint(file, t)
		if fc := f.TableFooter(); fc != nil {
			fc(file, t)
		}
		if i < len(tables)-1 {
			fmt.Fprintln(file)
		}
	}
	if fc := f.Footer(); fc != nil {
		fc(file, tables)
	}
	return nil
}

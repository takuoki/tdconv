package tdconv

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/iancoleman/strcase"
)

var createFile func(name string) (io.WriteCloser, error)

func init() {
	createFile = func(name string) (io.WriteCloser, error) {
		return os.Create(name)
	}
}

// Formatter is an interface for formatting.
type Formatter interface {
	Fprint(w io.Writer, t *Table)
	Header(w io.Writer, ts *TableSet)
	TableHeader(w io.Writer, t *Table)
	TableFooter(w io.Writer, t *Table)
	Footer(w io.Writer, ts *TableSet)
	Extension() string
}

type formatter struct {
	header      func(w io.Writer, ts *TableSet)
	tableHeader func(w io.Writer, t *Table)
	tableFooter func(w io.Writer, t *Table)
	footer      func(w io.Writer, ts *TableSet)
}

func (f *formatter) Header(w io.Writer, ts *TableSet) {
	if f.header != nil {
		f.header(w, ts)
	}
}

func (f *formatter) TableHeader(w io.Writer, t *Table) {
	if f.tableHeader != nil {
		f.tableHeader(w, t)
	}
}

func (f *formatter) TableFooter(w io.Writer, t *Table) {
	if f.tableFooter != nil {
		f.tableFooter(w, t)
	}
}

func (f *formatter) Footer(w io.Writer, ts *TableSet) {
	if f.footer != nil {
		f.footer(w, ts)
	}
}

func (f *formatter) setHeader(fc func(w io.Writer, ts *TableSet)) {
	f.header = fc
}

func (f *formatter) setTableHeader(fc func(w io.Writer, t *Table)) {
	f.tableHeader = fc
}

func (f *formatter) setTableFooter(fc func(w io.Writer, t *Table)) {
	f.tableFooter = fc
}

func (f *formatter) setFooter(fc func(w io.Writer, ts *TableSet)) {
	f.footer = fc
}

// Output outputs file(s) using Formatter.
func Output(f Formatter, ts *TableSet, multi bool, outdir string) error {

	if ts == nil {
		return errors.New("Table set is nil")
	}

	if !multi {
		err := output(f, ts, 0, len(ts.Tables), fmt.Sprintf("%s/%s.%s", outdir, strcase.ToSnake(ts.Name), f.Extension()))
		if err != nil {
			return err
		}
	} else {
		for i := 0; i < len(ts.Tables); i++ {
			err := output(f, ts, i, i+1, fmt.Sprintf("%s/%s.%s", outdir, strcase.ToSnake(ts.Tables[i].Name), f.Extension()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func output(f Formatter, ts *TableSet, from, to int, filepath string) error {
	file, err := createFile(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	f.Header(file, ts)
	for i := from; i < to; i++ {
		f.TableHeader(file, ts.Tables[i])
		f.Fprint(file, ts.Tables[i])
		f.TableFooter(file, ts.Tables[i])
		if i < to-1 {
			fmt.Fprintln(file)
		}
	}
	f.Footer(file, ts)

	return nil
}

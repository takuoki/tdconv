package tdconv

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

const (
	importStr = `import(
	"time"
)

`
)

// GoFormatter is a formatter to output the table definision as Go struct.
type GoFormatter struct {
	formatter
}

// NewGoFormatter creates a new GoFormatter.
// You can change some parameters of the GoFormatter with GoFormatOption.
func NewGoFormatter(options ...GoFormatOption) (*GoFormatter, error) {

	f := GoFormatter{}
	f.setHeader(func(w io.Writer, tables []*Table) {
		fmt.Fprint(w,
			"// This file generated by tdconv. DO NOT EDIT.\n"+
				"// See more details at https://github.com/takuoki/tdconv.\n"+
				"package main\n\n"+importStr)
	})
	for _, opt := range options {
		err := opt(&f)
		if err != nil {
			return nil, err
		}
	}
	return &f, nil
}

// GoFormatOption changes some parameters of the GoFormatter.
type GoFormatOption func(*GoFormatter) error

// GoHeader changes the header.
func GoHeader(fc func(w io.Writer, tables []*Table)) GoFormatOption {
	return func(f *GoFormatter) error {
		f.setHeader(fc)
		return nil
	}
}

// GoTableHeader changes the header of each table.
func GoTableHeader(fc func(w io.Writer, table *Table)) GoFormatOption {
	return func(f *GoFormatter) error {
		f.setTableHeader(fc)
		return nil
	}
}

// GoTableFooter changes the footer of each table.
func GoTableFooter(fc func(w io.Writer, table *Table)) GoFormatOption {
	return func(f *GoFormatter) error {
		f.setTableFooter(fc)
		return nil
	}
}

// GoFooter changes the footer.
func GoFooter(fc func(w io.Writer, tables []*Table)) GoFormatOption {
	return func(f *GoFormatter) error {
		f.setFooter(fc)
		return nil
	}
}

// Extension returns the extension of Go file.
func (f *GoFormatter) Extension() string {
	return "go"
}

// Fprint outputs the table definision as Go struct.
func (f *GoFormatter) Fprint(w io.Writer, t *Table) {

	if f == nil || t == nil {
		return
	}

	fmt.Fprintf(w, "type %s struct {\n", f.structName(t.Name))

	for _, c := range t.Columns {
		fmt.Fprintf(w, "\t%s %s\n", f.propertyName(c.Name), f.convType(c.Type))
	}

	fmt.Fprintln(w, "}")
}

func (f *GoFormatter) structName(s string) string {
	return strcase.ToLowerCamel(s)
}

func (f *GoFormatter) propertyName(s string) string {
	return f.convWords(strcase.ToCamel(s))
}

var tRegexp = regexp.MustCompile("^([a-zA-Z]+)[ (]{1}.*$")

func (f *GoFormatter) convType(t string) string {
	var r string
	switch strings.ToUpper(tRegexp.ReplaceAllString(t, "$1")) {
	case "INT", "TINYINT", "BIGINT":
		r = "*int"
	case "DOUBLE":
		r = "*float32"
	case "CHAR", "VARCHAR", "TEXT", "ENUM":
		r = "*string"
	case "BOOLEAN":
		r = "*bool"
	case "TIMESTAMP", "DATE", "TIME":
		r = "*time.Time"
	case "DECIMAL":
		r = "*decimal.Decimal"
	default:
		r = "UNKNOWN"
	}
	return r
}

func (f *GoFormatter) convWords(s string) string {
	s = strings.Replace(s, "Id", "ID", -1)
	s = strings.Replace(s, "Ip", "IP", -1)
	s = strings.Replace(s, "Api", "API", -1)
	s = strings.Replace(s, "Url", "URL", -1)
	s = strings.Replace(s, "Http", "HTTP", -1)
	return s
}

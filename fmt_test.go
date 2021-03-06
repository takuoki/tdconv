package tdconv_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/takuoki/tdconv"
)

type testFormatter struct{}

func (*testFormatter) Fprint(w io.Writer, t *tdconv.Table) {
	fmt.Fprintf(w, "table contents: %s\n", t.Name)
}

func (*testFormatter) Header(w io.Writer, tableSet *tdconv.TableSet) {
	fmt.Fprintf(w, "header: %s\n", tableSet.Name)
}

func (*testFormatter) TableHeader(w io.Writer, table *tdconv.Table) {
	fmt.Fprintf(w, "table header: %s\n", table.Name)
}

func (*testFormatter) TableFooter(w io.Writer, table *tdconv.Table) {
	fmt.Fprintf(w, "table footer: %s\n", table.Name)
}

func (*testFormatter) Footer(w io.Writer, tableSet *tdconv.TableSet) {
	fmt.Fprintf(w, "footer: %s\n", tableSet.Name)
}

func (*testFormatter) Extension() string {
	return "test"
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func TestOutput(t *testing.T) {

	var outputMap map[string]*bytes.Buffer
	errFileName := "output_dir/error.test"
	createFile := func(name string) (io.WriteCloser, error) {
		if name == errFileName {
			return nil, errors.New("error")
		}
		b := &bytes.Buffer{}
		outputMap[name] = b
		return nopCloser{b}, nil
	}
	resetFunc := tdconv.SetCreateFile(createFile)
	defer resetFunc()

	optFunc := func(s string) func(io.Writer, *tdconv.TableSet) {
		return func(w io.Writer, _ *tdconv.TableSet) {
			fmt.Fprintln(w, s)
		}
	}
	optTableFunc := func(s string) func(io.Writer, *tdconv.Table) {
		return func(w io.Writer, _ *tdconv.Table) {
			fmt.Fprintln(w, s)
		}
	}

	cases := []struct {
		caseName string
		f        tdconv.Formatter
		tableSet *tdconv.TableSet
		multi    bool
		expected map[string]string
		errMsg   string
	}{
		{
			caseName: "success: non-multi",
			f:        &testFormatter{},
			tableSet: &tdconv.TableSet{
				Name: "sample_table_set",
				Tables: []*tdconv.Table{
					{Name: "sample_table_1"},
					{Name: "sample_table_2"},
				},
			},
			expected: map[string]string{
				"output_dir/sample_table_set.test": "header: sample_table_set\n" +
					"table header: sample_table_1\n" +
					"table contents: sample_table_1\n" +
					"table footer: sample_table_1\n\n" +
					"table header: sample_table_2\n" +
					"table contents: sample_table_2\n" +
					"table footer: sample_table_2\n" +
					"footer: sample_table_set\n",
			},
		},
		{
			caseName: "success: multi",
			f:        &testFormatter{},
			tableSet: &tdconv.TableSet{
				Name: "sample_table_set",
				Tables: []*tdconv.Table{
					{Name: "sample_table_1"},
					{Name: "sample_table_2"},
				},
			},
			multi: true,
			expected: map[string]string{
				"output_dir/sample_table_1.test": "header: sample_table_set\n" +
					"table header: sample_table_1\n" +
					"table contents: sample_table_1\n" +
					"table footer: sample_table_1\n" +
					"footer: sample_table_set\n",
				"output_dir/sample_table_2.test": "header: sample_table_set\n" +
					"table header: sample_table_2\n" +
					"table contents: sample_table_2\n" +
					"table footer: sample_table_2\n" +
					"footer: sample_table_set\n",
			},
		},
		{
			caseName: "success: SQL formatter",
			f: mustSQLFormatter(
				tdconv.SQLTableHeader(optTableFunc("# table header")),
				tdconv.SQLTableFooter(optTableFunc("# table footer")),
				tdconv.SQLFooter(optFunc("# footer")),
			),
			tableSet: &tdconv.TableSet{
				Name: "sample_table_set",
				Tables: []*tdconv.Table{
					{
						Name: "sample_table",
						Columns: []tdconv.Column{
							{Name: "id", Type: "INT UNSIGNED", PKey: true, NotNull: true, Unique: false, Index: false, Option: "AUTO_INCREMENT", Comment: "this is id!", IsCommon: false},
							{Name: "foo", Type: "VARCHAR(32)", PKey: false, NotNull: true, Unique: true, Index: false, Option: "", Comment: "", IsCommon: false},
						},
						PKeyColumns: []string{"id"},
					},
				},
			},
			expected: map[string]string{
				"output_dir/sample_table_set.sql": "# This file generated by tdconv. DO NOT EDIT.\n" +
					"# See more details at https://github.com/takuoki/tdconv.\n" +
					"# table header\n" +
					"DROP TABLE IF EXISTS sample_table;\n" +
					"CREATE TABLE `sample_table` (\n" +
					"    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'this is id!',\n" +
					"    `foo` VARCHAR(32) NOT NULL UNIQUE,\n" +
					"    PRIMARY KEY (id)\n" +
					");\n" +
					"# table footer\n" +
					"# footer\n",
			},
		},
		{
			caseName: "success: Go formatter",
			f: mustGoFormatter(
				tdconv.GoTableHeader(optTableFunc("// table header")),
				tdconv.GoTableFooter(optTableFunc("// table footer")),
				tdconv.GoFooter(optFunc("// footer")),
			),
			tableSet: &tdconv.TableSet{
				Name: "sample_table_set",
				Tables: []*tdconv.Table{
					{
						Name: "sample_table",
						Columns: []tdconv.Column{
							{Name: "id", Type: "INT UNSIGNED", PKey: true, NotNull: true, Unique: false, Index: false, Option: "AUTO_INCREMENT", Comment: "this is id!", IsCommon: false},
							{Name: "foo", Type: "VARCHAR(32)", PKey: false, NotNull: true, Unique: true, Index: false, Option: "", Comment: "", IsCommon: false},
						},
						PKeyColumns: []string{"id"},
					},
				},
			},
			expected: map[string]string{
				"output_dir/sample_table_set.go": "// This file generated by tdconv. DO NOT EDIT.\n" +
					"// See more details at https://github.com/takuoki/tdconv.\n" +
					"package main\n\nimport(\n\t\"time\"\n)\n\n" +
					"// table header\n" +
					"type sampleTable struct {\n" +
					"	ID *int\n" +
					"	Foo *string\n" +
					"}\n" +
					"// table footer\n" +
					"// footer\n",
			},
		},
		{
			caseName: "failure: table set is nil",
			f:        &testFormatter{},
			tableSet: nil,
			errMsg:   "Table set is nil",
		},
		{
			caseName: "failure: error case of non-multi",
			f:        &testFormatter{},
			tableSet: &tdconv.TableSet{
				Name: "error",
				Tables: []*tdconv.Table{
					{Name: "sample_table_1"},
					{Name: "sample_table_2"},
				},
			},
			errMsg: "error",
		},
		{
			caseName: "failure: error case of multi",
			f:        &testFormatter{},
			tableSet: &tdconv.TableSet{
				Name: "sample_table_set",
				Tables: []*tdconv.Table{
					{Name: "error"},
				},
			},
			multi:  true,
			errMsg: "error",
		},
	}

	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			outputMap = map[string]*bytes.Buffer{}
			err := tdconv.Output(c.f, c.tableSet, c.multi, "output_dir")

			if c.errMsg == "" {
				if err != nil {
					t.Errorf("error must not occur: %v", err)
					return
				}
				if len(c.expected) != len(outputMap) {
					t.Errorf("the number of output files doesn't match (expected=%d, actual=%d)", len(c.expected), len(outputMap))
					return
				}
				for k, v := range c.expected {
					if a, ok := outputMap[k]; !ok {
						t.Errorf("file doesn't exist in output files (filename=%s)", k)
						return
					} else if a.String() != v {
						t.Errorf("file contents don't match (expected=%s, actual=%s)", v, a.String())
						return
					}
				}
			} else {
				if err == nil {
					t.Errorf("error must occur")
					return
				}
				if endIndex := strings.Index(err.Error(), ":"); endIndex < 0 {
					if err.Error() != c.errMsg {
						t.Errorf("error message doesn't match (expected=%s, actual=%s)", c.errMsg, err.Error())
						return
					}
				} else if err.Error()[:endIndex] != c.errMsg {
					t.Errorf("error message doesn't match (expected=%s, actual=%s)", c.errMsg, err.Error()[:endIndex])
					return
				}
			}
		})
	}
}

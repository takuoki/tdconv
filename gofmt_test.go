package tdconv_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/takuoki/tdconv"
)

func TestNewGoFormatter(t *testing.T) {

	errOptionFunc := func(*tdconv.GoFormatter) error {
		return errors.New("error")
	}

	cases := []struct {
		caseName string
		opts     []tdconv.GoFormatOption
		errMsg   string
	}{
		{
			caseName: "success: default",
		},
		{
			caseName: "success: set all options",
			opts: []tdconv.GoFormatOption{
				tdconv.GoHeader(nil),
				tdconv.GoTableHeader(nil),
				tdconv.GoTableFooter(nil),
				tdconv.GoFooter(nil),
			},
		},
		{
			caseName: "failure: option error",
			opts:     []tdconv.GoFormatOption{errOptionFunc},
			errMsg:   "error",
		},
	}

	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			_, err := tdconv.NewGoFormatter(c.opts...)

			if c.errMsg == "" {
				if err != nil {
					t.Errorf("error must not occur: %v", err)
					return
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

func TestGoFormatter_Extension(t *testing.T) {
	var f *tdconv.GoFormatter
	if f.Extension() != "go" {
		t.Errorf("value doesn't match (expected=go, actual=%s)", f.Extension())
	}
}

func TestGoFormatter_Fprint(t *testing.T) {

	mustGoFormatter := func(options ...tdconv.GoFormatOption) *tdconv.GoFormatter {
		f, err := tdconv.NewGoFormatter(options...)
		if err != nil {
			panic(err)
		}
		return f
	}

	cases := []struct {
		caseName string
		f        *tdconv.GoFormatter
		t        *tdconv.Table
		expected string
	}{
		{
			caseName: "nil formatter",
			f:        nil,
			t: &tdconv.Table{
				Name: "sample_table",
				Columns: []tdconv.Column{
					{Name: "id", Type: "INT UNSIGNED", PKey: true, NotNull: true, Unique: false, Index: false, Option: "AUTO_INCREMENT", Comment: "this is id!", IsCommon: false},
					{Name: "foo", Type: "VARCHAR(32)", PKey: false, NotNull: true, Unique: true, Index: false, Option: "", Comment: "", IsCommon: false},
					{Name: "bar", Type: "VARCHAR(32)", PKey: false, NotNull: false, Unique: false, Index: true, Option: "", Comment: "", IsCommon: false},
					{Name: "created_at", Type: "TIMESTAMP NULL", PKey: false, NotNull: false, Unique: false, Index: false, Option: "DEFAULT CURRENT_TIMESTAMP", Comment: "", IsCommon: true},
					{Name: "updated_at", Type: "TIMESTAMP NULL", PKey: false, NotNull: false, Unique: false, Index: false, Option: "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP", Comment: "", IsCommon: true},
					{Name: "deleted_at", Type: "TIMESTAMP NULL", PKey: false, NotNull: false, Unique: false, Index: false, Option: "", Comment: "", IsCommon: true},
				},
				PKeyColumns: []string{"id"},
				UniqueKeys:  nil,
				IndexKeys:   []tdconv.Key{{Name: "bar_key", Columns: []string{"bar"}}},
			},
			expected: "",
		},
		{
			caseName: "table is nil",
			f:        mustGoFormatter(),
			t:        nil,
			expected: "",
		},
		{
			caseName: "standard output",
			f:        mustGoFormatter(),
			t: &tdconv.Table{
				Name: "sample_table",
				Columns: []tdconv.Column{
					{Name: "id", Type: "INT UNSIGNED", PKey: true, NotNull: true, Unique: false, Index: false, Option: "AUTO_INCREMENT", Comment: "this is id!", IsCommon: false},
					{Name: "foo", Type: "VARCHAR(32)", PKey: false, NotNull: true, Unique: true, Index: false, Option: "", Comment: "", IsCommon: false},
					{Name: "bar", Type: "VARCHAR(32)", PKey: false, NotNull: false, Unique: false, Index: true, Option: "", Comment: "", IsCommon: false},
					{Name: "created_at", Type: "TIMESTAMP NULL", PKey: false, NotNull: false, Unique: false, Index: false, Option: "DEFAULT CURRENT_TIMESTAMP", Comment: "", IsCommon: true},
					{Name: "updated_at", Type: "TIMESTAMP NULL", PKey: false, NotNull: false, Unique: false, Index: false, Option: "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP", Comment: "", IsCommon: true},
					{Name: "deleted_at", Type: "TIMESTAMP NULL", PKey: false, NotNull: false, Unique: false, Index: false, Option: "", Comment: "", IsCommon: true},
				},
				PKeyColumns: []string{"id"},
				UniqueKeys:  nil,
				IndexKeys:   []tdconv.Key{{Name: "bar_key", Columns: []string{"bar"}}},
			},
			expected: "type sampleTable struct {\n" +
				"	ID *int\n" +
				"	Foo *string\n" +
				"	Bar *string\n" +
				"	CreatedAt *time.Time\n" +
				"	UpdatedAt *time.Time\n" +
				"	DeletedAt *time.Time\n" +
				"}\n",
		},
		{
			caseName: "unique key",
			f:        mustGoFormatter(),
			t: &tdconv.Table{
				Name: "sample_table",
				Columns: []tdconv.Column{
					{Name: "id", Type: "INT UNSIGNED", PKey: true, NotNull: true, Unique: false, Index: false, Option: "AUTO_INCREMENT", Comment: "this is id!", IsCommon: false},
					{Name: "foo", Type: "VARCHAR(32)", PKey: false, NotNull: true, Unique: true, Index: false, Option: "", Comment: "", IsCommon: false},
					{Name: "bar", Type: "VARCHAR(32)", PKey: false, NotNull: false, Unique: false, Index: false, Option: "", Comment: "", IsCommon: false},
					{Name: "baz", Type: "VARCHAR(32)", PKey: false, NotNull: false, Unique: false, Index: false, Option: "", Comment: "", IsCommon: false},
				},
				PKeyColumns: []string{"id"},
				UniqueKeys:  []tdconv.Key{{Name: "bar_key", Columns: []string{"bar", "baz"}}},
			},
			expected: "type sampleTable struct {\n" +
				"	ID *int\n" +
				"	Foo *string\n" +
				"	Bar *string\n" +
				"	Baz *string\n" +
				"}\n",
		},
		{
			caseName: "type double, boolean, unknown",
			f:        mustGoFormatter(),
			t: &tdconv.Table{
				Name: "sample_table",
				Columns: []tdconv.Column{
					{Name: "id", Type: "INT UNSIGNED", PKey: true, NotNull: true, Unique: false, Index: false, Option: "AUTO_INCREMENT", Comment: "this is id!", IsCommon: false},
					{Name: "foo", Type: "DOUBLE", PKey: false, NotNull: true, Unique: true, Index: false, Option: "", Comment: "", IsCommon: false},
					{Name: "bar", Type: "BOOLEAN", PKey: false, NotNull: false, Unique: false, Index: false, Option: "", Comment: "", IsCommon: false},
					{Name: "baz", Type: "LONGBLOB", PKey: false, NotNull: false, Unique: false, Index: false, Option: "", Comment: "", IsCommon: false},
				},
				PKeyColumns: []string{"id"},
				UniqueKeys:  []tdconv.Key{{Name: "bar_key", Columns: []string{"bar1", "bar2"}}},
			},
			expected: "type sampleTable struct {\n" +
				"	ID *int\n" +
				"	Foo *float32\n" +
				"	Bar *bool\n" +
				"	Baz UNKNOWN\n" +
				"}\n",
		},
	}

	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {

			b := &bytes.Buffer{}
			c.f.Fprint(b, c.t)

			if b.String() != c.expected {
				t.Errorf("value doesn't match (expected=%s, actual=%s)", c.expected, b.String())
			}
		})
	}
}

package tdconv_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/takuoki/tdconv"
)

func TestNewSQLFormatter(t *testing.T) {

	errOptionFunc := func(*tdconv.SQLFormatter) error {
		return errors.New("error")
	}

	cases := []struct {
		caseName string
		opts     []tdconv.SQLFormatOption
		errMsg   string
	}{
		{
			caseName: "success: default",
		},
		{
			caseName: "success: set all options",
			opts: []tdconv.SQLFormatOption{
				tdconv.SQLHeader(nil),
				tdconv.SQLTableHeader(nil),
				tdconv.SQLTableFooter(nil),
				tdconv.SQLFooter(nil),
			},
		},
		{
			caseName: "failure: option error",
			opts:     []tdconv.SQLFormatOption{errOptionFunc},
			errMsg:   "error",
		},
	}

	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			_, err := tdconv.NewSQLFormatter(c.opts...)

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

func TestSQLFormatter_Extension(t *testing.T) {
	var f *tdconv.SQLFormatter
	if f.Extension() != "sql" {
		t.Errorf("value doesn't match (expected=sql, actual=%s)", f.Extension())
	}
}

func TestSQLFormatter_Fprint(t *testing.T) {

	mustSQLFormatter := func(options ...tdconv.SQLFormatOption) *tdconv.SQLFormatter {
		f, err := tdconv.NewSQLFormatter(options...)
		if err != nil {
			panic(err)
		}
		return f
	}

	cases := []struct {
		caseName string
		f        *tdconv.SQLFormatter
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
			f:        mustSQLFormatter(),
			t:        nil,
			expected: "",
		},
		{
			caseName: "standard output",
			f:        mustSQLFormatter(),
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
			expected: "DROP TABLE IF EXISTS sample_table;\n" +
				"CREATE TABLE `sample_table` (\n" +
				"    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'this is id!',\n" +
				"    `foo` VARCHAR(32) NOT NULL UNIQUE,\n" +
				"    `bar` VARCHAR(32),\n" +
				"    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,\n" +
				"    `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n" +
				"    `deleted_at` TIMESTAMP NULL,\n" +
				"    PRIMARY KEY (id),\n" +
				"    INDEX `bar_key` (bar)\n" +
				");\n",
		},
		{
			caseName: "unique key",
			f:        mustSQLFormatter(),
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
			expected: "DROP TABLE IF EXISTS sample_table;\n" +
				"CREATE TABLE `sample_table` (\n" +
				"    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'this is id!',\n" +
				"    `foo` VARCHAR(32) NOT NULL UNIQUE,\n" +
				"    `bar` VARCHAR(32),\n" +
				"    `baz` VARCHAR(32),\n" +
				"    PRIMARY KEY (id),\n" +
				"    UNIQUE KEY `bar_key` (bar, baz)\n" +
				");\n",
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

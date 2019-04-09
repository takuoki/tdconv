# tdconv

[![CircleCI](https://circleci.com/gh/takuoki/tdconv/tree/master.svg?style=shield&circle-token=9c99add95b184cb77460481d93e0e7e5d9f7a943)](https://circleci.com/gh/takuoki/tdconv/tree/master)
[![codecov](https://codecov.io/gh/takuoki/tdconv/branch/master/graph/badge.svg)](https://codecov.io/gh/takuoki/tdconv)
[![GoDoc](https://godoc.org/github.com/takuoki/tdconv?status.svg)](https://godoc.org/github.com/takuoki/tdconv)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

A golang package and tool for converting table definitions to SQL and Go struct etc.

<!-- vscode-markdown-toc -->
* [Description](#Description)
	* [Package `tdconv`](#Packagetdconv)
	* [Tool `tdconverter`](#Tooltdconverter)
* [Install](#Install)
* [Requirements](#Requirements)
* [Usage](#Usage)

<!-- vscode-markdown-toc-config
	numbering=false
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

## <a name='Description'></a>Description

### <a name='Packagetdconv'></a>Package `tdconv`

This package converts table definitions to SQL and Go struct etc.
Currently this package supports SQL and Go format.

For example, if the table definition is like...

![sample table](docs/images/sample_table.png)

SQL is output as follows.

```sql
# This file generated by tdconv. DO NOT EDIT.
# See more details at https://github.com/takuoki/tdconv.
DROP TABLE IF EXISTS sample_table;
CREATE TABLE `sample_table` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'this is id!',
    `foo` VARCHAR(32) NOT NULL UNIQUE,
    `bar` VARCHAR(32),
    PRIMARY KEY (id),
    INDEX `bar_key` (bar)
);
```

Go struct is output as follows.

```go
// This file generated by tdconv. DO NOT EDIT.
// See more details at https://github.com/takuoki/tdconv.
package main

import(
  "time"
)

type sampleTable struct {
  ID *int
  Foo *string
  Bar *string
}
```

You can change the header and footer text as you want.
If you want to output with new format, you can do it with creating a new Formatter.
And more, if the parsed `Table` data are not enough for you, you can modify them as you want.
For usage of this package, see [the standard tool](tools/tdconverter) which uses this package.

### <a name='Tooltdconverter'></a>Tool `tdconverter`

This is a standard tool uses the package `tdconv`.
For more details, see [`README.md`](tools/tdconverter/README.md).

## <a name='Install'></a>Install

You can install this tool using `go get`.
Before installing, enable the Go module feature.

```bash
go get github.com/takuoki/tdconv/tools/tdconverter
```

## <a name='Requirements'></a>Requirements

On parsing, this package uses `Sheet` objects in [`github.com/takuoki/gsheets`](https://github.com/takuoki/gsheets) package.
For more details, see `README.md` in this package.

## <a name='Usage'></a>Usage

First, create a new `Parser`.
If your sheet is different from the default sheet format, set `ParseOption` to match your sheet.
And if you need some common columns to all tables, set them with `SetCommonColumns` method.

```go
p, err := tdconv.NewParser()
if err != nil {
  return nil, fmt.Errorf("Unable to create new parser: %v", err)
}

err = p.SetCommonColumns(commonSheet)
if err != nil {
  return nil, fmt.Errorf("Unable to parse common sheet information: %v", err)
}
```

Then, parse your sheet with `Parse` method.
Basically, just specify the sheet value returns by `GetSheet` method of the `gsheets` package.
In case of parsing multiple sheets, loop it in your application.

```go
var tables []*tdconv.Table
for _, sheetname := range sheets {
  sheet, err := gclient.GetSheet(ctx, id, sheetname)
  if err != nil {
    return nil, fmt.Errorf("Unable to get sheet values (sheetname=%s): %v", sheetname, err)
  }
  table, err := p.Parse(sheet)
  if err != nil {
    return nil, fmt.Errorf("Unable to parse sheet information (sheetname=%s): %v", sheetname, err)
  }
  tables = append(tables, table)
}
```

Finally, create `TableSet` based on some `Table`s you get above step, and output file(s) with formatter you need.
For `SQLFormatter` or `GoFormatter`, you can change the header and footer text with `SQLFormatOption` or `GoFormatOption`.
If the parsed `Table` data are not enough for you, you can modify them as you want before calling the `Output` function.
If the `multi` flag, which is one of the arguments of the `Output` function, is `true`, one file is output for each `Table`.
If `false`, one file is output for each `TableSet`.

```go
tableSet := &tdconv.TableSet{
  Name:   title,
  Tables: tables,
}

f, err := tdconv.NewSQLFormatter()
if err != nil {
  return fmt.Errorf("Fail to create SQL formatter: %v", err)
}

if err := tdconv.Output(f, tableSet, multi, outdir); err != nil {
  return fmt.Errorf("Fail to output table definitions: %v", err)
}
```

If you create a new formatter, follow the `Formatter` interface below.

```go
type Formatter interface {
  Fprint(w io.Writer, t *Table)
  Header(w io.Writer, ts *TableSet)
  TableHeader(w io.Writer, t *Table)
  TableFooter(w io.Writer, t *Table)
  Footer(w io.Writer, ts *TableSet)
  Extension() string
}
```

* Fprint: Output table contents. This is main method of this interface.
* Header: Output the header. This method is called only once for each file.
* TableHeader: Output the table headers. This method is called before calling `Fprint` method, and called multiple times if the `multi` flag is `false`.
* TableFooter: Output the table footers. This method is called after calling `Fprint` method, and called multiple times if the `multi` flag is `false`.
* Footer: Output the footer. This method is called only once for each file.
* Extension: Return file extension.

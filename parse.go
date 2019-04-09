package tdconv

import (
	"errors"
	"fmt"

	"github.com/takuoki/clmconv"
	"github.com/takuoki/gsheets"
)

// Parser is a struct to parse the sheet values to the table object.
// Create it using NewParser function.
type Parser struct {

	// table name
	tableNameRow,
	tableNameColumn,

	// columns
	startRow,
	noColumn,
	nameColumn,
	typeColumn,
	pKeyColumn,
	notNullColumn,
	uniqueColumn,
	indexColumn,
	optionColumn,
	commentColumn int

	// other properties
	boolString  string
	keyNameFunc func(string) string

	// non-initialized properties
	commonColumns []Column
}

// NewParser creates a new Parser.
// You can change some parameters of the Parser with ParseOption.
func NewParser(options ...ParseOption) (*Parser, error) {
	p := Parser{
		tableNameRow:    1,
		tableNameColumn: clmconv.MustAtoi("C"),
		startRow:        4,
		noColumn:        clmconv.MustAtoi("B"),
		nameColumn:      clmconv.MustAtoi("C"),
		typeColumn:      clmconv.MustAtoi("D"),
		pKeyColumn:      clmconv.MustAtoi("E"),
		notNullColumn:   clmconv.MustAtoi("F"),
		uniqueColumn:    clmconv.MustAtoi("G"),
		indexColumn:     clmconv.MustAtoi("H"),
		optionColumn:    clmconv.MustAtoi("I"),
		commentColumn:   clmconv.MustAtoi("J"),
		boolString:      "yes",
		keyNameFunc: func(s string) string {
			return s + "_key"
		},
	}
	for _, opt := range options {
		err := opt(&p)
		if err != nil {
			return nil, err
		}
	}
	return &p, nil
}

// ParseOption changes some parameters of the Parser.
type ParseOption func(*Parser) error

// TableNamePos changes the position (row and column) of table name.
func TableNamePos(row int, clm string) ParseOption {
	return func(p *Parser) error {
		if row >= p.startRow {
			return errors.New("Table name row must be smaller than the start row")
		}
		p.tableNameRow = row
		i, err := clmconv.Atoi(clm)
		if err != nil {
			return fmt.Errorf("Unable to convert column string: %v", err)
		}
		p.tableNameColumn = i
		return nil
	}
}

// StartRow changes the start row of column list.
func StartRow(row int) ParseOption {
	return func(p *Parser) error {
		if row <= p.tableNameRow {
			return errors.New("Start row must be greater than the table name row")
		}
		p.startRow = row
		return nil
	}
}

// BoolString changes the bool string in the sheet.
func BoolString(str string) ParseOption {
	return func(p *Parser) error {
		p.boolString = str
		return nil
	}
}

// KeyNameFunc changes the function to convert the column name to the key name.
func KeyNameFunc(f func(string) string) ParseOption {
	return func(p *Parser) error {
		if f == nil {
			return errors.New("Key name function must not be nil")
		}
		p.keyNameFunc = f
		return nil
	}
}

// SetCommonColumns parses the common sheet values and sets them as common columns.
func (p *Parser) SetCommonColumns(s *gsheets.Sheet) error {
	if p == nil {
		return nil
	}
	if len(p.commonColumns) > 0 {
		return errors.New("The common columns are already set")
	}
	t, err := p.parse(s, true)
	if err != nil {
		return err
	}
	p.commonColumns = t.Columns
	return nil
}

// Parse parses the sheet values to the table object.
func (p *Parser) Parse(s *gsheets.Sheet) (*Table, error) {

	if p == nil {
		return nil, nil
	}

	if s.Value(p.tableNameRow, p.tableNameColumn) == "" {
		return nil, errors.New("Table name is required")
	}

	return p.parse(s, false)
}

func (p *Parser) parse(s *gsheets.Sheet, common bool) (*Table, error) {

	t := Table{
		Name:        s.Value(p.tableNameRow, p.tableNameColumn),
		Columns:     make([]Column, 0, 16),
		PKeyColumns: make([]string, 0, 4),
	}

	for i, r := range s.Rows() {

		if i < p.startRow {
			continue
		}
		if r.Value(p.noColumn) == "" {
			break
		}
		if r.Value(p.typeColumn) == "" {
			continue
		}

		if common {
			if r.Value(p.pKeyColumn) == p.boolString {
				return nil, errors.New("The common column must not be PK")
			}
			if r.Value(p.indexColumn) == p.boolString {
				return nil, errors.New("The common column must not have index")
			}
		}

		c := Column{
			Name:     r.Value(p.nameColumn),
			Type:     r.Value(p.typeColumn),
			PKey:     r.Value(p.pKeyColumn) == p.boolString,
			NotNull:  r.Value(p.notNullColumn) == p.boolString,
			Unique:   r.Value(p.uniqueColumn) == p.boolString,
			Index:    r.Value(p.indexColumn) == p.boolString,
			Option:   r.Value(p.optionColumn),
			Comment:  r.Value(p.commentColumn),
			IsCommon: common,
		}
		t.Columns = append(t.Columns, c)

		if c.PKey {
			t.PKeyColumns = append(t.PKeyColumns, c.Name)
		}
		if c.Index {
			t.IndexKeys = append(t.IndexKeys, Key{Name: p.keyNameFunc(c.Name), Columns: []string{c.Name}})
		}
	}

	if len(t.Columns) == 0 {
		return nil, errors.New("The length of table columns must not be zero")
	}

	if len(p.commonColumns) > 0 {
		t.Columns = append(t.Columns, p.commonColumns...)
	}

	return &t, nil
}

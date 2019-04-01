package tdconv

import (
	"errors"

	"github.com/takuoki/clmconv"
	"github.com/takuoki/gsheets"
)

// Parser is a struct to parse the sheet values to the table object.
// Create it using NewParser function.
type Parser struct {
	rowTableName,
	columnTableName,
	rowStart,
	columnNo,
	columnName,
	columnType,
	columnPKey,
	columnNotNull,
	columnUnique,
	columnIndex,
	columnOption,
	columnComment int
	boolString  string
	keyNameFunc func(string) string

	commonColumns []Column
}

// NewParser creates a new Parser.
// You can change some parameters of the Parser with ParseOption.
func NewParser(options ...ParseOption) (*Parser, error) {
	p := Parser{
		rowTableName:    1,
		columnTableName: clmconv.MustAtoi("C"),
		rowStart:        4,
		columnNo:        clmconv.MustAtoi("B"),
		columnName:      clmconv.MustAtoi("C"),
		columnType:      clmconv.MustAtoi("D"),
		columnPKey:      clmconv.MustAtoi("E"),
		columnNotNull:   clmconv.MustAtoi("F"),
		columnUnique:    clmconv.MustAtoi("G"),
		columnIndex:     clmconv.MustAtoi("H"),
		columnOption:    clmconv.MustAtoi("I"),
		columnComment:   clmconv.MustAtoi("J"),
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
	return func(s *Parser) error {
		s.rowTableName = row
		i, err := clmconv.Atoi(clm)
		if err != nil {
			return err
		}
		s.columnTableName = i
		return nil
	}
}

// StartRow changes the start row of column list.
func StartRow(row int) ParseOption {
	return func(s *Parser) error {
		if row < 0 {
			return errors.New("start row must be positive number")
		}
		s.rowStart = row
		return nil
	}
}

// BoolString changes the bool string in the sheet.
func BoolString(str string) ParseOption {
	return func(s *Parser) error {
		s.boolString = str
		return nil
	}
}

// KeyNameFunc changes the function to convert the column name to the key name.
func KeyNameFunc(f func(string) string) ParseOption {
	return func(s *Parser) error {
		s.keyNameFunc = f
		return nil
	}
}

// SetCommonColumns parses the common sheet values and sets them as common columns.
func (p *Parser) SetCommonColumns(s *gsheets.Sheet) error {
	if p != nil && len(p.commonColumns) > 0 {
		return errors.New("the common columns are already set")
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

	if s.Value(p.rowTableName, p.columnTableName) == "" {
		return nil, errors.New("table name is required")
	}

	return p.parse(s, false)
}

func (p *Parser) parse(s *gsheets.Sheet, common bool) (*Table, error) {

	if p == nil {
		return nil, nil
	}

	t := Table{
		Name:        s.Value(p.rowTableName, p.columnTableName),
		Columns:     make([]Column, 0, 16),
		PKeyColumns: make([]string, 0, 4),
	}

	for i, r := range s.Rows() {

		if i < p.rowStart {
			continue
		}
		if r.Value(p.columnNo) == "" {
			break
		}
		if r.Value(p.columnType) == "" {
			continue
		}

		c := Column{
			Name:     r.Value(p.columnName),
			Type:     r.Value(p.columnType),
			PKey:     r.Value(p.columnPKey) == p.boolString,
			NotNull:  r.Value(p.columnNotNull) == p.boolString,
			Unique:   r.Value(p.columnUnique) == p.boolString,
			Index:    r.Value(p.columnIndex) == p.boolString,
			Option:   r.Value(p.columnOption),
			Comment:  r.Value(p.columnComment),
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
		return nil, errors.New("the length of table columns must not be zero")
	}

	if len(p.commonColumns) > 0 {
		t.Columns = append(t.Columns, p.commonColumns...)
	}

	return &t, nil
}

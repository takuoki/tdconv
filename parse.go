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
}

// NewParser creates a new Parser.
// You can change some parameters of the Parser with ParseOption.
func NewParser(options ...ParseOption) (*Parser, error) {
	sp := Parser{
		rowTableName:    1,
		columnTableName: clmconv.MustAtoi("C"),
		rowStart:        5,
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
		err := opt(&sp)
		if err != nil {
			return nil, err
		}
	}
	return &sp, nil
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

// Parse parse the sheet values to the table object.
func (sp *Parser) Parse(s *gsheets.Sheet) (*Table, error) {

	if sp == nil {
		return nil, nil
	}

	tableName := s.Value(sp.rowTableName, sp.columnTableName)
	if tableName == "" {
		return nil, errors.New("table name is required")
	}

	t := Table{
		Name:        tableName,
		Columns:     make([]Column, 0, 16),
		PKeyColumns: make([]string, 0, 4),
	}

	for i, r := range s.Rows() {

		if i < sp.rowStart {
			continue
		}
		if r.Value(sp.columnNo) == "" {
			break
		}
		if r.Value(sp.columnType) == "" {
			continue
		}

		c := Column{
			Name:    r.Value(sp.columnName),
			Type:    r.Value(sp.columnType),
			PKey:    r.Value(sp.columnPKey) == sp.boolString,
			NotNull: r.Value(sp.columnNotNull) == sp.boolString,
			Unique:  r.Value(sp.columnUnique) == sp.boolString,
			Index:   r.Value(sp.columnIndex) == sp.boolString,
			Option:  r.Value(sp.columnOption),
			Comment: r.Value(sp.columnComment),
		}
		t.Columns = append(t.Columns, c)
		if c.PKey {
			t.PKeyColumns = append(t.PKeyColumns, c.Name)
		}
		if c.Index {
			t.IndexKeys = append(t.IndexKeys, Key{Name: sp.keyNameFunc(c.Name), Columns: []string{c.Name}})
		}
	}

	if len(t.Columns) == 0 {
		return nil, errors.New("the length of table columns must not be zero")
	}

	return &t, nil
}

// ParseColumns parse the sheet values to the column object list.
func (sp *Parser) ParseColumns(s *gsheets.Sheet) ([]Column, error) {
	t, err := sp.Parse(s)
	if err != nil {
		return nil, err
	}
	return t.Columns, nil
}

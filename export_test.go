package tdconv

func (p *Parser) TableNameRow() int {
	if p == nil {
		return 0
	}
	return p.tableNameRow
}

func (p *Parser) TableNameColumn() int {
	if p == nil {
		return 0
	}
	return p.tableNameColumn
}

func (p *Parser) StartRow() int {
	if p == nil {
		return 0
	}
	return p.startRow
}

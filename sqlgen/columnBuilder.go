package sqlgen

import "fmt"

// VarChar sets VarChar(given size) type
// to the column
func (cb *ColumnBuilder) VarChar(size int64) *CreateTableBuilder {
	return cb.tableBuilder.appendColumn(&tableColumn{
		name:    cb.name,
		sqlType: fmt.Sprintf("VARCHAR(%v)", size),
	})
}

// Integer sets INTEGER type to the column
func (cb *ColumnBuilder) Integer() *CreateTableBuilder {
	return cb.tableBuilder.appendColumn(&tableColumn{
		name:    cb.name,
		sqlType: "INTEGER",
	})
}

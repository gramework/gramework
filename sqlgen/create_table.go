package sqlgen

import (
	"fmt"
	"strings"
)

// Column initialize a column builder
// that requires you to choose a type
// before you can Build the statement
func (tb *CreateTableBuilder) Column(name string) *ColumnBuilder {
	return &ColumnBuilder{
		tableBuilder: tb,
		name:         name,
	}
}

// appendColumn to the table column list
func (tb *CreateTableBuilder) appendColumn(c *tableColumn) *CreateTableBuilder {
	tb.columns = append(tb.columns, *c)
	return tb
}

// Build the statement
func (tb *CreateTableBuilder) Build() string {
	return fmt.Sprintf("CREATE TABLE %s(\n    %s);\n", tb.name, tb.buildColumns())
}

func (tb *CreateTableBuilder) buildColumns() string {
	queryParts := []string{}
	for _, c := range tb.columns {
		queryParts = append(queryParts, fmt.Sprintf("%s %s", c.name, c.sqlType))
	}

	return strings.Join(queryParts, ",\n\t")
}

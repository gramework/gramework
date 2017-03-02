package sqlgen

import (
	"fmt"
	"strings"
	"sync"
)

// Insert statement builder generates
// an insert statement using `?` placeholders
// for values
func Insert(table string) *InsertBuilder {
	return &InsertBuilder{
		tableName: table,
		query:     fmt.Sprintf(`INSERT INTO %s `, table),
		lock:      &sync.Mutex{},
	}
}

// PreparedInsert statement builder generates
// the insert SQL statement with values built
// in the statement without using placeholders
func PreparedInsert(table string) *InsertBuilder {
	i := Insert(table)
	i.prepared = true
	return i
}

// Columns defines column list
func (b *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	b.lock.Lock()
	b.columns = columns
	b.query = fmt.Sprintf(`%s(`, b.query)
	for i, column := range columns {
		b.query = fmt.Sprintf(`%s%s`, b.query, column)
		if i < len(columns)-1 {
			b.query = fmt.Sprintf(`%s,`, b.query)
		}
	}
	b.query = fmt.Sprintf(`%s)`, b.query)
	b.lock.Unlock()
	return b
}

// Values appends column values to the query
func (b *InsertBuilder) Values(columnValues ...interface{}) *InsertBuilder {
	b.lock.Lock()
	sqlValue := "("
	if b.prepared {
		for k, columnValue := range columnValues {
			switch v := columnValue.(type) {
			case string:
				sqlValue = fmt.Sprintf("%s'%s'", sqlValue, strings.Replace(v, "'", "''", -1))
			default:
				sqlValue = fmt.Sprintf("%s%v", sqlValue, v)
			}
			if k < len(columnValues)-1 {
				sqlValue = fmt.Sprintf("%s, ", sqlValue)
			}
		}
	} else {
		for k := range columnValues {
			sqlValue = fmt.Sprintf("%s?", sqlValue)
			if k < len(columnValues)-1 {
				sqlValue = fmt.Sprintf("%s, ", sqlValue)
			}
		}
	}
	b.sqlValues = append(b.sqlValues, fmt.Sprintf(`%s)`, sqlValue))

	b.lock.Unlock()

	return b
}

// Build the query
func (b *InsertBuilder) Build() string {
	b.lock.Lock()
	defer b.lock.Unlock()
	return fmt.Sprintf("%s\n    VALUES %s;", b.query, strings.Join(b.sqlValues, ", \n        "))
}

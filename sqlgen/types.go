package sqlgen

import "sync"

// InsertBuilder defines internal implementation
// of a INSERT statement builder
type InsertBuilder struct {
	tableName string
	query     string
	columns   []string
	sqlValues []string
	prepared  bool
	lock      *sync.Mutex
}

// CreateBuilder stands for the API
// of CREATE statement builder
type CreateBuilder struct{}

// CreateDatabaseBuilder handles internal
// info about CREATE DATABASE statement
// that now builds
type CreateDatabaseBuilder struct {
	name  string
	useIt bool
}

// CreateTableBuilder handles internal
// info about create table statement
// that now builds
type CreateTableBuilder struct {
	name    string
	columns []tableColumn
}

// ColumnBuilder handles internal
// column info
type ColumnBuilder struct {
	name         string
	sqlType      string
	tableBuilder *CreateTableBuilder
}

type tableColumn struct {
	name    string
	sqlType string
}

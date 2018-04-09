package sqlgen

import (
	"sync"
)

type (

	// InsertBuilder defines internal implementation
	// of a INSERT statement builder
	InsertBuilder struct {
		tableName string
		query     string
		columns   []string
		sqlValues []string
		prepared  bool
		lock      *sync.Mutex
	}

	// CreateBuilder stands for the API
	// of CREATE statement builder
	CreateBuilder struct{}

	// CreateDatabaseBuilder handles internal
	// info about CREATE DATABASE statement
	// that now builds
	CreateDatabaseBuilder struct {
		name  string
		useIt bool
	}

	// CreateTableBuilder handles internal
	// info about create table statement
	// that now builds
	CreateTableBuilder struct {
		name    string
		columns []tableColumn
	}

	// ColumnBuilder handles internal
	// column info
	ColumnBuilder struct {
		name         string
		sqlType      string
		tableBuilder *CreateTableBuilder
	}

	tableColumn struct {
		name    string
		sqlType string
	}
)

package sqlgen

// Create starts construction
// of CREATE query
func Create() *CreateBuilder {
	return &CreateBuilder{}
}

// Database says that we are CREATE-ing
// a DATABASE with given name
// It's an alias of (*CreateBuilder).DB
func (cb *CreateBuilder) Database(name string) *CreateDatabaseBuilder {
	return cb.DB(name)
}

// DB says that we are CREATE-ing
// a DATABASE with given name
func (cb *CreateBuilder) DB(name string) *CreateDatabaseBuilder {
	return &CreateDatabaseBuilder{
		name: name,
	}
}

// Table says that we are CREATE-ing
// a TABLE with given name
func (cb *CreateBuilder) Table(name string) *CreateTableBuilder {
	return &CreateTableBuilder{
		name:    name,
		columns: make([]tableColumn, 0),
	}
}

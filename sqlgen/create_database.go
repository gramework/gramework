package sqlgen

import "fmt"

// UseIt says that we need to USE
// the database we just created
func (cb *CreateDatabaseBuilder) UseIt() *CreateDatabaseBuilder {
	cb.useIt = true
	return cb
}

// Build the statement
func (cb *CreateDatabaseBuilder) Build() string {
	q := fmt.Sprintf("CREATE DATABASE %s;", cb.name)
	if cb.useIt {
		q = fmt.Sprintf("%s\nUSE %s;\n", q, cb.name)
	}
	return q
}

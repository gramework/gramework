package healthchecks

// Register both ping and healthcheck endpoints
func Register(r interface{}, collectors ...func() (statKey string, stats interface{})) error {
	return doReg(r, collectors, true, true)
}

// RegisterPing registers ping endpoint
func RegisterPing(r interface{}) error {
	return doReg(r, nil, true, false)
}

// RegisterHealthcheck registers healthcheck endpoint
func RegisterHealthcheck(r interface{}, collectors ...func() (statKey string, stats interface{})) error {
	return doReg(r, collectors, false, true)
}

// ServeHealthcheck serves healthcheck
func ServeHealthcheck(collectors ...func() (statKey string, stats interface{})) func() interface{} {
	return func() interface{} {
		return check(collectors...)
	}
}

package gqlhandler

import (
	"strings"

	"github.com/gramework/gramework"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
)

var (
	// ErrNilSchema returned when New() called with nil argument
	ErrNilSchema = errors.New("could not initialize gqlhandler state: schema is nil")
)

// State is the GQL handler state
type State struct {
	NoIntrospection bool
	schema          *graphql.Schema
}

// New returns gql handler state based on given schema
func New(schema *graphql.Schema) (*State, error) {
	if schema == nil {
		return nil, ErrNilSchema
	}

	s := &State{
		schema: schema,
	}
	return s, nil
}

// Handler is the simplest
func (s *State) Handler(ctx *gramework.Context) {
	if s.schema == nil {
		ctx.Logger.Error("schema is nil")
		ctx.Err500()
		return
	}

	// ok, we have the schema. try to decode request
	req, err := ctx.DecodeGQL()
	if err != nil {
		ctx.Logger.Warn("gql request decoding failed")
		ctx.Error("Invalid request", 400)
		return
	}

	if s.NoIntrospection && strings.HasPrefix(strings.TrimSpace(req.Query), "query IntrospectionQuery") {
		ctx.Forbidden()
		return
	}

	// check if we got invalid content type
	if req == nil {
		ctx.Logger.Error("GQL request is nil: invalid content type")
		ctx.Error("Invalid content type", 400)
		return
	}

	ctx.Logger.Info("processing request")
	if _, err := ctx.Encode(s.schema.Exec(ctx.ToContext(), req.Query, req.OperationName, req.Variables)); err != nil {
		ctx.SetStatusCode(415)
	}
	ctx.Logger.Info("processing request done")
}

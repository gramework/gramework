package gramework

import (
	"strings"
)

// DecodeGQL parses GraphQL request and returns data from it
func (ctx *Context) DecodeGQL() (*GQLRequest, error) {
	r := &GQLRequest{}

	if string(ctx.Method()) == GET {
		query := ctx.GETParam("query")
		if len(query) == 0 {
			return nil, ErrInvalidGQLRequest
		}
		r.Query = query[0]

		if operationName := ctx.GETParam("operationName"); len(operationName) != 0 {
			r.OperationName = operationName[0]
		}

		if variables := ctx.GETParam("variables"); len(variables) != 0 {
			if _, err := ctx.UnJSONBytes([]byte(variables[0]), &r.Variables); err != nil {
				return nil, ErrInvalidGQLRequest
			}
		}

		return r, nil
	}

	if strings.HasPrefix(ctx.ContentType(), jsonCTshort) {
		if err := ctx.UnJSON(&r); err != nil {
			return nil, err
		}
	}

	switch ctx.ContentType() {
	case gqlCT:
		r.Query = string(ctx.PostBody())
	}

	return r, nil
}

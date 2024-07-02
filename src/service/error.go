package service

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AddError(ctx context.Context, msg, err, code string) {
	graphql.AddError(ctx, &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: msg,
		Extensions: map[string]interface{}{
			"code": code,
			"err":  err,
		},
	})
}

package middleware

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/alphaflow/api-core/buffalo/trace"
)

type UpFunction func(connection *gorm.DB) bool

// ServiceCheck returns a middleware that checks if a required service is up, given an UpFunction and a DB Connection.
func ServiceCheck(upFunc UpFunction, connection *gorm.DB, errorMessage string) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			if !upFunc(connection) {
				return trace.WithStatus(c, http.StatusServiceUnavailable, errors.New(errorMessage))
			}

			return next(c)
		}
	}
}

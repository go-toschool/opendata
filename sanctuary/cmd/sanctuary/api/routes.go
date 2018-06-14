package api

import (
	"net/http"

	router "github.com/Finciero/httprouter"
)

// Routes ...
func Routes(ctx *Context) http.Handler {
	r := router.New()

	r.POST("/extract", ctx.Handle(createExtraction))
	return r
}

package auth

import (
	"net/http"

	router "github.com/Finciero/httprouter"
)

// Routes ...
func Routes(ctx *Context) http.Handler {
	r := router.New()

	r.POST("/session", ctx.Handle(createSession))

	return r
}

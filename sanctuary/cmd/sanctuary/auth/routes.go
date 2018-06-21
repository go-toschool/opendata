package auth

import (
	"net/http"

	router "github.com/Finciero/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// Routes ...
func Routes(ctx *Context) http.Handler {
	r := router.New()

	r.POST("/session", ctx.Handle(createSession))

	tokenAuth := NewMiddleware(ctx.ShuraClient)
	routes := negroni.Wrap(r)
	cors := cors.New(cors.Options{
		AllowedOrigins:     ctx.AllowedOrigins,
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowedMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
	})

	return negroni.New(tokenAuth, cors, routes)
}

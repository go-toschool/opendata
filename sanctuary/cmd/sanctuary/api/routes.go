package api

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/urfave/negroni"

	router "github.com/go-toschool/httprouter"
)

// Routes ...
func Routes(ctx *Context) http.Handler {
	r := router.New()

	r.POST("/extract", ctx.Handle(createExtraction))

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

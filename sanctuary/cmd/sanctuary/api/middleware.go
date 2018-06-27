package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Finciero/opendata/capricornius/shura"
	"google.golang.org/grpc/metadata"
)

const (
	tokenTypePrefix = "Bearer "
	tokenHeaderKey  = "X-FIN-CLIENT-TOKEN"
)

// Middleware provides a middleware to authenticate an incoming request.
type Middleware struct {
	ShuraClient shura.ServiceClient
}

// NewMiddleware creates a new AuthMiddleware with the given user session service.
func NewMiddleware(ssc shura.ServiceClient) *Middleware {
	return &Middleware{ShuraClient: ssc}
}

// Handle authenticate the incoming request, if the authentication process fails then an
// ErrUnauthorizedAccess is returned.
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			// router handles the OPTIONS request to obtain the list of allowed methods.
			next.ServeHTTP(w, r)
			return
		}

		token, err := parseAdminAuthToken(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("unauthorized access: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		meta := map[string]string{
			"service":  "sanctuary",
			"endpoint": "api",
			"version":  "v0.0.1",
			"time":     time.Now().String(),
		}
		md := metadata.New(meta)
		ctx = metadata.NewIncomingContext(ctx, md)
		rr, err := m.ShuraClient.GetToken(ctx, &shura.Token{
			PartnerToken: token,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !rr.Valid {
			http.Error(w, "unauthorized access: invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ServeHTTP implements a negroni compatible signature.
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m.Handle(next).ServeHTTP(w, r)
}

func parseAdminAuthToken(r *http.Request) (string, error) {
	header := r.Header.Get(tokenHeaderKey)
	if !strings.HasPrefix(header, tokenTypePrefix) {
		return "", errors.New("Invalid token format")
	}
	return header[len(tokenTypePrefix):], nil
}

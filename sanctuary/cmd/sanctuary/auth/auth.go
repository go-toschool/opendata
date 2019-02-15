package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-toschool/opendata/taurus/aldebaran"
	"google.golang.org/grpc/metadata"
)

type userData struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type session struct {
	Token string `json:"token,omitempty"`
}

// Createsession creates a new user token on Aldebaran service.
// This method return a user key, this user key is approved by the user
// that belong the new token.
func createSession(ctx *Context, w http.ResponseWriter, r *http.Request) (*Response, error) {
	var user *userData

	if err := ctx.ToJSON(r.Body, user); err != nil {
		return nil, err
	}

	c := newContext()
	rr, err := ctx.AldebaranClient.CreateToken(c, &aldebaran.Create{
		Email:       user.Email,
		ClientToken: ctx.clientToken,
	})
	if err != nil {
		return nil, err
	}

	if rr.StatusCode == 201 {
		return &Response{
			Data: &session{
				Token: rr.Token,
			},
		}, nil
	}

	return nil, errors.New("something went wrong")
}

func newContext() context.Context {
	ctx := context.Background()
	meta := map[string]string{
		"service":  "sanctuary",
		"endpoint": "auth/createSession",
		"version":  "v0.0.1",
		"time":     time.Now().String(),
	}
	md := metadata.New(meta)
	return metadata.NewIncomingContext(ctx, md)
}

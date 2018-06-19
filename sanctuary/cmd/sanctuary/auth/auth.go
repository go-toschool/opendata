package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/Finciero/opendata/taurus/aldebaran"
)

type userData struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type session struct {
	Token string `json:"token,omitempty"`
}

func createSession(ctx *Context, w http.ResponseWriter, r *http.Request) (*Response, error) {
	clientToken, ok := r.Context().Value("client_token").(string)
	if !ok {
		return nil, errors.New("invalid client token")
	}

	var user *userData

	if err := ctx.ToJSON(r.Body, user); err != nil {
		return nil, err
	}

	c := context.Background()
	rr, err := ctx.TaurusClient.CreateToken(c, &aldebaran.Create{
		Email:       user.Email,
		ClientToken: clientToken,
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

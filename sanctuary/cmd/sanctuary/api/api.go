package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Finciero/opendata/aries/mu"
	"github.com/Finciero/opendata/sanctuary"
)

func createExtraction(ctx *Context, w http.ResponseWriter, r *http.Request) (*Response, error) {
	// validate token
	// request aries
	var request sanctuary.ExtractionRequest

	if err := ctx.ToJSON(r.Body, request); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c := context.Background()
	res, err := ctx.MuClient.Extract(c, &mu.Request{
		PartnerToken: ctx.clientToken,
		UserToken:    request.UserToken,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(res.Balance)
	fmt.Println(res.Account)
	fmt.Println(res.Message)
	fmt.Println(res.StatusCode)

	return &Response{}, nil
}

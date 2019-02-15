package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-toschool/opendata/aries/mu"
	"github.com/go-toschool/opendata/sanctuary"
	"google.golang.org/grpc/metadata"
)

func createExtraction(ctx *Context, w http.ResponseWriter, r *http.Request) (*Response, error) {
	// validate token
	// request aries
	var request sanctuary.ExtractionRequest

	if err := ctx.ToJSON(r.Body, request); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c := newContext()
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

func newContext() context.Context {
	ctx := context.Background()
	meta := map[string]string{
		"service":  "sanctuary",
		"endpoint": "api/createExtraction",
		"version":  "v0.0.1",
		"time":     time.Now().String(),
	}
	md := metadata.New(meta)
	return metadata.NewIncomingContext(ctx, md)
}

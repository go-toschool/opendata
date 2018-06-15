package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/Finciero/opendata/gemini"
	"github.com/Finciero/opendata/gemini/castor"
)

func main() {
	gc := gemini.NewCastor(&gemini.Config{
		Host: "localhost",
		Port: 4000,
		Cert: "certs/castor/castor.crt",
	})

	ctx := context.Background()
	m := map[string]string{
		"X-Client-Id": "rodrwan",
	}
	md := metadata.New(m)
	ctx = metadata.NewIncomingContext(ctx, md)

	rr, err := gc.Card(ctx, &castor.Request{
		ClientId:    "Rodrwan",
		UserId:      "user-id-241243234",
		ReferenceId: "21423424",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(rr.StatusCode)
	fmt.Println(rr.Balance)
}

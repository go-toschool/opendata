package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Finciero/opendata/sagittarius/aiolos"
	"google.golang.org/grpc"

	"github.com/Finciero/opendata/gemini/brickwall"
	"github.com/Finciero/opendata/gemini/kanon"
)

// KanonService implements kanon interface of gemini package
type KanonService struct {
	BrickwallClient   *brickwall.Client
	SagittariusClient aiolos.SaggitariusClient
}

// GetTransactions ...
func (ks *KanonService) GetTransactions(ctx context.Context, r *kanon.Request) (*kanon.Response, error) {
	path := fmt.Sprintf("/cards/%s/transactions", r.ReferenceId)
	resp, err := ks.BrickwallClient.Get(path)
	if err != nil {
		return nil, err
	}

	//var trxs kanon.Transactions
	fmt.Println(string(resp))

	rr, err := ks.SagittariusClient.Dispatch(ctx, &aiolos.Request{
		Id: r.ReferenceId,
	})
	if err != nil {
		return nil, err
	}

	return &kanon.Response{
		StatusCode: 200,
		Message:    rr.Message,
	}, nil
}

func main() {
	srv := grpc.NewServer()

	brickwallClient := brickwall.NewClient(&brickwall.Config{
		Token: "",
	})

	// Dial with saggitarius
	conn, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	sc := aiolos.NewSaggitariusClient(conn)

	ks := &KanonService{
		BrickwallClient:   brickwallClient,
		SagittariusClient: sc,
	}

	kanon.RegisterServiceServer(srv, ks)

	lis, err := net.Listen("tcp", ":4001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("starting kanon service...")
	fmt.Println("listening on: 4001")
	srv.Serve(lis)
}

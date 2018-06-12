package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/Finciero/opendata/gemini/brickwall"
	"github.com/Finciero/opendata/gemini/saga"
)

// SagaService implements gemini interface of gemini package
type SagaService struct {
	BrickwallClient *brickwall.Client
}

// GetBalance ...
func (ss *SagaService) GetBalance(ctx context.Context, r *saga.Request) (*saga.Response, error) {
	resp, err := ss.BrickwallClient.Get(fmt.Sprintf("/cards/%s/balance", r.ReferenceId))
	if err != nil {
		return nil, err
	}

	var balance *saga.Balance
	if err := json.Unmarshal(resp, &balance); err != nil {
		return nil, err
	}

	return &saga.Response{
		StatusCode: 200,
		Balance: &saga.Balance{
			Email:   "",
			Balance: 0,
			UserId:  "",
		},
	}, nil
}

func main() {
	srv := grpc.NewServer()

	brickwallClient := brickwall.NewClient(&brickwall.Config{
		Token: "",
	})

	gs := &SagaService{
		BrickwallClient: brickwallClient,
	}

	saga.RegisterServiceServer(srv, gs)

	lis, err := net.Listen("tcp", ":4002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("starting gemini service...")
	fmt.Println("listening on: 4002")
	srv.Serve(lis)
}

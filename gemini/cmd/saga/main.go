package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Finciero/opendata/gemini/brickwall"
	"github.com/Finciero/opendata/gemini/saga"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if CASTOR_SERVICE_HOST env var is set)")
	port = flag.Int("port", 4002, "Service port (Overwriten if SAGA_SERVICE_PORT env var is set)")

	sagaCert = flag.String("saga-cert", "", "Saga cert (Overwriten if SAGA_CERT env var is set)")
	sagaKey  = flag.String("saga-key", "", "Saga key (Overwriten if SAGA_KEY env var is set)")
)

func main() {
	flag.Parse()

	brickwallClient := brickwall.NewClient(&brickwall.Config{
		Token: "",
	})

	gs := &SagaService{
		BrickwallClient: brickwallClient,
	}

	creds, err := credentials.NewServerTLSFromFile(*sagaCert, *sagaKey)
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	srv := grpc.NewServer(opts...)
	saga.RegisterServiceServer(srv, gs)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("starting saga service...")
	log.Println("listening on: 4002")
	srv.Serve(lis)
}

// SagaService implements gemini interface of gemini package
type SagaService struct {
	BrickwallClient *brickwall.Client
}

// GetBalance ...
func (ss *SagaService) GetBalance(ctx context.Context, r *saga.Request) (*saga.Response, error) {
	// resp, err := ss.BrickwallClient.Get(fmt.Sprintf("/cards/%s/balance", r.ReferenceId))
	// if err != nil {
	// 	return nil, err
	// }

	// var balance *saga.Balance
	// if err := json.Unmarshal(resp, &balance); err != nil {
	// 	return nil, err
	// }

	return &saga.Response{
		StatusCode: 200,
		Balance: &saga.Balance{
			Email:   "",
			Balance: 0,
			UserId:  "",
		},
	}, nil
}

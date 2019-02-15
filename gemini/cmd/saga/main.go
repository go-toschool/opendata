package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/go-toschool/opendata/gemini/brickwall"
	"github.com/go-toschool/opendata/gemini/saga"
)

var (
	host = flag.String("host", "saga", "Service host (Overwriten if CASTOR_SERVICE_HOST env var is set)")
	port = flag.Int("port", 4002, "Service port (Overwriten if SAGA_SERVICE_PORT env var is set)")

	sagaCert = flag.String("saga-cert", "", "Saga cert (Overwriten if SAGA_CERT env var is set)")
	sagaKey  = flag.String("saga-key", "", "Saga key (Overwriten if SAGA_KEY env var is set)")

	withTLS = flag.Bool("with-tls", false, "service with TLS")
)

func main() {
	flag.Parse()

	brickwallClient := brickwall.NewClient(&brickwall.Config{
		Token: "",
	})

	gs := &SagaService{
		BrickwallClient: brickwallClient,
	}

	var srv *grpc.Server
	if *withTLS {
		creds, err := credentials.NewServerTLSFromFile(*sagaCert, *sagaKey)
		if err != nil {
			log.Fatalf("could not load TLS keys: %s", err)
		}
		opts := []grpc.ServerOption{grpc.Creds(creds)}
		srv = grpc.NewServer(opts...)
	} else {
		srv = grpc.NewServer()
	}

	saga.RegisterServiceServer(srv, gs)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Starting saga service...")
	log.Println(fmt.Sprintf("Listening on: %d", *port))
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// SagaService implements gemini interface of gemini package
type SagaService struct {
	BrickwallClient *brickwall.Client
}

// GetBalance ...
func (ss *SagaService) GetBalance(ctx context.Context, r *saga.Request) (*saga.Response, error) {
	log.Println("Calling GetBalance...")

	log.Printf("Account ID: %s\n", r.AccountId)
	log.Printf("Reference ID: %s", r.ReferenceId)
	log.Printf("User ID: %s\n", r.UserId)
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
			AccountId: r.AccountId,
			UserId:    r.UserId,
			Balance:   float32(rand.Int31n(2000)),
		},
	}, nil
}

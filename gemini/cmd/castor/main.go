package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/Finciero/opendata/gemini"
	"github.com/Finciero/opendata/gemini/castor"
	"github.com/Finciero/opendata/gemini/kanon"
	"github.com/Finciero/opendata/gemini/saga"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if CASTOR_SERVICE_HOST env var is set)")
	port = flag.Int("port", 4000, "Service port (Overwriten if CASTOR_SERVICE_PORT env var is set)")

	kanonHost = flag.String("kanon-host", "kanon", "Service host (Overwriten if KANON_SERVICE_HOST env var is set)")
	kanonPort = flag.Int("kanon-port", 4001, "Service port (Overwriten if KANON_SERVICE_PORT env var is set)")

	sagaHost = flag.String("saga-host", "saga", "Service host (Overwriten if SAGA_SERVICE_HOST env var is set)")
	sagaPort = flag.Int("saga-port", 4002, "Service port (Overwriten if SAGA_SERVICE_PORT env var is set)")
)

func main() {
	flag.Parse()
	srv := grpc.NewServer()

	// Dial with Kanon
	ks := gemini.NewKanon(&gemini.Config{
		Host: *kanonHost,
		Port: *kanonPort,
	})

	// Dial with Saga
	ss := gemini.NewSaga(&gemini.Config{
		Host: *sagaHost,
		Port: *sagaPort,
	})

	gs := &CastorService{
		KanonClient: ks,
		SagaClient:  ss,
	}

	castor.RegisterServiceServer(srv, gs)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("starting Castor service...")
	log.Println(fmt.Sprintf("listening on: %s:%d", *host, *port))
	srv.Serve(lis)
}

// CastorService implements gemini interface of gemini package
type CastorService struct {
	KanonClient kanon.ServiceClient
	SagaClient  saga.ServiceClient
}

// Card ...
func (ss *CastorService) Card(ctx context.Context, r *castor.Request) (*castor.Response, error) {
	go ss.KanonClient.GetTransactions(ctx, &kanon.Request{
		ReferenceId: r.ReferenceId,
		UserId:      r.UserId,
	})

	balance, err := ss.SagaClient.GetBalance(ctx, &saga.Request{
		ReferenceId: r.ReferenceId,
		UserId:      r.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &castor.Response{
		StatusCode: balance.StatusCode,
		Balance: &saga.Balance{
			Email:   balance.Balance.Email,
			Balance: balance.Balance.Balance,
			UserId:  balance.Balance.UserId,
		},
	}, nil
}

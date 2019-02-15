package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/metadata"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/go-toschool/opendata/gemini"
	"github.com/go-toschool/opendata/gemini/castor"
	"github.com/go-toschool/opendata/gemini/kanon"
	"github.com/go-toschool/opendata/gemini/saga"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if CASTOR_SERVICE_HOST env var is set)")
	port = flag.Int("port", 4000, "Service port (Overwriten if CASTOR_SERVICE_PORT env var is set)")

	kanonHost = flag.String("kanon-host", "localhost", "Service host (Overwriten if KANON_SERVICE_HOST env var is set)")
	kanonPort = flag.Int("kanon-port", 4001, "Service port (Overwriten if KANON_SERVICE_PORT env var is set)")

	sagaHost = flag.String("saga-host", "localhost", "Service host (Overwriten if SAGA_SERVICE_HOST env var is set)")
	sagaPort = flag.Int("saga-port", 4002, "Service port (Overwriten if SAGA_SERVICE_PORT env var is set)")

	castorCert = flag.String("castor-cert", "", "Castor cert (Overwriten if CASTOR_CERT env var is set)")
	castorKey  = flag.String("castor-key", "", "Castor key (Overwriten if CASTOR_KEY env var is set)")

	kanonCert = flag.String("kanon-cert", "", "Kanon cert (Overwriten if KANON_CERT env var is set)")
	sagaCert  = flag.String("saga-cert", "", "Saga cert (Overwriten if SAGA_CERT env var is set)")

	withTLS = flag.Bool("with-tls", false, "service with TLS")
)

func main() {
	flag.Parse()
	// Dial with Kanon
	ks := gemini.NewKanon(&gemini.Config{
		Host: *kanonHost,
		Port: *kanonPort,
		Cert: *kanonCert,
	})

	// Dial with Saga
	ss := gemini.NewSaga(&gemini.Config{
		Host: *sagaHost,
		Port: *sagaPort,
		Cert: *sagaCert,
	})

	gs := &CastorService{
		KanonClient: ks,
		SagaClient:  ss,
	}

	var srv *grpc.Server
	if *withTLS {
		creds, err := credentials.NewServerTLSFromFile(*castorCert, *castorKey)
		if err != nil {
			log.Fatalf("could not load TLS keys: %s", err)
		}
		opts := []grpc.ServerOption{grpc.Creds(creds)}
		srv = grpc.NewServer(opts...)
	} else {
		srv = grpc.NewServer()
	}

	castor.RegisterServiceServer(srv, gs)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Starting Castor service...")
	log.Println(fmt.Sprintf("Listening on: %d", *port))
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// CastorService implements gemini interface of gemini package
type CastorService struct {
	KanonClient kanon.ServiceClient
	SagaClient  saga.ServiceClient
}

// Card ...
func (ss *CastorService) Card(ctx context.Context, r *castor.Request) (*castor.Response, error) {
	log.Println("Calling Card...")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("Bad metadata format")
	}

	fmt.Printf("%#v\n", md)
	fmt.Printf("X-FIN-CLIENT-ID: %s\n", md.Get("X-FIN-CLIENT-ID"))

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
			AccountId: balance.Balance.AccountId,
			Balance:   balance.Balance.Balance,
			UserId:    balance.Balance.UserId,
		},
	}, nil
}

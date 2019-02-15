package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/go-toschool/opendata/sagittarius"
	"github.com/go-toschool/opendata/sagittarius/aiolos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/go-toschool/opendata/gemini/brickwall"
	"github.com/go-toschool/opendata/gemini/kanon"
)

var (
	host = flag.String("host", "kanon", "Service host (Overwriten if KANON_SERVICE_HOST env var is set)")
	port = flag.Int("port", 4001, "Service port (Overwriten if KANON_SERVICE_PORT env var is set)")

	aiolosHost = flag.String("aiolos-host", "aiolos", "Service host (Overwriten if AIOLOS_SERVICE_HOST env var is set)")
	aiolosPort = flag.Int("aiolos-port", 3000, "Service port (Overwriten if AIOLOS_SERVICE_PORT env var is set)")

	kanonCert = flag.String("kanon-cert", "", "Kanon cert (Overwriten if KANON_CERT env var is set)")
	kanonKey  = flag.String("kanon-key", "", "Kanon key (Overwriten if KANON_KEY env var is set)")

	brickwallToken = flag.String("brickwall-token", "", "Token to access brickwall service.")

	withTLS = flag.Bool("with-tls", false, "service with TLS")
)

func main() {
	flag.Parse()

	brickwallClient := brickwall.NewClient(&brickwall.Config{
		Token: *brickwallToken,
	})

	// Dial with saggitarius
	sc := sagittarius.NewAiolos(&sagittarius.Config{
		Host: *aiolosHost,
		Port: *aiolosPort,
	})

	ks := &KanonService{
		BrickwallClient:   brickwallClient,
		SagittariusClient: sc,
	}

	var srv *grpc.Server
	if *withTLS {
		creds, err := credentials.NewServerTLSFromFile(*kanonCert, *kanonKey)
		if err != nil {
			log.Fatalf("could not load TLS keys: %s", err)
		}
		opts := []grpc.ServerOption{grpc.Creds(creds)}
		srv = grpc.NewServer(opts...)
	} else {
		srv = grpc.NewServer()
	}

	kanon.RegisterServiceServer(srv, ks)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Starting kanon service...")
	log.Println(fmt.Sprintf("Listening on: %d", *port))
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// KanonService implements kanon interface of gemini package
type KanonService struct {
	BrickwallClient   *brickwall.Client
	SagittariusClient aiolos.ServiceClient
}

// GetTransactions ...
func (ks *KanonService) GetTransactions(ctx context.Context, r *kanon.Request) (*kanon.Response, error) {
	log.Println("Calling GetTransactions...")

	log.Printf("Account ID: %s\n", r.AccountId)
	log.Printf("Reference ID: %s", r.ReferenceId)
	log.Printf("User ID: %s\n", r.UserId)
	// path := fmt.Sprintf("/cards/%s/transactions", r.ReferenceId)
	// resp, err := ks.BrickwallClient.Get(path)
	// if err != nil {
	// 	return nil, err
	// }

	// //var trxs kanon.Transactions
	// log.Println(string(resp))

	// rr, err := ks.SagittariusClient.Dispatch(ctx, &aiolos.Request{
	// 	Id: r.ReferenceId,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	return &kanon.Response{
		StatusCode: 200,
		Message:    "GetTransactions", //rr.Message,
	}, nil
}

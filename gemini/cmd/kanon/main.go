package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/Finciero/opendata/sagittarius/aiolos"
	"google.golang.org/grpc"

	"github.com/Finciero/opendata/gemini/brickwall"
	"github.com/Finciero/opendata/gemini/kanon"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if KANON_SERVICE_HOST env var is set)")
	port = flag.Int("port", 4001, "Service port (Overwriten if KANON_SERVICE_PORT env var is set)")

	aiolosHost = flag.String("aiolos-host", "aiolos", "Service host (Overwriten if AIOLOS_SERVICE_HOST env var is set)")
	aiolosPort = flag.Int("aiolos-port", 3000, "Service port (Overwriten if AIOLOS_SERVICE_PORT env var is set)")

	brickwallToken = flag.String("brickwall-token", "", "Token to access brickwall service.")
)

func main() {
	flag.Parse()
	srv := grpc.NewServer()

	brickwallClient := brickwall.NewClient(&brickwall.Config{
		Token: *brickwallToken,
	})

	// Dial with saggitarius
	aiolosIP := fmt.Sprintf("%s:%d", *aiolosHost, *aiolosPort)
	conn, err := grpc.Dial(aiolosIP, grpc.WithInsecure())
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("starting kanon service...")
	log.Println(fmt.Sprintf("listening on: %s:%d", *host, *port))
	srv.Serve(lis)
}

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
	log.Println(string(resp))

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

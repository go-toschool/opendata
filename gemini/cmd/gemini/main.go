package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/Finciero/opendata/gemini/gemini"
	"github.com/Finciero/opendata/gemini/kanon"
	"github.com/Finciero/opendata/gemini/saga"
	"github.com/Finciero/opendata/gemini/sigiriya"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if GEMINI_SERVICE_HOST env var is set)")
	port = flag.Int("port", 4000, "Service port (Overwriten if GEMINI_SERVICE_PORT env var is set)")

	kanonHost = flag.String("kanon-host", "kanon", "Service host (Overwriten if KANON_SERVICE_HOST env var is set)")
	kanonPort = flag.Int("kanon-port", 4001, "Service port (Overwriten if KANON_SERVICE_PORT env var is set)")

	sagaHost = flag.String("saga-host", "saga", "Service host (Overwriten if SAGA_SERVICE_HOST env var is set)")
	sagaPort = flag.Int("saga-port", 4002, "Service port (Overwriten if SAGA_SERVICE_PORT env var is set)")

	sigiriyaToken = flag.String("sigiriya-token", "", "Token to access sigiriya service.")
)

func main() {
	flag.Parse()
	srv := grpc.NewServer()

	sigiriyaClient := sigiriya.NewClient(&sigiriya.Config{
		Token: *sigiriyaToken,
	})

	// Dial with Kanon
	kanonIP := fmt.Sprintf("%s:%d", *kanonHost, *kanonPort)
	conn, err := grpc.Dial(kanonIP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	ks := kanon.NewServiceClient(conn)

	// Dial with Saga
	sagaIP := fmt.Sprintf("%s:%d", *sagaHost, *sagaPort)
	conn, err = grpc.Dial(sagaIP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	ss := saga.NewServiceClient(conn)

	gs := &GeminiService{
		SigiriyaClient: sigiriyaClient,
		KanonClient:    ks,
		SagaClient:     ss,
	}

	gemini.RegisterServiceServer(srv, gs)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("starting gemini service...")
	log.Println(fmt.Sprintf("listening on: %s:%d", *host, *port))
	srv.Serve(lis)
}

// GeminiService implements gemini interface of gemini package
type GeminiService struct {
	SigiriyaClient *sigiriya.Client
	KanonClient    kanon.ServiceClient
	SagaClient     saga.ServiceClient
}

// Card ...
func (ss *GeminiService) Card(ctx context.Context, r *gemini.Request) (*gemini.Response, error) {
	resp, err := ss.SigiriyaClient.Get(fmt.Sprintf("/cards?user_id=%s", r.UserId))
	if err != nil {
		return nil, err
	}

	var acc *gemini.Account
	if err := json.Unmarshal(resp, acc); err != nil {
		return nil, err
	}

	go ss.KanonClient.GetTransactions(ctx, &kanon.Request{
		ReferenceId: acc.ReferenceId,
	})

	balance, err := ss.SagaClient.GetBalance(ctx, &saga.Request{
		ReferenceId: acc.ReferenceId,
	})
	if err != nil {
		return nil, err
	}

	return &gemini.Response{
		StatusCode: balance.StatusCode,
		Balance: &saga.Balance{
			Email:   balance.Balance.Email,
			Balance: balance.Balance.Balance,
			UserId:  balance.Balance.UserId,
		},
	}, nil
}

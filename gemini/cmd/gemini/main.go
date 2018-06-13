package main

import (
	"encoding/json"
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

// GeminiService implements gemini interface of gemini package
type GeminiService struct {
	SigiriyaClient *sigiriya.Client
	KanonClient    kanon.ServiceClient
	SagaClient     saga.ServiceClient
}

// Card ...
func (ss *GeminiService) Card(ctx context.Context, r *gemini.Request) (*gemini.Response, error) {
	resp, err := ss.SigiriyaClient.Get(fmt.Sprintf("/cards?email=%s", r.Email))
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
		StatusCode: 200,
		Balance: &saga.Balance{
			Email:   balance.Balance.Email,
			Balance: balance.Balance.Balance,
			UserId:  balance.Balance.UserId,
		},
	}, nil
}

func main() {
	srv := grpc.NewServer()

	sigiriyaClient := sigiriya.NewClient(&sigiriya.Config{
		Token: "",
	})

	// Dial with Kanon
	kanonIP := "kanon:4001"
	conn, err := grpc.Dial(kanonIP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	ks := kanon.NewServiceClient(conn)

	// Dial with Saga
	sagaIP := "saga:4002"
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

	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("starting gemini service...")
	log.Println("listening on: 4000")
	srv.Serve(lis)
}

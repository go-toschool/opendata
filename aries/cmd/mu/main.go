package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Finciero/opendata/gemini/gemini"

	"github.com/Finciero/opendata/aries/mu"
	"github.com/Finciero/opendata/aries/sigiriya"
	"google.golang.org/grpc"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if MU_SERVICE_HOST env var is set)")
	port = flag.Int("port", 2000, "Service port (Overwriten if MU_SERVICE_PORT env var is set)")

	geminiHost = flag.String("gemini-host", "gemini", "Service host (Overwriten if GEMINI_SERVICE_HOST env var is set)")
	geminiPort = flag.Int("gemini-port", 4000, "Service port (Overwriten if GEMINI_SERVICE_PORT env var is set)")

	sigiriyaToken = flag.String("sigiriya-token", "", "Token to access sigiriya service.")
)

// MuService ...
type MuService struct {
	SigiriyaClient *sigiriya.Client
	GeminiClient   gemini.ServiceClient
}

// Auth ...
func (ms *MuService) Auth(ctx context.Context, r *mu.Request) (*mu.Response, error) {
	resp, err := ms.SigiriyaClient.Get(fmt.Sprintf("/session?token=%s", r.UserToken))
	if err != nil {
		return nil, err
	}

	var response *gemini.Request
	if err := json.Unmarshal(resp, response); err != nil {
		return nil, err
	}

	response.ClientId = r.PartnerToken
	rr, err := ms.GeminiClient.Card(ctx, response)
	if err != nil {
		return nil, err
	}

	if rr.StatusCode == 200 {
		return &mu.Response{
			StatusCode: rr.StatusCode,
			Message:    "success",
			Balance: &mu.Balance{
				Balance: rr.Balance.Balance,
			},
		}, nil
	}

	return nil, errors.New("Invalid response")
}

func main() {
	flag.Parse()
	srv := grpc.NewServer()

	// Dial with Gemini
	geminiIP := fmt.Sprintf("%s:%d", *geminiHost, *geminiPort)
	conn, err := grpc.Dial(geminiIP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	gs := gemini.NewServiceClient(conn)

	sc := sigiriya.NewClient(&sigiriya.Config{
		Token: *sigiriyaToken,
	})

	ms := &MuService{
		SigiriyaClient: sc,
		GeminiClient:   gs,
	}

	mu.RegisterServiceServer(srv, ms)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Starting mu service...")
	log.Println(fmt.Sprintf("listening on: %s:%d", *host, *port))
	srv.Serve(lis)

	fmt.Println("Service Mu de Aries")
}

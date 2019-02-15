package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-toschool/opendata/gemini"
	"github.com/go-toschool/opendata/gemini/castor"

	"github.com/go-toschool/opendata/aries/mu"
	"github.com/go-toschool/opendata/aries/sigiriya"
	"google.golang.org/grpc"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if MU_SERVICE_HOST env var is set)")
	port = flag.Int("port", 2001, "Service port (Overwriten if MU_SERVICE_PORT env var is set)")

	castorHost = flag.String("castor-host", "extraction.castor", "Service host (Overwriten if CASTOR_SERVICE_HOST env var is set)")
	castorPort = flag.Int("castor-port", 4000, "Service port (Overwriten if CASTOR_SERVICE_PORT env var is set)")
	castorCert = flag.String("castor-cert", "", "Castor cert (Overwriten if CASTOR_CERT env var is set)")

	sigiriyaToken = flag.String("sigiriya-token", "", "Token to access sigiriya service.")
)

func main() {
	flag.Parse()
	srv := grpc.NewServer()

	// Dial with Gemini
	cs := gemini.NewCastor(&gemini.Config{
		Host: *castorHost,
		Port: *castorPort,
		Cert: *castorCert,
	})

	sc := sigiriya.NewClient(&sigiriya.Config{
		Token: *sigiriyaToken,
	})

	ms := &MuService{
		SigiriyaClient: sc,
		CastorClient:   cs,
	}

	mu.RegisterServiceServer(srv, ms)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting mu service...")
	log.Println(fmt.Sprintf("listening on: %d", *port))
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// MuService ...
type MuService struct {
	SigiriyaClient *sigiriya.Client
	CastorClient   castor.ServiceClient
}

type sigiriyaResponse struct {
	ID              string `json:"id,omitempty"`
	VanityNumber    string `json:"vanity_number,omitempty"`
	ReferenceID     string `json:"reference_id,omitempty"`
	ReferenceUserID string `json:"reference_user_id,omitempty"`
}

// Extract ...
func (ms *MuService) Extract(ctx context.Context, r *mu.Request) (*mu.Response, error) {
	ms.SigiriyaClient.SetUserToken(r.UserToken)
	resp, err := ms.SigiriyaClient.Get("/auth")
	if err != nil {
		return nil, err
	}

	var srr *sigiriyaResponse
	if err := json.Unmarshal(resp, srr); err != nil {
		return nil, err
	}

	rr, err := ms.CastorClient.Card(ctx, &castor.Request{
		ClientId:    r.PartnerToken,
		UserId:      srr.ReferenceUserID,
		ReferenceId: srr.ReferenceID,
	})
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
			Account: &mu.Account{
				Id:       srr.ID,
				VanityId: srr.VanityNumber,
			},
		}, nil
	}

	return nil, errors.New("invalid response")
}

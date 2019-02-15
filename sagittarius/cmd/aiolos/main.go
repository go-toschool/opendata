package main

import (
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/go-toschool/opendata/sagittarius/aiolos"
)

func main() {
	srv := grpc.NewServer()
	ss := &SaggitariusService{}

	aiolos.RegisterServiceServer(srv, ss)

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("starting sagittarius service...")
	log.Println("listening on: 3000")
	srv.Serve(lis)
}

// SaggitariusService implements saggitarios interface of aiolos package
type SaggitariusService struct {
}

// Create ...
func (ss *SaggitariusService) Create(ctx context.Context, c *aiolos.Callback) (*aiolos.Response, error) {
	return &aiolos.Response{}, nil
}

// Update ...
func (ss *SaggitariusService) Update(ctx context.Context, c *aiolos.Callback) (*aiolos.Response, error) {
	return &aiolos.Response{}, nil
}

// Delete ...
func (ss *SaggitariusService) Delete(ctx context.Context, c *aiolos.Callback) (*aiolos.Response, error) {
	return &aiolos.Response{}, nil
}

// Dispatch ...
func (ss *SaggitariusService) Dispatch(ctx context.Context, r *aiolos.Request) (*aiolos.Response, error) {
	// Get callback url
	fmt.Println(r.Id)
	// dispatch data
	return &aiolos.Response{
		Message: "hello ",
	}, nil
}

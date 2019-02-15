package aries

import (
	"fmt"
	"log"

	"github.com/go-toschool/opendata/aries/mu"
	"google.golang.org/grpc"
)

type Config struct {
	Host string
	Port int
}

func NewMu(c *Config) mu.ServiceClient {
	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	return mu.NewServiceClient(conn)
}

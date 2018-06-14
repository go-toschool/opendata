package taurus

import (
	"fmt"
	"log"

	"github.com/Finciero/opendata/taurus/aldebaran"
	"google.golang.org/grpc"
)

type Config struct {
	Host string
	Port int
}

func NewMu(c *Config) aldebaran.ServiceClient {
	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return aldebaran.NewServiceClient(conn)
}

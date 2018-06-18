package taurus

import (
	"fmt"
	"log"

	"github.com/Finciero/opendata/taurus/aldebaran"
	"google.golang.org/grpc"
)

// Config hold data to connect to this service
type Config struct {
	Host string
	Port int
}

func NewAldebaran(c *Config) aldebaran.ServiceClient {
	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return aldebaran.NewServiceClient(conn)
}

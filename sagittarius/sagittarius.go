package sagittarius

import (
	"fmt"
	"log"

	"github.com/Finciero/opendata/sagittarius/aiolos"
	"google.golang.org/grpc"
)

type Config struct {
	Host string
	Port int
}

func NewAiolos(c *Config) aiolos.ServiceClient {
	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return aiolos.NewServiceClient(conn)
}

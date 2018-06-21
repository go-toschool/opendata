package capricornius

import (
	"fmt"
	"log"

	"github.com/Finciero/opendata/capricornius/shura"
	"google.golang.org/grpc"
)

type Config struct {
	Host string
	Port int
}

func NewShura(c *Config) shura.ServiceClient {
	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	return shura.NewServiceClient(conn)
}

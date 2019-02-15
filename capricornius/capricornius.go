package capricornius

import (
	"fmt"
	"log"

	"github.com/go-toschool/opendata/capricornius/shura"
	"google.golang.org/grpc"
)

// Config ...
type Config struct {
	Host string
	Port int
}

// NewShura ...
func NewShura(c *Config) shura.ServiceClient {
	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	return shura.NewServiceClient(conn)
}

package gemini

import (
	"fmt"
	"log"

	"github.com/Finciero/opendata/gemini/castor"
	"github.com/Finciero/opendata/gemini/kanon"
	"github.com/Finciero/opendata/gemini/saga"
	"google.golang.org/grpc"
)

type Config struct {
	Host string
	Port int
}

func NewGemini(c *Config) castor.ServiceClient {
	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return castor.NewServiceClient(conn)
}

func NewKanon(c *Config) kanon.ServiceClient {
	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return kanon.NewServiceClient(conn)
}

func NewSaga(c *Config) saga.ServiceClient {
	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return saga.NewServiceClient(conn)
}

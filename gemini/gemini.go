package gemini

import (
	"fmt"
	"log"

	"github.com/Finciero/opendata/gemini/castor"
	"github.com/Finciero/opendata/gemini/kanon"
	"github.com/Finciero/opendata/gemini/saga"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	Host string
	Port int
	Cert string
}

func NewCastor(c *Config) castor.ServiceClient {
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(c.Cert, "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	// defer conn.Close()
	return castor.NewServiceClient(conn)
}

func NewKanon(c *Config) kanon.ServiceClient {
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(c.Cert, "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	// defer conn.Close()
	return kanon.NewServiceClient(conn)
}

func NewSaga(c *Config) saga.ServiceClient {
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(c.Cert, "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	IP := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := grpc.Dial(IP, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	// defer conn.Close()
	return saga.NewServiceClient(conn)
}

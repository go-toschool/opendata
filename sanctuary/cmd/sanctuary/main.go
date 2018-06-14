package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Finciero/opendata/aries/mu"
	"github.com/Finciero/opendata/sanctuary/cmd/sanctuary/api"
	"github.com/Finciero/opendata/sanctuary/cmd/sanctuary/auth"
	"github.com/Finciero/opendata/taurus/aldebaran"
	"github.com/urfave/negroni"
	"google.golang.org/grpc"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if SANCTUARY_SERVICE_HOST env var is set)")
	port = flag.Int("port", 2000, "Service port (Overwriten if SANCTUARY_SERVICE_PORT env var is set)")

	muHost = flag.String("mu-host", "mu", "Service host (Overwriten if MU_SERVICE_HOST env var is set)")
	muPort = flag.Int("mu-port", 2001, "Service port (Overwriten if MU_SERVICE_PORT env var is set)")

	aldebaranHost = flag.String("aldebaran-host", "aldebaran", "Service host (Overwriten if ALDEBARAN_SERVICE_HOST env var is set)")
	aldebaranPort = flag.Int("aldebaran-port", 2002, "Service port (Overwriten if ALDEBARAN_SERVICE_PORT env var is set)")
)

func main() {
	flag.Parse()
	srv := grpc.NewServer()

	// Dial with Mu
	muIP := fmt.Sprintf("%s:%d", *muHost, *muPort)
	conn, err := grpc.Dial(muIP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	ms := mu.NewServiceClient(conn)

	// Dial with Aldebaran
	aldebaranIP := fmt.Sprintf("%s:%d", *aldebaranHost, *aldebaranPort)
	conn, err = grpc.Dial(aldebaranIP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	as := aldebaran.NewServiceClient(conn)

	authCtx := &auth.Context{
		AldebaranClient: as,
	}

	apiCtx := &api.Context{
		MuClient: ms,
	}

	n := negroni.New(
		negroni.NewRecovery(),
	)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api.Handle(aCtx)))
	mux.Handle("/auth/", http.StripPrefix("/auth", auth.Handle(aCtx)))
	n.UseHandler(mux)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	logger.Printf("Start listening on %s", addr)
	logger.Fatal(http.ListenAndServe(addr, n))
}

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Finciero/opendata/aries"
	"github.com/Finciero/opendata/sanctuary/cmd/sanctuary/api"
	"github.com/Finciero/opendata/sanctuary/cmd/sanctuary/auth"
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
	ms := aries.NewMu(&aries.Config{
		Host: *muHost,
		Port: *muPort,
	})

	// Dial with Aldebaran
	as := taurus.NewAldebaran(&taurus.Config{
		Host: *aldebaranHost,
		Port: *aldebaranPort,
	})

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

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/Finciero/opendata/capricornius"
	"github.com/Finciero/opendata/capricornius/shura"

	"github.com/Finciero/opendata/aries"
	"github.com/Finciero/opendata/sanctuary/cmd/sanctuary/api"
	"github.com/Finciero/opendata/sanctuary/cmd/sanctuary/auth"
	"github.com/Finciero/opendata/taurus"
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

	shuraHost = flag.String("shura-host", "shura", "Service host (Overwriten if SHURA_SERVICE_HOST env var is set)")
	shuraPort = flag.Int("shura-port", 2002, "Service port (Overwriten if SHURA_SERVICE_PORT env var is set)")
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

	// Dial with Shura
	ss := capricornius.NewShura(&capricornius.Config{
		Host: *shuraHost,
		Port: *shuraPort,
	})
	ctx := context.Background()
	origins, err := ss.GetOrigins(ctx, &shura.Origin{})
	if err != nil {
		panic(err)
	}

	authCtx := &auth.Context{
		AldebaranClient: as,
		ShuraClient:     ss,
		AllowedOrigins:  origins.AllowedOrigins,
	}

	apiCtx := &api.Context{
		ShuraClient:    ss,
		MuClient:       ms,
		AllowedOrigins: origins.AllowedOrigins,
	}

	n := negroni.New(
		negroni.NewRecovery(),
	)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api.Routes(apiCtx)))
	mux.Handle("/auth/", http.StripPrefix("/auth", auth.Routes(authCtx)))
	n.UseHandler(mux)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	logger.Printf("Start listening on %s", addr)
	logger.Fatal(http.ListenAndServe(addr, n))
}

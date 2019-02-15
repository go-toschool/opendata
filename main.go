package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/go-toschool/opendata/gemini"
	"github.com/go-toschool/opendata/gemini/castor"
	"google.golang.org/grpc/metadata"
)

var (
	port = flag.Int("port", 8080, "port to contact service")
)

func main() {
	flag.Parse()

	gc := gemini.NewCastor(&gemini.Config{
		Host: "localhost",
		Port: *port,
		Cert: "certs/castor/castor.crt",
	})

	ctx := context.Background()
	m := map[string]string{
		"X-FIN-CLIENT-ID": "rodrwan",
	}
	md := metadata.New(m)
	ctx = metadata.NewOutgoingContext(ctx, md)

	rr, err := gc.Card(ctx, &castor.Request{
		ClientId:    "Rodrwan",
		UserId:      "user-id-241243234",
		ReferenceId: "21423424",
	})
	if err != nil {
		log.Fatalf("Could not retrieve cards: %v", err)
	}

	fmt.Println(rr.StatusCode)
	fmt.Printf("Balance: %f\n", rr.Balance.Balance)
	fmt.Printf("UserID: %s\n", rr.Balance.UserId)
	fmt.Printf("Account ID: %s\n", rr.Balance.AccountId)
}

// func main() {
// 	flag.Parse()

// 	bodyRequest := &sanctuary.ExtractionRequest{
// 		UserToken: "user-token",
// 	}

// 	bodyByte, err := json.Marshal(bodyRequest)
// 	if err != nil {
// 		log.Fatalf("could not marshal struct, %v", err)
// 	}

// 	URL := fmt.Sprintf("http://localhost:%d/api/extract", *port)
// 	req, err := http.NewRequest("POST", URL, strings.NewReader(string(bodyByte)))
// 	if err != nil {
// 		log.Fatalf("could not create request, %v", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("X-FIN-CLIENT-TOKEN", "Bearer user-token")

// 	httpClient := &http.Client{}

// 	resp, err := httpClient.Do(req)
// 	if err != nil {
// 		log.Fatalf("could not do http request, %v", err)
// 	}

// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("could not read response body, %v", err)
// 	}

// 	fmt.Println(string(body))
// }

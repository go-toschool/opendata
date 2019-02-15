package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	sigiriyaPkg "github.com/go-toschool/sigiriya"

	"golang.org/x/net/context"

	"github.com/go-toschool/opendata/taurus/aldebaran"
	"github.com/go-toschool/opendata/taurus/sigiriya"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if ALDEBARAN_SERVICE_HOST env var is set)")
	port = flag.Int("port", 2002, "Service port (Overwriten if ALDEBARAN_SERVICE_PORT env var is set)")

	aldebaranCert = flag.String("aldebaran-cert", "", "Aldebaran cert (Overwriten if ALDEBARAN_CERT env var is set)")
	aldebaranKey  = flag.String("aldebaran-key", "", "Aldebaran key (Overwriten if ALDEBARAN_KEY env var is set)")

	sigiriyaToken = flag.String("sigiriya-token", "", "Token to access sigiriya service.")

	withTLS = flag.Bool("with-tls", false, "service with TLS")
)

func main() {
	flag.Parse()

	as := &AldebaranService{}

	var srv *grpc.Server
	if *withTLS {
		creds, err := credentials.NewServerTLSFromFile(*aldebaranCert, *aldebaranKey)
		if err != nil {
			log.Fatalf("could not load TLS keys: %s", err)
		}
		opts := []grpc.ServerOption{grpc.Creds(creds)}
		srv = grpc.NewServer(opts...)
	} else {
		srv = grpc.NewServer()
	}

	aldebaran.RegisterServiceServer(srv, as)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Starting mu service...")
	log.Println(fmt.Sprintf("listening on: %d", *port))

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// AldebaranService implements aldebaran interface.
type AldebaranService struct {
	SigiriyaClient *sigiriya.Client
}

// CreateToken create a new Auth token
func (as *AldebaranService) CreateToken(ctx context.Context, r *aldebaran.Create) (*aldebaran.CreateResponse, error) {
	now := time.Now()
	text := fmt.Sprintf("%s#%s#%s", r.Email, r.ClientToken, now.String())
	token, err := encryptString(text, r.ClientToken)
	if err != nil {
		return nil, err
	}

	data := &sigiriyaPkg.AuthApplication{
		Email:     r.Email,
		Client:    r.ClientToken,
		Token:     token,
		CreatedAt: now,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := as.SigiriyaClient.Post("/auth", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	fmt.Println(string(resp))

	return &aldebaran.CreateResponse{
		StatusCode: 200,
		Message:    "Token generated",
		Token:      token,
	}, nil
}

// CheckToken validate user token
func (as *AldebaranService) CheckToken(ctx context.Context, r *aldebaran.Check) (*aldebaran.CheckResponse, error) {
	text, err := decryptString(r.UserToken, r.ClientToken)
	if err != nil {
		return nil, err
	}

	resp, err := as.SigiriyaClient.Get(fmt.Sprintf("/auth?token=%s", r.UserToken))
	if err != nil {
		return nil, err
	}

	payload := new(struct {
		StatusCode int32                       `json:"status_code"`
		Data       sigiriyaPkg.AuthApplication `json:"data"`
	})
	if err := json.Unmarshal(resp, &payload); err != nil {
		return nil, err
	}

	if payload.StatusCode == 404 {
		return &aldebaran.CheckResponse{
			StatusCode: 404,
			Message:    "Invalid token",
			Valid:      false,
		}, nil
	}

	splitted := strings.Split(text, "#")
	cutOffTime, err := time.Parse(time.RFC3339, splitted[2])
	if err != nil {
		return nil, err
	}

	valid := strings.Contains(splitted[0], payload.Data.Email) &&
		strings.Contains(splitted[1], payload.Data.Token) &&
		payload.Data.CreatedAt == cutOffTime

	return &aldebaran.CheckResponse{
		StatusCode: 200,
		Message:    "Validated token",
		Valid:      valid,
	}, nil
}

// Takes two strings, cryptoText and keyString.
// cryptoText is the text to be decrypted and the keyString is the key to use for the decryption.
// The function will output the resulting plain text string with an error variable.
func decryptString(cryptoText string, keyString string) (plainTextString string, err error) {
	encrypted, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	if len(encrypted) < aes.BlockSize {
		return "", fmt.Errorf("cipherText too short. It decodes to %v bytes but the minimum length is 16", len(encrypted))
	}

	decrypted, err := decryptAES(hashTo32Bytes(keyString), encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func decryptAES(key, data []byte) ([]byte, error) {
	// split the input up in to the IV seed and then the actual encrypted data.
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(data, data)
	return data, nil
}

// Takes two string, plainText and keyString.
// plainText is the text that needs to be encrypted by keyString.
// The function will output the resulting crypto text and an error variable.
func encryptString(plainText string, keyString string) (cipherTextString string, err error) {
	key := hashTo32Bytes(keyString)
	encrypted, err := encryptAES(key, []byte(plainText))
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func encryptAES(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// create two 'windows' in to the output slice.
	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	// populate the IV slice with random data.
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	// note that encrypted is still a window in to the output slice
	stream.XORKeyStream(encrypted, data)
	return output, nil
}

func hashTo32Bytes(input string) []byte {
	data := sha256.Sum256([]byte(input))
	return data[0:]
}

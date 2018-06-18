package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/Finciero/opendata/taurus/aldebaran"
	"google.golang.org/grpc"
)

var (
	host = flag.String("host", "", "Service host (Overwriten if ALDEBARAN_SERVICE_HOST env var is set)")
	port = flag.Int("port", 2000, "Service port (Overwriten if ALDEBARAN_SERVICE_PORT env var is set)")

	sigiriyaToken = flag.String("sigiriya-token", "", "Token to access sigiriya service.")
)

func main() {
	flag.Parse()
	srv := grpc.NewServer()

	as := &AldebaranService{}

	aldebaran.RegisterServiceServer(srv, as)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Starting mu service...")
	log.Println(fmt.Sprintf("listening on: %s:%d", *host, *port))
	srv.Serve(lis)

	fmt.Println("Service Mu de Aries")
}

// AldebaranService implements aldebaran interface.
type AldebaranService struct {
}

// CreateToken create a new Auth token
func (as *AldebaranService) CreateToken(ctx context.Context, r *aldebaran.Request) (*aldebaran.Response, error) {
	return &aldebaran.Response{}, nil
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

// As we cannot use a variable length key, we must cut the users key
// up to or down to 32 bytes. To do this the function takes a hash
// of the key and cuts it down to 32 bytes.
func hashTo32Bytes(input string) []byte {
	data := sha256.Sum256([]byte(input))
	return data[0:]
}

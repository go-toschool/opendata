package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"unicode"

	"github.com/Finciero/opendata/taurus/aldebaran"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Context contains all services and variables of the applications.
type Context struct {
	TaurusClient aldebaran.ServiceClient
}

// Handle creates a new bounded Handler with context.
func (c *Context) Handle(h HandlerFunc) *Handler {
	return &Handler{c, h}
}

// ToJSON unmarshal the given reader into v.
func (c *Context) ToJSON(r io.Reader, v interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// NfcString clean characters
func (c *Context) NfcString(nfd string) (nfc string) {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	nfc, _, _ = transform.String(t, nfd)
	return
}

// Response represents a response of the Application REST API.
type Response struct {
	Status int         `json:"-"`
	Data   interface{} `json:"data,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
}

// Write writes a ApplicationResposne to the given response writer encoded as JSON.
func (r *Response) Write(w http.ResponseWriter) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	_, err = w.Write(b)
	return err
}

// HandlerFunc function handler signature used by brickwall application.
type HandlerFunc func(*Context, http.ResponseWriter, *http.Request) (*Response, error)

// Handler is an http.Handler that provides access to the Context to the given HandlerFunc.
type Handler struct {
	ctx    *Context
	handle HandlerFunc
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	resp, err := h.handle(h.ctx, w, r)
	if err != nil {
		fmt.Printf("[ERROR]: unexpected unhandled error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := resp.Write(w); err != nil {
		fmt.Printf("[ERROR]: %v, encoding response: %v", err, resp)
	}
}

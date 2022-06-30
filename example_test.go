package jackett_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/SkYNewZ/go-jackett"
)

func Example() {
	// Read JACKETT_API_URL and JACKETT_API_KEY by default
	client := jackett.New(nil)
	results, err := client.Fetch(context.Background(), jackett.WithQuery("foo"))
	if err != nil {
		panic(err)
	}

	_ = results // use results
}

func Example_disableSSL() {
	httpClient := &http.Client{Timeout: time.Minute}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient.Transport = transport

	client := jackett.New(&jackett.Config{
		APIUrl:     "https://",
		APIKey:     "foo",
		HTTPClient: httpClient,
	})
	results, err := client.Fetch(context.Background(), jackett.WithQuery("foo"))
	if err != nil {
		panic(err)
	}

	_ = results // use results
}

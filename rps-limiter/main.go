package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/kogutich/go-concurrency-test/rps-limiter/transport"
)

func main() {
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport.NewRPSLimitedTransport(nil, 10),
	}

	r, err := client.Get("https://httpbin.org/status/200")

	if err == nil {
		if fullResp, err := httputil.DumpResponse(r, true); err == nil {
			fmt.Printf("%s\n", fullResp)
		}
	}
}

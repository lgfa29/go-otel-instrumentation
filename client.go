package main

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func makeRequest(ctx context.Context, url string, name string) error {
	// Run HTTP client.
	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport),
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create client request: %v", err)
	}

	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to run client request: %v", err)
	}
	resp.Body.Close()

	fmt.Printf("User registration request for %q sent successfully\n", name)
	return nil
}

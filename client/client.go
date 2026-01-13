package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	urls := []string{
		"http://localhost:8088/v1/books/1",
		"http://localhost:8088/v1/books/2",
	}

	for _, url := range urls {
		go func(reqURL string) {
			defer wg.Done()
			defer func() {
				cancel()
				fmt.Printf("reqURL %s has been canceled.\n", reqURL)
			}()

			var resp Response
			if err := doRequest(ctx, reqURL, http.MethodGet, nil, &resp); err != nil {
				return
			}

			fmt.Printf("%s has received reponse successfully\n", reqURL)
			fmt.Printf("Response: %+v\n", resp)
		}(url)
	}

	wg.Wait()
	fmt.Println("Any request-response has been done.")
}

func doRequest(ctx context.Context, reqURL string, method string, body io.Reader, ret interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(ret); err != nil {
		return err
	}

	return nil
}

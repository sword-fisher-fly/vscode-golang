package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestResponse struct {
	Message string `json:"message"`
	Value   int    `json:"value"`
}

// go.exe test -test.fullpath=true -timeout 30s -run ^TestDoRequest$ github.com/sword-fisher-fly/vscode-golang/client
func TestDoRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := TestResponse{Message: "success", Value: 42}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	t.Run("Normal request with valid response", func(t *testing.T) {
		var result TestResponse
		err := doRequest(context.Background(), server.URL, "GET", nil, &result)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.Message != "success" {
			t.Errorf("Expected message 'success', got '%s'", result.Message)
		}

		if result.Value != 42 {
			t.Errorf("Expected value 42, got %d", result.Value)
		}
	})
}

// TestDoRequestWithTimeout test timeout function.
func TestDoRequestWithTimeout(t *testing.T) {
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		response := TestResponse{Message: "slow", Value: 100}
		json.NewEncoder(w).Encode(response)
	}))
	defer slowServer.Close()

	// 使用带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	var result TestResponse
	err := doRequest(ctx, slowServer.URL, "GET", nil, &result)

	if err == nil {
		t.Error("Expected a timeout error, but got none")
	}
}

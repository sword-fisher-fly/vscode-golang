package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// go test -v -coverprofile=coverage.out && go tool cover -html=coverage.out
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// /opt/go/bin/go test -test.fullpath=true -timeout 30s -run ^TestGetBookByID$ gitlab.lionaitech.com/golang-vscode/server
// go test -timeout 30s -run ^TestGetBookByID$ ./server -v
func TestGetBookByID(t *testing.T) {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/books/:id", getBookByID)
	}

	testCases := []struct {
		name           string
		bookID         string
		expectedStatus int
		expectedCode   int
		expectedTitle  string
		expectedDelay  time.Duration
	}{
		{
			name:           "获取存在的图书（ID=1，带延时）",
			bookID:         "1",
			expectedStatus: http.StatusOK,
			expectedCode:   200,
			expectedTitle:  "《Go语言实战》",
			expectedDelay:  10 * time.Second,
		},
		{
			name:           "获取存在的图书（ID=2，无延时）",
			bookID:         "2",
			expectedStatus: http.StatusOK,
			expectedCode:   200,
			expectedTitle:  "《Gin框架入门与实战》",
			expectedDelay:  0,
		},
		{
			name:           "获取不存在的图书（ID=4）",
			bookID:         "4",
			expectedStatus: http.StatusNotFound,
			expectedCode:   404,
			expectedTitle:  "",
			expectedDelay:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/books/"+tc.bookID, nil)
			startTime := time.Now()
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			elapsed := time.Since(startTime)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			// err := gin.UnmarshalJSON(w.Body.Bytes(), &response)
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "解析响应 JSON 失败")

			assert.Equal(t, float64(tc.expectedCode), response["code"])
			if tc.expectedTitle != "" {
				data := response["data"].(map[string]interface{})
				assert.Equal(t, tc.expectedTitle, data["title"])
				assert.Equal(t, tc.bookID, data["id"])
			} else {
				assert.Nil(t, response["data"])
				assert.Equal(t, "Book not found", response["message"])
			}

			if tc.expectedDelay > 0 {
				assert.GreaterOrEqual(t, elapsed, tc.expectedDelay, "ID=1 的请求延时不足10秒")
			} else {
				assert.Less(t, elapsed, 1*time.Second, "非ID=1的请求不应有长延时")
			}
		})
	}
}

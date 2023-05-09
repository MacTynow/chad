package cmd

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockImageResponseBody struct {
	mock.Mock
}

func (m *mockImageResponseBody) Error() OpenAIError {
	args := m.Called()
	return args.Get(0).(OpenAIError)
}

type mockChatResponseBody struct {
	mock.Mock
}

func (m *mockChatResponseBody) Error() OpenAIError {
	args := m.Called()
	return args.Get(0).(OpenAIError)
}

type mockEditResponseBody struct {
	mock.Mock
}

func (m *mockEditResponseBody) Error() OpenAIError {
	args := m.Called()
	return args.Get(0).(OpenAIError)
}

// func APIhandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	switch r.URL.Path {
// 	case "/images":
// 		w.Write([]byte(`{"created": 1622096000, "data": [{"url": "https://cdn.openai.com/blah"}]}`))
// 	case "/chat":
// 		w.Write([]byte(`{"choices": [{"text": "Hello, world!"}]}`))
// 	case "/edit":
// 		w.Write([]byte(`{"choices": [{"text": "This is edited text."}]}`))
// 	case "/invalid":
// 		w.Write([]byte(`invalid response body`))
// 	default:
// 		w.WriteHeader(http.StatusNotFound)
// 	}
// }

func TestSendRequesttoOpenAI(t *testing.T) {
	testCases := []struct {
		name           string
		requestURL     string
		requestBody    RequestBody
		expectedResult string
		expectedError  error
	}{
		{
			name:           "invalid request body",
			requestURL:     "http://localhost:8080/images",
			requestBody:    make(chan int),
			expectedResult: "",
			expectedError:  &json.MarshalerError{},
		},
		{
			name:           "invalid request",
			requestURL:     ":",
			requestBody:    nil,
			expectedResult: "",
			expectedError:  errors.New("parse :: missing protocol scheme"),
		},
		{
			name:           "missing api key",
			requestURL:     "http://localhost:8080/images",
			requestBody:    &ImageResponseBody{},
			expectedResult: "",
			expectedError:  errors.New("Please set the OPENAI_API_KEY environment variable"),
		},
		{
			name:           "image response success",
			requestURL:     "http://localhost:8080/images",
			requestBody:    &ImageRequestBody{},
			expectedResult: "https://example.com/image.jpg",
			expectedError:  nil,
		},
		{
			name:           "chat response success",
			requestURL:     "http://localhost:8080/chat",
			requestBody:    &ChatRequestBody{},
			expectedResult: "Hello, world!",
			expectedError:  nil,
		},
		{
			name:           "edit response success",
			requestURL:     "http://localhost:8080/edit",
			requestBody:    &EditRequestBody{},
			expectedResult: "This is edited text.",
			expectedError:  nil,
		},
		{
			name:           "invalid response",
			requestURL:     "http://localhost:8080/invalid",
			requestBody:    &ChatRequestBody{},
			expectedResult: "invalid response body",
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock response body based on request URL
			var responseBody ResponseBody
			switch {
			case strings.Contains(tc.requestURL, "images"):
				responseBody = &mockImageResponseBody{}
			case strings.Contains(tc.requestURL, "chat"):
				responseBody = &mockChatResponseBody{}
			case strings.Contains(tc.requestURL, "edit"):
				responseBody = &mockEditResponseBody{}
			default:
				responseBody = nil
			}

			// Set up mock HTTP server to handle requests
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.requestURL, r.URL.String())

				// Assert that the request body matches the expected value
				var requestBody RequestBody
				err := json.NewDecoder(r.Body).Decode(&requestBody)
				require.NoError(t, err)
				assert.Equal(t, tc.requestBody, requestBody)

				// Set up mock response
				responseData := map[string]interface{}{}
				switch responseBody.(type) {
				case *mockImageResponseBody:
					responseData = map[string]interface{}{
						"created": 0,
						"data": []interface{}{
							map[string]interface{}{
								"url": "https://example.com/image.jpg",
							},
						},
					}
				case *mockChatResponseBody:
					responseData = map[string]interface{}{
						"choices": []interface{}{
							map[string]interface{}{
								"message": map[string]interface{}{
									"content": "Hello, world!",
								},
							},
						},
					}
				case *mockEditResponseBody:
					responseData = map[string]interface{}{
						"choices": []interface{}{
							map[string]interface{}{
								"text": "This is edited text.",
							},
						},
					}
				}
				responseBodyJSON, err := json.Marshal(responseData)
				require.NoError(t, err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err = w.Write(responseBodyJSON)
				require.NoError(t, err)
			}))
			defer ts.Close()

			// Set up environment variable
			os.Setenv("OPENAI_API_KEY", "test-api-key")

			// Call function under test
			result, err := sendRequesttoOpenAI(ts.URL+strings.TrimPrefix(tc.requestURL, "http://localhost:8080"), tc.requestBody)

			// Assert that the result and error match the expected values
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

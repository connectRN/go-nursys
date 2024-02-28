package nursys

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * This suite tests the Nursys client itself, the constructors and the  invokeEndpoint endpoint
 * verifying that it sends the correct JSON request bodies to the right endpoints.
 */

// Test basic operation
func TestInvokeEndpoint_Basic(t *testing.T) {
	assert := assert.New(t)

	expectedBody := `{ "message": "Hello World" }`

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Validate request parameters and body
		assert.Equal("POST", req.Method)
		assert.Equal("/endpoint", req.URL.String())
		assert.Equal("application/json", req.Header.Get("Content-Type"))
		assert.Equal("acme", req.Header.Get("username"))
		assert.Equal("1234!", req.Header.Get("password"))
		bodyBytes, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.JSONEq(expectedBody, string(bodyBytes), "Actual: %s", bodyBytes)
		// Write response
		rw.Write([]byte(`{"ok": true,"operation_id": "df6a6b50","push_ids": ["9d78a53b"],"message_ids": [], "content_urls": []}`))
	}))
	t.Cleanup(server.Close)

	testConnection := New(server.URL, "acme", "1234!").(*nsHTTPClient)

	type TestResponseMessage struct {
		OperationID string `json:"operation_id"`
	}

	// Invoke!
	var result TestResponseMessage
	body := map[string]string{"message": "Hello World"}
	err := testConnection.invokeEndpoint(context.Background(), http.MethodPost, "/endpoint", body, &result)
	require.Nil(t, err)
	assert.Equal(TestResponseMessage{OperationID: "df6a6b50"}, result)
}

// Make sure HTTP error codes returned by httpClient.Do() don't panic
func TestInvokeEndpoint_HttpError(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(`{"error": "Forbidden"}`))
	}))
	t.Cleanup(server.Close)

	testConnection := New("https://example.com/api", "acme", "1234!").(*nsHTTPClient)

	// Invoke!
	body := map[string]string{"message": "Hello World"}
	err := testConnection.invokeEndpoint(context.Background(), http.MethodPost, "/endpoint", body, nil)
	assert.Error(err) // HTTP errors return a go error
}

// Make sure errors returned by httpClient.Do() don't panic and are returned.
func TestInvokeEndpoint_GoError(t *testing.T) {
	assert := assert.New(t)

	client := &http.Client{Transport: RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("oh no an error")
	})}

	testConnection := New("https://example.com/api", "acme", "1234!", WithHTTPClient(client)).(*nsHTTPClient)

	// Invoke!
	body := map[string]string{"message": "Hello World"}
	err := testConnection.invokeEndpoint(context.Background(), http.MethodPost, "/endpoint", body, nil)
	assert.Error(err)
}

// The RoundTripperFunc type is an adapter to allow the use of ordinary functions as HTTP Client transports.
// This really ought to already be in the net/http/httptest package, it's just like http.HandlerFunc
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

// RoundTrip makes RoundTripFunc implement http.RoundTripper. It calls f(req)
func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Make sure cancelling the context cancels the operation.
func TestInvokeEndpoint_ContextCancel(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(5 * time.Second) // Take a long time
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(`{"error": "Forbidden"}`))
	}))
	t.Cleanup(server.Close)

	testConnection := New("https://example.com/api", "acme", "1234!").(*nsHTTPClient)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	t.Cleanup(cancel)

	// Invoke!
	start := time.Now()
	body := map[string]string{"message": "Hello World"}
	err := testConnection.invokeEndpoint(ctx, http.MethodPost, "/endpoint", body, nil)
	assert.ErrorIs(err, context.DeadlineExceeded) // Should return this one!
	assert.Less(time.Since(start), 2*time.Second) // Just extra sure that we didn't wait 5 seconds
}

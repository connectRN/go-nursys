package nursys_test

import (
	"context"
	_ "embed"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/connectRN/go-nursys"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/Generic_SubmitResponseMessage.json
var submitResponseJSON []byte

func Test_ChangePassword(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	const chpwSubmitExpectedRequestJSON = `{"NewPassword": "MyN3wR34llyStr0ngAP1P4ssw0rd$1!0"}`
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Validate request parameters and body
		assert.Equal("POST", req.Method)
		assert.Equal("/changepassword", req.URL.String())
		assertBodyJSONEqual(t, chpwSubmitExpectedRequestJSON, req.Body)

		rw.Write(submitResponseJSON)
	}))
	t.Cleanup(server.Close)

	testConnection := nursys.New(server.URL, "acme", "1234!")

	request := nursys.ChangePasswordSubmitRequestMessage{
		NewPassword: "MyN3wR34llyStr0ngAP1P4ssw0rd$1!0",
	}

	postResp, err := testConnection.ChangePassword(ctx, request)
	require.NoError(t, err)

	assert.Equal("a523e0d4-01e1-4c8d-8dd9-54b269c315b7", postResp.TransactionID)
	assert.Equal(time.Date(2024, 1, 4, 10, 18, 38, 966883900, time.FixedZone("", -21600)), postResp.TransactionDate)
	assert.Equal("Submit lookup successful. ", postResp.TransactionComment)
	assert.Equal(true, postResp.TransactionSuccessFlag)
	assert.Equal([]nursys.TransactionError{}, postResp.TransactionErrors)
}

//go:embed testdata/ChangePasswordSubmitResponseMessage_Failed.json
var chpwSubmitFailedResponseJSON []byte

func Test_ChangePassword_Failed(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	const chpwSubmitExpectedRequestJSON = `{"NewPassword": "short"}`
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Validate request parameters and body
		assert.Equal("POST", req.Method)
		assert.Equal("/changepassword", req.URL.String())
		assertBodyJSONEqual(t, chpwSubmitExpectedRequestJSON, req.Body)

		rw.Write(chpwSubmitFailedResponseJSON)
	}))
	t.Cleanup(server.Close)

	testConnection := nursys.New(server.URL, "acme", "1234!")

	request := nursys.ChangePasswordSubmitRequestMessage{
		NewPassword: "short",
	}

	postResp, err := testConnection.ChangePassword(ctx, request)
	require.NoError(t, err)

	assert.Equal("xfaca386-x034-41x8-90xf-96xd23318348", postResp.TransactionID)
	assert.Equal(time.Date(2021, 8, 31, 15, 47, 0, 265942100, time.FixedZone("", -18000)), postResp.TransactionDate)
	assert.Equal("", postResp.TransactionComment)
	assert.Equal(false, postResp.TransactionSuccessFlag)
	assert.Equal([]nursys.TransactionError{
		{
			ErrorID:      210,
			ErrorMessage: "Password must be between 8 and 50 characters in length. ",
		},
	}, postResp.TransactionErrors)

}

// Helper to read and compare the request body.
func assertBodyJSONEqual(t testing.TB, expected string, body io.ReadCloser, msgAndArgs ...interface{}) bool {
	// Read the body and check for error while reading.
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	require.Nil(t, err, msgAndArgs...)
	return assert.JSONEq(t, expected, string(bodyBytes), msgAndArgs...)
}

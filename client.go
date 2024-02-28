package nursys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// HeaderUsername is the non-standard HTTP header in which to send the API username.
	HeaderUsername = "username"
	// HeaderPassword is the non-standard HTTP header in which to send the API password.
	HeaderPassword = "password"
)

//go:generate mockery --name Client
// Install mockery from https://github.com/vektra/mockery

// Client is the API for interacting with Nursys
// 1. Manage Nurse List: add, update, and remove nurses from an institution’s nurse list. [POST, GET]
// 2. Change Password: changing an institution’s API password. [POST]
// 3. Nurse Lookup: institutions can retrieve detailed license and discipline/final orders status information about their enrolled nurses. [POST, GET]
// 4. Notification Lookup: institutions can get changes to the license and discipline/final orders status of their enrolled licenses. [POST, GET]
// 5. Retrieve Documents: institutions can retrieve license and discipline/final orders related documents. [GET]
type Client interface {
	ManageNurseList(ctx context.Context, request ManageNurseListSubmitRequestMessage) (ManageNurseListSubmitResponseMessage, error)
	GetManageNurseListResult(ctx context.Context, txID string) (ManageNurseListRetrieveResponseMessage, error)
	ChangePassword(ctx context.Context, request ChangePasswordSubmitRequestMessage) (ChangePasswordSubmitResponseMessage, error)
	NurseLookup(ctx context.Context, request NurseLookupSubmitRequestMessage) (ManageNurseListSubmitResponseMessage, error)
	GetNurseLookupResult(ctx context.Context, txID string) (NurseLookupRetrieveResponseMessage, error)
	NotificationLookup(ctx context.Context, request NotificationLookupSubmitRequestMessage) (NotificationLookupSubmitResponseMessage, error)
	GetNotificationLookupResult(ctx context.Context, txID string) (NotificationLookupRetrieveResponseMessage, error)
	RetrieveDocuments(ctx context.Context, documentIDs []string) (RetrieveDocumentsRetrieveResponseMessage, error)
}

// Urban Airship HTTP API Client implementation
type nsHTTPClient struct {
	httpClient  *http.Client
	username    string
	pasword     string
	endpointURL string
}

// ClientOption are configuration functions that can be passed to New to configure the client.
type ClientOption func(c *nsHTTPClient)

// New creates a new client instance configured with the given options.
// For example:
//
//	conn := nursys.New("https://api.whatever/", "username", "password")
func New(baseURL, username, password string, options ...ClientOption) Client {
	client := nsHTTPClient{
		endpointURL: baseURL,
		username:    username,
		pasword:     password,
	}
	for _, opt := range options {
		opt(&client)
	}
	if client.httpClient == nil {
		client.httpClient = &http.Client{}
	}
	return &client
}

// WithHTTPClient overrides the http.Client instance used by the Airship Client.
// This is useful for unit tests of the client itself, but not much else.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *nsHTTPClient) {
		c.httpClient = httpClient
	}
}

// InvokeEndpoint invokes the airship API endpoint by sending <body> to <endpoint> using HTTP <method>.
// The response body is discarded unless an error status is returned.
func (cfg *nsHTTPClient) invokeEndpoint(ctx context.Context, method string, endpoint string, body interface{}, target interface{}) error {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, cfg.endpointURL+endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("username", cfg.username) // The non-standard HTTP header in which to send the API username.
	req.Header.Add("password", cfg.pasword)  // The non-standard HTTP header in which to send the API password.

	resp, err := cfg.httpClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			respBody, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("nursys: request returned %d: %s", resp.StatusCode, respBody)
		} else if target != nil {
			return json.NewDecoder(resp.Body).Decode(target)
		}
	}

	return err
}

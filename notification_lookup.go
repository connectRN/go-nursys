package nursys

import (
	"context"
	"net/url"
	"time"
)

// NotificationLookup is an asynchronous method for retrieving a list of status changes for licenses
// enrolled in an institution’s nurse list.
// Clients send a date range to the API using the Notification Lookup HTTP POST method. The server will
// respond with a TransactionId. After a period of time to allow for processing, the client calls the
// Notification Lookup HTTP GET method with that TransactionId to retrieve their results.
func (c *nsHTTPClient) NotificationLookup(ctx context.Context, request NotificationLookupSubmitRequestMessage) (NotificationLookupSubmitResponseMessage, error) {
	var response NotificationLookupSubmitResponseMessage
	err := c.invokeEndpoint(ctx, "POST", "/notificationlookup", request, &response)
	return response, err
}

// The NotificationLookupSubmitRequestMessage object is the input into the Notification Lookup HTTP POST method.
// Please note:
//   - The start date must be on or before the end date
//   - The start date must be on or before the current date
//   - The end date must be on or before the current date
type NotificationLookupSubmitRequestMessage struct {
	StartDate string `json:"StartDate"` // Required Start date of the date range in YYYY-MM- DD format
	EndDate   string `json:"EndDate"`   // Required End date of the date range in YYYY-MM-DD format
}

// SetStartDate formats the StartDate field with date
func (nlsrm *NotificationLookupSubmitRequestMessage) SetStartDate(date time.Time) {
	nlsrm.StartDate = date.Format(time.DateOnly)
}

// SetEndDate formats the EndDate field with date
func (nlsrm *NotificationLookupSubmitRequestMessage) SetEndDate(date time.Time) {
	nlsrm.EndDate = date.Format(time.DateOnly)
}

// The NotificationLookupSubmitResponseMessage models the response from the Notification Lookup HTTP POST method.
type NotificationLookupSubmitResponseMessage struct {
	Transaction `json:"Transaction"`
}

// **** **** **** **** **** **** **** ****

// GetNotificationLookupResult retrieves the results of a NotificationLookup transaction.
func (c *nsHTTPClient) GetNotificationLookupResult(ctx context.Context, txID string) (NotificationLookupRetrieveResponseMessage, error) {
	var response NotificationLookupRetrieveResponseMessage
	err := c.invokeEndpoint(ctx, "GET", "/notificationlookup?transactionId="+url.QueryEscape(txID), nil, &response)
	return response, err
}

// NotificationLookupRetrieveResponseMessage models the response from the Notification Lookup HTTP GET method
type NotificationLookupRetrieveResponseMessage struct {
	ProcessingCompleteFlag      bool                         `json:"ProcessingCompleteFlag"` // Indicates if the API has finished processing the request.
	Transaction                 Transaction                  `json:"Transaction"`
	NotificationLookupResponses []NotificationLookupResponse `json:"NotificationLookupResponses"`
}

// NotificationLookupResponse is an individual response to the Notification Lookup GET method.
type NotificationLookupResponse struct {
	NcsbnID                     string `json:"NcsbnId,omitempty"`           // Optional 10 NCSBN ID is the public, globally unique identifier for all nurses from participating boards of nursing.
	JurisdictionAbbreviation    string `json:"JurisdictionAbbreviation"`    // Required 4 State board of nursing. See appendix for a list of valid values.
	Jurisdiction                string `json:"Jurisdiction"`                // Required 50 State board of nursing. See appendix for a list of valid values.
	LicenseNumber               string `json:"LicenseNumber"`               // Required 15 License number.
	LicenseType                 string `json:"LicenseType"`                 // Required 4 License type. See appendix for a list of valid values.
	FirstName                   string `json:"FirstName"`                   // Required 50 Licensee first name.
	LastName                    string `json:"LastName"`                    // Required 50 Licensee last name.
	RecordId                    string `json:"RecordId"`                    // Optional 50 Client-provided identifier.
	NotificationDate            Time   `json:"NotificationDate"`            // Required Date the status change was reported.
	LicenseStatusChange         string `json:"LicenseStatusChange"`         // Optional 500 License status changes affecting the enrolled license.
	DisciplineStatusChange      string `json:"DisciplineStatusChange"`      // Optional 500 Discipline/final orders status changes affecting the enrolled license.
	DisciplineStatusChangeOther string `json:"DisciplineStatusChangeOther"` // Optional 500 Discipline/final orders status changes affecting licenses the nurse may hold that may not be enrolled.
}

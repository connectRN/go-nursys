package nursys

import (
	"context"
	"net/url"
)

// ManageNurseList is an asynchronous method for adding, updating, and removing batches of nurses
// from an institution’s nurse list.
//
// Clients sent batches of nurse information to the API using the Manage Nurse List HTTP POST method.
// The server will respond with a TransactionId. After a period of time to allow for processing, the
// client calls the Manage Nurse List HTTP GET method with that TransactionId to retrieve their results.
func (c *nsHTTPClient) ManageNurseList(ctx context.Context, request ManageNurseListSubmitRequestMessage) (ManageNurseListSubmitResponseMessage, error) {
	var response ManageNurseListSubmitResponseMessage
	err := c.invokeEndpoint(ctx, "POST", "/managenurselist", request, &response)
	return response, err
}

// The ManageNurseListSubmitRequestMessage object is the input into the Manage Nurse List HTTP POST method.
type ManageNurseListSubmitRequestMessage struct {
	ManageNurseListRequests []ManageNurseListRequest
}

// ManageNurseListRequest is a request entry passed to the Nurse Lookup POST method.
//
// Institutions must submit some combination of jurisdiction abbreviation, license type, license number,
// and NCSBN ID. Different combinations of licenses will be affected. This applies to license being added,
// updated, and removed
type ManageNurseListRequest struct {
	SubmissionActionCode         string      `json:"SubmissionActionCode"`                   // Required 1 Submission action code
	JurisdictionAbbreviation     string      `json:"JurisdictionAbbreviation,omitempty"`     // Optional 4 State board of nursing. Please see section 3.2.2 for matching rules.
	LicenseNumber                string      `json:"LicenseNumber,omitempty"`                // Optional 15 License number. Please see section 3.2.2 for matching rules.
	LicenseType                  string      `json:"LicenseType,omitempty"`                  // Optional 4 License type. Please see section 3.2.2 for matching rules.
	NcsbnID                      interface{} `json:"NcsbnId,omitempty"`                      // Optional 10 NCSBN ID is the public, globally unique identifier for all nurses from participating boards of nursing. Please see section 3.2.2 for matching rules. Unfortunately empty NcsbnID is coming as empty string "" type in ManageNurseListResponse
	Email                        string      `json:"Email,omitempty"`                        // Optional 50 E-mail address.
	Address1                     string      `json:"Address1"`                               // Required 50 Address line 1.
	Address2                     string      `json:"Address2,omitempty"`                     // Optional 50 Address line 2.
	City                         string      `json:"City"`                                   // Required 50 City.
	State                        string      `json:"State"`                                  // Required 2 State.
	Zip                          string      `json:"Zip"`                                    // Required 10 Zip code.
	LastFourSSN                  string      `json:"LastFourSSN"`                            // Required 4 Last four digits of social security number.
	BirthYear                    interface{} `json:"BirthYear"`                              // Required 4 Birth year. Unfortunately empty BirthYear is coming as empty string "" type in ManageNurseListResponse
	HospitalPracticeSetting      string      `json:"HospitalPracticeSetting,omitempty"`      // Required 2 Hospital practice setting.
	HospitalPracticeSettingOther string      `json:"HospitalPracticeSettingOther,omitempty"` // Optional 50 Hospital practice setting (other).
	NotificationsEnabled         string      `json:"NotificationsEnabled"`                   // Required 1 License and state licensure action notifications enabled.
	RemindersEnabled             string      `json:"RemindersEnabled"`                       // Required 1 License expiration reminders enabled.
	RecordID                     string      `json:"RecordId,omitempty"`                     // Optional 50 Client-provided identifier echoed back as part of the response.
	LocationList                 string      `json:"LocationList,omitempty"`                 // Optional 100 Pipe delimited list of location codes.
}

// The ManageNurseListSubmitResponseMessage models the response from the Manage Nurse List HTTP POST method.
type ManageNurseListSubmitResponseMessage struct {
	Transaction `json:"Transaction"`
}

// **** **** **** **** **** **** **** ****

// GetManageNurseListResult retrieves the results of a ManageNurseList transaction.
func (c *nsHTTPClient) GetManageNurseListResult(ctx context.Context, txID string) (ManageNurseListRetrieveResponseMessage, error) {
	var response ManageNurseListRetrieveResponseMessage
	err := c.invokeEndpoint(ctx, "GET", "/managenurselist?transactionId="+url.QueryEscape(txID), nil, &response)
	return response, err
}

// ManageNurseListRetrieveResponseMessage models the response from the Manage Nurse List HTTP GET method
type ManageNurseListRetrieveResponseMessage struct {
	ProcessingCompleteFlag   bool                     `json:"ProcessingCompleteFlag"` // Indicates if the API has finished processing the request.
	Transaction              Transaction              `json:"Transaction"`
	ManageNurseListResponses []MangeNurseListResponse `json:"ManageNurseListResponses"`
}

// MangeNurseListResponse is an individual response to the Manage Nurse List thing.
type MangeNurseListResponse struct {
	SuccessFlag            bool                   `json:"SuccessFlag"`            // True or False indicator if the Nurse was successfully added/updated/removed.
	Errors                 []TransactionError     `json:"Errors"`                 // A collection of validation errors that may have occurred during the processing of the Nurse.
	ManageNurseListRequest ManageNurseListRequest `json:"ManageNurseListRequest"` // Nurse sent via request.
}

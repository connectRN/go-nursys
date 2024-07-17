package nursys

import (
	"context"
	"net/url"
)

// NurseLookup is an asynchronous method for retrieving public license and discipline/final orders status
// information for batches of nurses enrolled in an institution’s nurse list. Institutions send batches of
// license information for their enrolled nurses and Nurse Lookup will return all publicly available
// information for all licenses for those nurses within Nursys.
func (c *nsHTTPClient) NurseLookup(ctx context.Context, request NurseLookupSubmitRequestMessage) (ManageNurseListSubmitResponseMessage, error) {
	var response ManageNurseListSubmitResponseMessage
	err := c.invokeEndpoint(ctx, "POST", "/nurselookup", request, &response)
	return response, err
}

// The NurseLookupSubmitRequestMessage object is the input into the Nurse Lookup HTTP POST method.
type NurseLookupSubmitRequestMessage struct {
	NurseLookupRequests []NurseLookupRequest
}

// NurseLookupRequest is a request entry passed to the Nurse Lookup POST method.
// You must specify at least a few of these fields. Valid combinations of elements:
//
//	Jurisdiction Abbreviation, License Type, License Number
//	Jurisdiction Abbreviation, License Type, NCSBN ID
//	Jurisdiction Abbreviation, NCSBN ID
//	License Type, NCSBN ID
//	NCSBN ID
type NurseLookupRequest struct {
	JurisdictionAbbreviation string `json:"JurisdictionAbbreviation,omitempty"` // Optional 4 State board of nursing. Please see the appendix for a list of valid values.
	LicenseNumber            string `json:"LicenseNumber,omitempty"`            // Optional 15 License number.
	LicenseType              string `json:"LicenseType,omitempty"`              // Optional 4 License type. Please see the appendix for a list of valid values.
	NcsbnID                  string `json:"NcsbnId,omitempty"`                  // Optional 10 NCSBN ID is the public, globally unique identifier for all nurses from participating boards of nursing.
	RecordID                 string `json:"RecordId,omitempty"`                 // Optional 50 Client-provided id
}

// The NurseLookupSubmitResponseMessage models the response from the Nurse Lookup HTTP POST method.
type NurseLookupSubmitResponseMessage struct {
	Transaction `json:"Transaction"`
}

// GetManageNurseListResult retrieves the results of a NurseLookup transaction.
func (c *nsHTTPClient) GetNurseLookupResult(ctx context.Context, txID string) (NurseLookupRetrieveResponseMessage, error) {
	var response NurseLookupRetrieveResponseMessage
	err := c.invokeEndpoint(ctx, "GET", "/nurselookup?transactionId="+url.QueryEscape(txID), nil, &response)
	return response, err
}

// ManageNurseListRetrieveResponseMessage models the response from the Nurse Lookup HTTP GET method
type NurseLookupRetrieveResponseMessage struct {
	ProcessingCompleteFlag bool                  `json:"ProcessingCompleteFlag"` // Indicates if the API has finished processing the request.
	Transaction            Transaction           `json:"Transaction"`
	NurseLookupResponses   []NurseLookupResponse `json:"NurseLookupResponses"`
}

// NurseLookupResponse is an individual response to the Nurse Lookup GET method.
type NurseLookupResponse struct {
	SuccessFlag                           bool                      `json:"SuccessFlag"`                                     // Required True or False indicator if successfully found in the institution’s nurse list.
	Errors                                []TransactionError        `json:"Errors,omitempty"`                                // Optional A collection of validation errors that may have occurred during the processing of the Nurse.
	NurseLookupRequest                    NurseLookupRequest        `json:"NurseLookupRequest"`                              // Required Nurse sent via request.
	FirstName                             string                    `json:"FirstName"`                                       // Required 50 Nurse first name
	LastName                              string                    `json:"LastName"`                                        // Required 50 Nurse last name
	NcsbnID                               string                    `json:"NcsbnId,omitempty"`                               // Optional 10 NCSBN ID is the public, globally unique identifier for all nurses from participating boards of nursing.
	Messages                              []string                  `json:"Messages,omitempty"`                              // Optional A collection of notification messages regarding the nurse. It is vital to review these messages as they may contain important information.
	NurseLookupLicenses                   []NurseLookupLicense      `json:"NurseLookupLicenses,omitempty"`                   // Optional A collection of licenses.
	NurseLookupRNAuthorizationsToPractice []AuthorizationToPractice `json:"NurseLookupRNAuthorizationsToPractice,omitempty"` // Optional A collection of RN authorization to practice information for each state.
	NurseLookupPNAuthorizationsToPractice []AuthorizationToPractice `json:"NurseLookupPNAuthorizationsToPractice,omitempty"` // Optional A collection of PN authorization to practice information for each state.
}

// NurseLookupLicense is an element of NurseLookupResponse
type NurseLookupLicense struct {
	LastName                     string                        `json:"LastName"`                     //  Required 50 Licensee last name.
	FirstName                    string                        `json:"FirstName"`                    //  Required 50 Licensee first name.
	LicenseType                  string                        `json:"LicenseType"`                  //  Required 4 License type. See appendix for a list of valid values.
	JurisdictionAbbreviation     string                        `json:"JurisdictionAbbreviation"`     //  Required 4 State board of nursing abbreviation. See appendix for a list of valid values.
	Jurisdiction                 string                        `json:"Jurisdiction"`                 //  Required 50 State board of nursing description. See appendix for a list of valid values.
	LicenseNumber                string                        `json:"LicenseNumber"`                //  Required 15 License number.
	Active                       string                        `json:"Active"`                       //  Optional 50 Active status for the license.
	LicenseStatus                string                        `json:"LicenseStatus"`                //  Optional 500 Current status of the license.
	LicenseOriginalDate          string                        `json:"LicenseOriginalDate"`          //  Optional Original issue date for the license.
	LicenseExpirationDate        string                        `json:"LicenseExpirationDate"`        //  Optional Expiration date for the license.
	CompactStatus                string                        `json:"CompactStatus"`                //  Optional 50 Nurse Licensure Compact (NLC) status of the license. Please visit nursys.com for more information about the NLC.
	Messages                     []string                      `json:"Messages"`                     //  Optional A collection of notification messages regarding the license. It is vital to review these messages as they may contain important license information.
	NurseLookupDisciplines       []NurseLookupDiscipline       `json:"NurseLookupDisciplines"`       //  Optional Collection of discipline information associated with this license.
	NurseLookupNotifications     []NurseLookupNotification     `json:"NurseLookupNotifications"`     //  Optional Collection of member board notifications regarding this license.
	NurseLookupAdvancedPractices []NurseLookupAdvancedPractice `json:"NurseLookupAdvancedPractices"` //  Optional Collection of focus, specialty, and related information for advanced practice (APRN) licenses.
}

// NurseLookupDiscipline is an element of NurseLookupLicense
type NurseLookupDiscipline struct {
	JurisdictionAbbreviation          string                      `json:"JurisdictionAbbreviation"`          // Required 4 Abbreviation for the state board of nursing that took the discipline/final orders. See appendix for a list of valid values.
	Jurisdiction                      string                      `json:"Jurisdiction"`                      // Required 50 Description for the state board of nursing that took the discipline/final orders. See appendix for a list of valid values.
	DateActionWasTaken                Time                        `json:"DateActionWasTaken"`                // Required Date the discipline/final orders was taken.
	AgainstPrivilegeToPracticeFlag    bool                        `json:"AgainstPrivilegeToPracticeFlag"`    // Required Flag to indicate if the discipline/final orders was taken against the license’s Privilege To Practice (PTP) as part of the Nurse Licensure Compact (NLC). Please visit nursys.com for more information about the NLC.
	NurseLookupBasisForActions        []NurseLookupBasisForAction `json:"NurseLookupBasisForActions"`        // Optional Collection of NPDB basis for action codes and descriptions for this discipline/final order. Note: some state boards of nursing elect not to provide this information if board order documents are attached to the discipline/final order.
	NurseLookupInitialActions         []NurseLookupAction         `json:"NurseLookupInitialActions"`         // Optional Collection of NPDB action codes and descriptions for this discipline/final order. Note: some state boards of nursing elect not to provide this information if board order documents are attached to the discipline/final order.
	NurseLookupInitialActionDocuments []NurseLookupDocument       `json:"NurseLookupInitialActionDocuments"` // Optional Collection of board order documents associated with this discipline/final order.
	NurseLookupRevisionReports        []NurseLookupRevisionReport `json:"NurseLookupRevisionReports"`        // Optional Collection of revisions
}

// NurseLookupBasisForAction is an element of NurseLookupDiscipline
type NurseLookupBasisForAction struct {
	BasisForActionCode        string `json:"BasisForActionCode"`        // Required 2 NPDB code for this discipline/final order basis for action.
	BasisForActionDescription string `json:"BasisForActionDescription"` // Required 150 NPDB description for this discipline/final order
}

// NurseLookupAction is an element of NurseLookupDiscipline or NurseLookupRevisionReport
type NurseLookupAction struct {
	ActionDate             Time   `json:"ActionDate"`             // Required Date of the discipline/final order action
	ActionCode             string `json:"ActionCode"`             // Required 5 NPDB code for this discipline/final order action.
	ActionDescription      string `json:"ActionDescription"`      // Required 150 NPDB description for this discipline/final order action.
	ActionStayedFlag       bool   `json:"ActionStayedFlag"`       // Required Flag to indicate if this action has been stayed by the state board of nursing.
	StartDate              Time   `json:"StartDate"`              // Optional Start date for this discipline/final order action.
	EndDate                Time   `json:"EndDate"`                // Optional End date for this discipline/final order action.
	Duration               string `json:"Duration"`               // Optional 50 Duration for this discipline/final order action (Indefinite/Unspecified, Permanent, or Specified)
	AutomaticReinstatement string `json:"AutomaticReinstatement"` // Optional 50 Indicates if the license is automatically reinstated upon the conclusion of the discipline/final order (No, Yes, or Yes With Conditions)
}

// NurseLookupDocument is an element of NurseLookupDiscipline or NurseLookupRevisionReport or NurseLookupNotification
type NurseLookupDocument struct {
	ActionDate   Time   `json:"ActionDate"`   // Optional Date of the discipline/final order action associated with this board order document.
	DocumentId   string `json:"DocumentId"`   // Required 50 Unique identifier for this board order document. This identifier can be used to retrieve the document using the Retrieve Document HTTP GET method described in section 3.6.
	DocumentName string `json:"DocumentName"` // Required 50 Name of the board order document.
}

// NurseLookupRevisionReport is an element of NurseLookupDiscipline
type NurseLookupRevisionReport struct {
	RevisionReportDate                 Time                  `json:"RevisionReportDate"`                 // Optional Date of the discipline/final order revision actions.
	NurseLookupRevisionActions         []NurseLookupAction   `json:"NurseLookupRevisionActions"`         // Optional Collection of NPDB action codes and descriptions for the revisions to this discipline/final order. Note: some state boards of nursing elect not to provide this information if board order documents are attached to the discipline/final order.
	NurseLookupRevisionActionDocuments []NurseLookupDocument `json:"NurseLookupRevisionActionDocuments"` // Optional Collection of board order documents associated with these discipline/final order revision actions.
}

// NurseLookupNotification is an elemement of NurseLookupLicense
type NurseLookupNotification struct {
	JurisdictionAbbreviation string                `json:"JurisdictionAbbreviation"` //  Required 4 The state board of nursing that placed the member board notification on the license. See appendix for a list of valid values.
	Jurisdiction             string                `json:"Jurisdiction"`             //  Required 50 The state board of nursing that placed the member board notification on the license. See appendix for a list of valid values.
	NotificationDate         Time                  `json:"NotificationDate"`         //  Required The date the member board notification was placed on the license.
	NotificationMessage      string                `json:"NotificationMessage"`      //  Required 5000 The message text for the member board notification.
	NotificationDocuments    []NurseLookupDocument `json:"NotificationDocuments"`    //  Optional Collection of documents associated with this member board notification.
}

// NurseLookupAdvancedPractice is an element of NurseLookupLicense
type NurseLookupAdvancedPractice struct {
	FocusSpecialty               string `json:"FocusSpecialty"`               // Required 50 Description of the advanced practice (APRN) focus or specialty
	PrescriptionAuthority        string `json:"PrescriptionAuthority"`        // Required 50 Indicates if this advanced practice (APRN) focus or specialty has prescription authority.
	CertificationExpirationDate  Time   `json:"CertificationExpirationDate"`  // Optional Expiration date for the advanced practice (APRN) certification associated with this focus or specialty.
	FocusSpecialtyExpirationDate Time   `json:"FocusSpecialtyExpirationDate"` // Optional Expiration date for this focus or specialty
}

// AuthorizationToPractice is an element of NurseLookupResponse
type AuthorizationToPractice struct {
	StateAbbreviation                  string `json:"StateAbbreviation"`                  // Required 2 State. See appendix for a list of valid values.
	StateDescription                   string `json:"StateDescription"`                   // Required 50 State. See appendix for a list of valid values.
	AuthorizationToPracticeCode        string `json:"AuthorizationToPracticeCode"`        // Required 1 Code indicating the authorization to practice in the given state. See appendix for a list of valid values.
	AuthorizationToPracticeDescription string `json:"AuthorizationToPracticeDescription"` // Required 50 Description indicating the authorization to practice in the given state. See appendix for a list of valid values.
	AuthorizationToPracticeNarrative   string `json:"AuthorizationToPracticeNarrative"`   // Required 5000 Detailed information regarding the authorization to practice in the given state.
}

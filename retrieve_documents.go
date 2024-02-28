package nursys

import (
	"context"
	"net/url"
	"strings"
)

// Documents attached to discipline/final orders and member board notifications can be retrieved by
// calling the Retrieve Documents HTTP GET method. DocumentId values are returned as part of the Nurse
// Lookup HTTP GET method response. These DocumentIds are passed to the Retrieve Documents HTTP
// GET method to get the actual document itself. A maximum of five DocumentId values can be included in
// any one method call.
func (c *nsHTTPClient) RetrieveDocuments(ctx context.Context, documentIDs []string) (RetrieveDocumentsRetrieveResponseMessage, error) {
	var response RetrieveDocumentsRetrieveResponseMessage
	ids := strings.Join(documentIDs, ",") // Nursys expects DocumentId values separated by commas
	err := c.invokeEndpoint(ctx, "GET", "/retrievedocuments?documentIds="+url.QueryEscape(ids), nil, &response)
	return response, err
}

// RetrieveDocumentsRetrieveResponseMessage models the response of the Retrieve Documents HTTP GET method.
type RetrieveDocumentsRetrieveResponseMessage struct {
	Transaction `json:"Transaction"`
	Documents   []RetrieveDocumentResponse `json:"Documents"`
}

// RetrieveDocumentResponse models a document returned from the  Retrieve Documents HTTP GET method.
type RetrieveDocumentResponse struct {
	SuccessFlag      bool   `json:"SuccessFlag"`      // Required True or False indicator if the document was successfully found
	DocumentId       string `json:"DocumentId"`       // Required 50 Unique identifier for the document as supplied in the request.
	DocumentName     string `json:"DocumentName"`     // Required 50 Name of the document, including the extension.
	DocumentContents string `json:"DocumentContents"` // Required Contents of the binary file.  (base64?? the spec doesn't say)
}

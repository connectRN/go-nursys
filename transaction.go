package nursys

// The Transaction JSON object contains information about the API request and response.
type Transaction struct {
	TransactionID          string             `json:"TransactionId"`          // Required Unique transaction identifier. Used as an input for the Manage Nurse List HTTP GET method.
	TransactionDate        Time               `json:"TransactionDate"`        // Required System date and time when the request was processed.
	TransactionComment     string             `json:"TransactionComment"`     // Optional System provided comments regarding the processing of the request.
	TransactionSuccessFlag bool               `json:"TransactionSuccessFlag"` // Required True or False indicator if the request was successfully processed or not.
	TransactionErrors      []TransactionError `json:"TransactionErrors"`      // Optional A collection of system errors that may have occurred during the processing of the request.
}

// TransactionError models an error in the API response JSON.
type TransactionError struct {
	ErrorID      int64  `json:"ErrorID"`      // System assigned processing error identifier.
	ErrorMessage string `json:"ErrorMessage"` // System provided processing error message
}

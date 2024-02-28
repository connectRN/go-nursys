package nursys

import "context"

// ChangePassword changes an institution’s API password.
func (c *nsHTTPClient) ChangePassword(ctx context.Context, request ChangePasswordSubmitRequestMessage) (ChangePasswordSubmitResponseMessage, error) {
	var response ChangePasswordSubmitResponseMessage
	err := c.invokeEndpoint(ctx, "POST", "/changepassword", request, &response)
	return response, err
}

// The ChangePasswordSubmitRequestMessage object is the input into the Change Password HTTP POST method.
type ChangePasswordSubmitRequestMessage struct {
	NewPassword string `json:"NewPassword"`
}

// The ChangePasswordSubmitResponseMessage models the response from the Change Password HTTP POST method.
type ChangePasswordSubmitResponseMessage struct {
	Transaction `json:"Transaction"`
}

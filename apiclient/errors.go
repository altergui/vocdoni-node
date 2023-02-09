package apiclient

import "fmt"

var (
	// ErrAccountNotConfigured is returned when the client has not been configured with an account.
	ErrAccountNotConfigured = fmt.Errorf("account not configured")

	ErrUnmarshalFailed = fmt.Errorf("unmarshal failed")
)

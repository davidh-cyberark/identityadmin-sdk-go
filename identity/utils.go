package identity

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

func headersToString(header http.Header) string {
	var sb strings.Builder
	for key, values := range header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", key, strings.Join(values, ",")))
	}
	return sb.String()
}

func ValidateRoleName(s string) error {
	// A username can contain of any UTF-8 alphanumeric characters plus the
	// symbols + (plus), - (dash), _ (underscore), and . (period).
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && s != "+" && s != "-" && s != "." && s != "_" {
			return fmt.Errorf("role name contains bad chars")
		}
	}
	return nil
}

// ReturnErrorWhenBodySuccessIsFalse returns an error if the body's Success field is false (even though HTTP status code is 200 OK).
func ReturnErrorWhenBodySuccessIsFalse(body *RolesStoreRole) error {
	if body == nil {
		return fmt.Errorf("failed to create role: body is nil")
	}
	success := body.Success
	if *success {
		return nil
	}
	Message := body.Message
	if Message != nil {
		return fmt.Errorf("failed to create role: %s", *Message)
	}
	MessageID := body.MessageID
	if MessageID != nil {
		return fmt.Errorf("failed to create role: %s", *MessageID)
	}
	ErrorCode := body.ErrorCode
	if ErrorCode != nil {
		return fmt.Errorf("failed to create role (no message available) error code: %s", *ErrorCode)
	}
	ErrorID := body.ErrorID
	if ErrorID != nil {
		return fmt.Errorf("failed to create role (no message available) error id: %s", *ErrorID)
	}
	return fmt.Errorf("failed to create role: unknown error")
}

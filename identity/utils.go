package identity

import (
	"encoding/json"
	"fmt"
	"unicode"
)

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

// IsBodySuccessFieldFalse returns an error if the body's Success field is false (even though HTTP status code is 200 OK).
func IsBodySuccessFieldFalse(body *RolesStoreRole) error {
	if body == nil {
		return fmt.Errorf("unknown error: body is nil")
	}
	success := body.Success
	if *success {
		return nil
	}

	retMsg := make(map[string]string)
	Message := body.Message
	if Message != nil {
		retMsg["Message"] = *Message
	}
	MessageID := body.MessageID
	if MessageID != nil {
		retMsg["MessageID"] = *MessageID
	}
	ErrorCode := body.ErrorCode
	if ErrorCode != nil {
		retMsg["ErrorCode"] = *ErrorCode
	}
	ErrorID := body.ErrorID
	if ErrorID != nil {
		retMsg["ErrorID"] = *ErrorID
	}

	if len(retMsg) == 0 {
		return fmt.Errorf("unknown error")
	}

	jsonBytes, err := json.Marshal(retMsg)
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %v", err)
	}
	jsonString := string(jsonBytes)
	return fmt.Errorf("%s", jsonString)
}

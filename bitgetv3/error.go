package bitgetv3

import (
	"fmt"
)

// Error - Bitget API error: the code/msg pair of the response envelope
// Codes are strings; "00000" means success
// https://www.bitget.com/api-doc/uta/error-code/restapi
type Error struct {
	// Code - Bitget error code (string, "00000" on success)
	Code string
	// Text - human-readable error message (msg field of the envelope)
	Text string
}

func (o *Error) Std() error {
	if o.Empty() {
		return nil
	}
	return o
}

func (o *Error) Empty() bool {
	return o.Code == "00000"
}

func (o *Error) Error() string {
	return fmt.Sprintf("code[%s]: %s", o.Code, o.Text)
}

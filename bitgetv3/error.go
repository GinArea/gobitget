package bitgetv3

import (
	"fmt"
)

type Error struct {
	Code string
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

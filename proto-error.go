package gogram

import (
	"github.com/thewug/gogram/data"

	"fmt"
)

type ProtocolError struct {
	HTTPCode    int
	Description string
}

func (this ProtocolError) Error() (string) {
	return fmt.Sprintf("Failure indicated by API endpoint (%d: %s)\n", this.HTTPCode, this.Description)
}

func NewError(code int, description string) (*ProtocolError) {
	return &ProtocolError{
		HTTPCode: code,
		Description: description,
	}
}

func HandleSoftError(resp *data.TGenericResponse) (error) {
	if resp.Ok != true {
		return NewError(*resp.Error_code, *resp.Description)
	}
	return nil
}

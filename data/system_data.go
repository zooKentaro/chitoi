package data

import (
	"github.com/pkg/errors"
)

// SystemLoggerRequest is XXX
type SystemLoggerRequest struct {
	Tag  string `json:"tag,string"`
	Body string `json:"body,string"`
}

func (u *SystemLoggerRequest) Validate() error {
	if l := len(u.Tag); l > 30 {
		return errors.Errorf("tag is too long, current length: %d, max: %d", l, 30)
	}
	if bl := len(u.Body); bl > 5000 {
		return errors.Errorf("body is too long, current length: %d, max: %d", bl, 5000)
	}
	return nil
}

// SystemLoggerResponse is XXX
type SystemLoggerResponse struct{}

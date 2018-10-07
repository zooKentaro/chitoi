package data

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/constant"
)

// UserSignupRequest is XXX
type UserSignupRequest struct {
    Platform uint32 `json:"platform"`
}

func (u *UserSignupRequest) Validate() error {
    if u.Platform != constant.PlatformTypePC && u.Platform != constant.PlatformTypeIOS && u.Platform != constant.PlatformTypeAndroid {
        return errors.Errorf("error platform type %d is not found", u.Platform)
    }
    return nil
}

// UserSignupResponse is XXX
type UserSignupResponse struct {
    Token string `json:"token"`
}

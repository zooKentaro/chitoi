package data

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/constant"
    "github.com/uenoryo/chitoi/database/row"
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
    User      *row.User `json:"user"`
    SessionID string    `json:"session_id"`
}

// UserLoginRequest is XXX
type UserLoginRequest struct {
    Token string `json:"token"`
}

// UserLoginResponse is XXX
type UserLoginResponse struct {
    User      *row.User `json:"user"`
    SessionID string    `json:"session_id"`
}

// UserInfoRequest is XXX
type UserInfoRequest struct {
    SessionID string `json:"session_id"`
}

// UserInfoResponse is XXX
type UserInfoResponse struct {
    User *row.User `json:"user"`
}

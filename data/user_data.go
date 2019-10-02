package data

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/constant"
    "github.com/uenoryo/chitoi/database/row"
)

// UserSignupRequest ...
type UserSignupRequest struct {
    Platform uint32 `json:"platform"`
}

func (u *UserSignupRequest) Validate() error {
    if u.Platform != constant.PlatformTypePC && u.Platform != constant.PlatformTypeIOS && u.Platform != constant.PlatformTypeAndroid {
        return errors.Errorf("error platform type %d is not found", u.Platform)
    }
    return nil
}

// UserSignupResponse ...
type UserSignupResponse struct {
    User       *row.User       `json:"user"`
    SessionID  string          `json:"session_id"`
    Businesses []*row.Business `json:"businesses"`
    UserRanks  []*row.UserRank `json:"user_ranks"`
}

// UserLoginRequest ...
type UserLoginRequest struct {
    Token string `json:"token"`
}

// UserLoginResponse ...
type UserLoginResponse struct {
    User              *row.User           `json:"user"`
    SessionID         string              `json:"session_id"`
    UserBusinesses    []*row.UserBusiness `json:"user_businesses"`
    Businesses        []*row.Business     `json:"businesses"`
    UserRanks         []*row.UserRank     `json:"user_ranks"`
    IsTodayFirstLogin bool                `json:"is_today_first_login"`
}

// UserInfoRequest ...
type UserInfoRequest struct {
    SessionID string `json:"session_id"`
}

// UserInfoResponse ...
type UserInfoResponse struct {
    User           *row.User           `json:"user"`
    UserBusinesses []*row.UserBusiness `json:"user_businesses"`
}

// UserRecordRequest ...
type UserRecordRequest struct {
    SessionID      string `json:"session_id"`
    BestScore      uint64 `json:"best_score,string"`
    BestTotalScore uint64 `json:"best_total_score,string"`
}

// UserRecordResponse ...
type UserRecordResponse struct{}

// UserEditNameRequest ...
type UserEditNameRequest struct {
    SessionID string `json:"session_id"`
    Name      string `json:"name"`
}

// UserEditNameResponse ...
type UserEditNameResponse struct {
    User *row.User `json:"user"`
}

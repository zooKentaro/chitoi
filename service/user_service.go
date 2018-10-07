package service

import "github.com/uenoryo/chitoi/data"

type UserService interface {
    Signup() (*data.UserSignupResponse, error)
}

type userService struct{}

// NewUserService is XXX
func NewUserService() UserService {
    return &userService{}
}

// Signup is XXX
func (u *userService) Signup() (*data.UserSignupResponse, error) {
    return &data.UserSignupResponse{}, nil
}

package service

type UserService interface {
    Signup()
}

type userService struct{}

// NewUserService is XXX
func NewUserService() UserService {
    return &userService{}
}

// Signup is XXX
func (s *userService) Signup() {
    return
}

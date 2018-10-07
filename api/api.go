package api

import (
    "net/http"

    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/handler"
    "github.com/uenoryo/chitoi/service"
)

// NewServer is XXX
func NewServer() (*http.ServeMux, error) {
    server := http.NewServeMux()

    core, err := core.New()
    if err != nil {
        return nil, errors.Wrap(err, "error new core")
    }

    userService := service.NewUserService(core)
    userHandler := handler.NewUserServer(handler.NewUserHandler(userService))

    server.Handle("/user/", userHandler)

    return server, nil
}

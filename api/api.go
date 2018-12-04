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
    if err := core.LoadMasterdata(); err != nil {
        return nil, errors.Wrap(err, "error load masterdata")
    }

    userService := service.NewUserService(core)
    userHandler := handler.NewUserServer(handler.NewUserHandler(core, userService))

    gameService := service.NewGameService(core)
    gameHandler := handler.NewGameServer(handler.NewGameHandler(core, gameService))

    businessService := service.NewBusinessService(core)
    businessHandler := handler.NewBusinessServer(handler.NewBusinessHandler(core, businessService))

    server.Handle("/user/", userHandler)
    server.Handle("/game/", gameHandler)
    server.Handle("/business/", businessHandler)

    return server, nil
}

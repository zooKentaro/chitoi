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
    userHandler := handler.NewUserServer(handler.NewUserHandler(core, userService))

    gameService := service.NewGameService(core)
    gameHandler := handler.NewGameServer(handler.NewGameHandler(core, gameService))

    businessService := service.NewBusinessService(core)
    businessHandler := handler.NewBusinessServer(handler.NewBusinessHandler(core, businessService))

    roomService := service.NewRoomService(core)
    roomHandler := handler.NewRoomServer(handler.NewRoomHandler(core, roomService))

    systemHandler := handler.NewSystemServer(handler.NewSystemHandler(core))

    server.Handle("/user/", userHandler)
    server.Handle("/game/", gameHandler)
    server.Handle("/business/", businessHandler)
    server.Handle("/room/", roomHandler)
    server.Handle("/system/", systemHandler)

    return server, nil
}

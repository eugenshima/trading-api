package main

import (
	"github.com/eugenshima/trading-api/internal/handlers"
	"github.com/eugenshima/trading-api/internal/repository"
	"github.com/eugenshima/trading-api/internal/service"
	proto "github.com/eugenshima/trading-api/proto"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func main() {
	e := echo.New()
	conn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	client := proto.NewPriceServiceClient(conn)
	rps := repository.NewProfileRepository(client)
	srv := service.NewProfileService(rps)
	handler := handlers.NewProfileApiHandler(srv)

	api := e.Group("/profile")
	{
		api.POST("/login", handler.Login)
		api.POST("/signup", handler.SignUp)
	}

	e.Logger.Fatal(e.Start(":8089"))
}

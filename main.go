// Package main is the entry-point for the package
package main

import (
	"fmt"

	balanceProto "github.com/eugenshima/balance/proto"
	priceServiceProto "github.com/eugenshima/price-service/proto"
	profileProto "github.com/eugenshima/profile/proto"
	"github.com/eugenshima/trading-api/internal/handlers"
	"github.com/eugenshima/trading-api/internal/middleware"
	"github.com/eugenshima/trading-api/internal/repository"
	"github.com/eugenshima/trading-api/internal/service"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

// main is the main function
// nolint:gocritic, staticcheck
func main() {
	e := echo.New()

	profileConn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		return
	}

	priceServiceConn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		return
	}

	balanceConn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		return
	}

	defer func() {
		err = profileConn.Close()
		if err != nil {
			fmt.Println("Error closing profile connection")
		}
	}()
	defer func() {
		err = priceServiceConn.Close()
		if err != nil {
			fmt.Println("Error closing price-service connection")
		}
	}()
	defer func() {
		err = balanceConn.Close()
		if err != nil {
			fmt.Println("Error closing balance connection")
		}
	}()

	profileClient := profileProto.NewProfilesClient(profileConn)
	profileRps := repository.NewProfileRepository(profileClient)
	profileSrv := service.NewProfileService(profileRps)
	handler := handlers.NewProfileAPIHandler(profileSrv)

	priceServiceClient := priceServiceProto.NewPriceServiceClient(priceServiceConn)
	priceServiceRps := repository.NewPriceServiceRepository(priceServiceClient)

	balanceClient := balanceProto.NewBalanceServiceClient(balanceConn)
	balanceRps := repository.NewBalanceRepository(balanceClient)
	balanceSrv := service.NewBalanceService(balanceRps, priceServiceRps)
	balanceHandler := handlers.NewBalanceAPIHandler(balanceSrv)

	middlewr := middleware.UserIdentity()

	auth := e.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/signup", handler.SignUp)
		auth.POST("/refreshtokenpair", handler.RefreshTokenPair)
		auth.DELETE("/deleteprofile", handler.DeleteProfile)
	}

	balance := e.Group("/balance")
	{
		balance.POST("/deposit", balanceHandler.GetLatestPrice, middlewr)
		balance.POST("/getBalance", balanceHandler.GetBalance, middlewr)
		balance.POST("/withdraw", balanceHandler.Withdraw, middlewr)
	}
	// in progress...
	/*
		trading := e.Group("/trading")
		{
			trading.POST("/openPosition", tradingHandler.OpenPosition, middlewr)
			trading.POST("/closePosition", tradingHandler.ClosePosition, middlewr)
		}
	*/

	e.Logger.Fatal(e.Start(":8089"))
}

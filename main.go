package main

import (
	"github.com/eugenshima/trading-api/internal/handlers"
	"github.com/eugenshima/trading-api/internal/repository"
	"github.com/eugenshima/trading-api/internal/service"
	balanceProto "github.com/eugenshima/trading-api/proto/balance"
	priceProto "github.com/eugenshima/trading-api/proto/price-service"
	profileProto "github.com/eugenshima/trading-api/proto/profile"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func main() {
	e := echo.New()

	profileConn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer profileConn.Close()

	priceServiceConn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer priceServiceConn.Close()

	balanceConn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer balanceConn.Close()

	profileClient := profileProto.NewProfilesClient(profileConn)
	profileRps := repository.NewProfileRepository(profileClient)
	profileSrv := service.NewProfileService(profileRps)
	handler := handlers.NewProfileApiHandler(profileSrv)

	priceServiceClient := priceProto.NewPriceServiceClient(priceServiceConn)
	priceServiceRps := repository.NewPriceServiceRepository(priceServiceClient)

	balanceClient := balanceProto.NewBalanceServiceClient(balanceConn)
	balanceRps := repository.NewBalanceRepository(balanceClient)
	balanceSrv := service.NewBalanceService(balanceRps, priceServiceRps)
	balanceHandler := handlers.NewBalanceApiHandler(balanceSrv)

	auth := e.Group("/auth")
	{
		auth.POST("/login", handler.Login)   //without jwt logic
		auth.POST("/signup", handler.SignUp) //without jwt logic

		//auth.POST("/refreshtokenpair", handler.RefreshToken)
		//auth.DELETE("/deleteprofile", handler.Delete)
	}

	balance := e.Group("/balance")
	{
		balance.POST("/deposit", balanceHandler.Deposit)       //with jwt logic
		balance.POST("/getBalance", balanceHandler.GetBalance) //with jwt logic
		balance.POST("/withdraw", balanceHandler.Withdraw)     //with jwt logic

		//balance.POST("/GetLatest", balanceHandler.GetLatestPrice)
	}

	e.Logger.Fatal(e.Start(":8089"))
}

// Package repository contains methods to communicate with postgres and gRPC servers
package repository

import (
	"context"
	"fmt"

	priceServiceProto "github.com/eugenshima/price-service/proto"
	"github.com/eugenshima/trading-api/internal/model"
)

// priceServiceRepository strucct ....
type priceServiceRepo struct {
	client priceServiceProto.PriceServiceClient
}

// NewPriceServiceRepository creates a new priceServiceRepository
func NewPriceServiceRepository(client priceServiceProto.PriceServiceClient) *priceServiceRepo {
	return &priceServiceRepo{client: client}
}

// RecvShares receives a list of selected shares
func (r *priceServiceRepo) RecvShares(ctx context.Context, selectedShares []string) (*model.Shares, error) {
	req := &priceServiceProto.SubscribeRequest{
		ShareName: selectedShares,
	}
	stream, err := r.client.Subscribe(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}
	response, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("recv: %w", err)
	}
	shares := &model.Shares{
		ShareName:  response.Shares[0].ShareName,
		SharePrice: response.Shares[0].SharePrice,
	}
	return shares, nil
}

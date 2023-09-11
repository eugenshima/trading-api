package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/trading-api/internal/model"
	proto "github.com/eugenshima/trading-api/proto/price-service"
)

type priceServiceRepository struct {
	client proto.PriceServiceClient
}

func NewPriceServiceRepository(client proto.PriceServiceClient) *priceServiceRepository {
	return &priceServiceRepository{client: client}
}

func (r *priceServiceRepository) RecvShares(ctx context.Context, selectedShares []string) (*model.Shares, error) {
	req := &proto.SubscribeRequest{
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

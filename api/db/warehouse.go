package db

import "context"

type WarehouseClient struct{}

func NewWarehouseClient() *WarehouseClient {
	return &WarehouseClient{}
}

func (w *WarehouseClient) GetRevenue(ctx context.Context, startDate, endDate, accountID string) string {
	// TODO: Replace with actual warehouse query
	return "0"
}

func (w *WarehouseClient) GetConversionRate(ctx context.Context, startDate, endDate, accountID string) string {
	// TODO: Replace with actual warehouse query
	return "0"
}

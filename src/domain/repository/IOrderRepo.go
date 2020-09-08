package repository

import (
	"context"

	"microservice_gokit_base/src/domain/model"
)

// IOrderRepository decribes the repository of Order
type IOrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) (string, error)
	GetOrderByID(ctx context.Context, id string) (model.Order, error)
	ChangeOrderStatus(ctx context.Context, id string, status string) (int64, error)
	GetAll(ctx context.Context) ([]*model.Order, error)
	GetPage(ctx context.Context, page int64, size int64) ([]*model.Order, error)
	Count(ctx context.Context) (int64, error)
}

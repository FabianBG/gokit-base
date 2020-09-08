package service

import (
	"context"

	"microservice_gokit_base/src/domain/model"
	"microservice_gokit_base/src/domain/repository"
	"microservice_gokit_base/src/domain/utils"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/validator.v2"
)

// IOrderService describes the Order service.
type IOrderService interface {
	Create(ctx context.Context, order model.Order) (string, error)
	GetByID(ctx context.Context, id string) (model.Order, error)
	GetAll(ctx context.Context) ([]*model.Order, error)
	GetPage(ctx context.Context, page int64, size int64) ([]*model.Order, error)
	ChangeStatus(ctx context.Context, id string, status string) (int64, error)
	Count(ctx context.Context) (int64, error)
}

// OrderService instance
type OrderService struct {
	repository repository.IOrderRepository
	uuid       utils.IUUIDGenerator
	date       utils.IDateGenerator
	logger     log.Logger
}

// NewOrderService creates and returns a new Order service instance
func NewOrderService(rep repository.IOrderRepository, uuid utils.IUUIDGenerator,
	date utils.IDateGenerator, logger log.Logger) IOrderService {
	return &OrderService{
		repository: rep,
		uuid:       uuid,
		date:       date,
		logger:     logger,
	}
}

// Create makes an order
func (s *OrderService) Create(ctx context.Context, order model.Order) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	id := s.uuid.GenerateID()
	order.ID = id
	order.Status = "Pending"
	order.CreatedOn = s.date.NowTimestamp()
	created, err := s.repository.CreateOrder(ctx, order)
	if err := validator.Validate(order); err != nil {
		level.Debug(logger).Log("err", err)
		return "", err
	}
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	return created, nil
}

// GetByID returns an order given by id
func (s *OrderService) GetByID(ctx context.Context, id string) (model.Order, error) {
	logger := log.With(s.logger, "method", "GetByID")
	order, err := s.repository.GetOrderByID(ctx, id)
	if err != nil {
		level.Debug(logger).Log("msg", err)
		return order, err
	}
	return order, nil
}

// ChangeStatus changes the status of an order
func (s *OrderService) ChangeStatus(ctx context.Context, id string, status string) (int64, error) {
	logger := log.With(s.logger, "method", "ChangeStatus")
	changed, err := s.repository.ChangeOrderStatus(ctx, id, status)
	if err != nil && changed < 1 {
		level.Error(logger).Log("err", err)
		return 0, err
	}
	return changed, nil
}

// GetAll recive all orders
func (s *OrderService) GetAll(ctx context.Context) ([]*model.Order, error) {
	logger := log.With(s.logger, "method", "GetAll")
	orders, err := s.repository.GetAll(ctx)
	if err != nil {
		level.Debug(logger).Log("msg", err)
		return nil, err
	}
	return orders, nil
}

// GetPage returns paged orders
func (s *OrderService) GetPage(ctx context.Context, page int64, size int64) ([]*model.Order, error) {
	logger := log.With(s.logger, "method", "GetPage")
	orders, err := s.repository.GetPage(ctx, page, size)
	if err != nil {
		level.Debug(logger).Log("msg", err)
		return nil, err
	}
	return orders, nil
}

// Count returns the coutn of documents
func (s *OrderService) Count(ctx context.Context) (int64, error) {
	logger := log.With(s.logger, "method", "Count")
	counted, err := s.repository.Count(ctx)
	if err != nil {
		level.Error(logger).Log("msg", err)
		return -1, err
	}
	return counted, nil
}

package repository

import (
	"context"
	"errors"
	"strconv"

	"microservice_gokit_base/src/domain/model"
	domainRepo "microservice_gokit_base/src/domain/repository"

	"github.com/go-kit/kit/log"
)

var (
	ErrRepository = errors.New("unable to handle request")
)

type DbMemory struct {
	Data []model.Order
}

type repositoryMem struct {
	db     *DbMemory
	logger log.Logger
}

// NewOrderRepositoryMem returns a concrete repository backe by mem array
func NewOrderRepositoryMem(db *DbMemory, logger log.Logger) (domainRepo.IOrderRepository, error) {
	// return  repository
	return &repositoryMem{
		db:     db,
		logger: log.With(logger, "rep", "cockroachdb"),
	}, nil
}

// CreateOrder inserts a new order and its order items into db
func (repo *repositoryMem) CreateOrder(ctx context.Context, order model.Order) (string, error) {

	repo.db.Data = append(repo.db.Data, order)
	return strconv.Itoa(len(repo.db.Data)), nil
}

// ChangeOrderStatus changes the order status
func (repo *repositoryMem) ChangeOrderStatus(ctx context.Context, id string, status string) (int64, error) {
	for i := 0; i <= len(repo.db.Data); i++ {
		if repo.db.Data[i].ID == id {
			repo.db.Data[i].Status = status
			return 1, nil
		}
	}
	return 0, nil
}

// GetOrderByID query the order by given id
func (repo *repositoryMem) GetOrderByID(ctx context.Context, id string) (model.Order, error) {
	var orderRow = model.Order{}
	for i := 0; i <= len(repo.db.Data); i++ {
		if repo.db.Data[i].ID == id {
			return repo.db.Data[i], nil
		}
	}
	return orderRow, nil
}

// GetAll query all orders
func (repo *repositoryMem) GetAll(ctx context.Context) ([]*model.Order, error) {
	var results []*model.Order
	for _, elem := range repo.db.Data {
		results = append(results, &elem)
	}
	return results, nil
}

// GetPage query orders by a page
func (repo *repositoryMem) GetPage(ctx context.Context, page int64, size int64) ([]*model.Order, error) {
	var results []*model.Order
	for _, elem := range repo.db.Data {
		results = append(results, &elem)
	}
	return results, nil
}

// Count get the count of documents
func (repo *repositoryMem) Count(ctx context.Context) (int64, error) {
	return int64(len(repo.db.Data)), nil
}

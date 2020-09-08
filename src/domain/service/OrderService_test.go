package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/golang/mock/gomock"
	"gotest.tools/assert"

	"microservice_gokit_base/src/domain/model"
	"microservice_gokit_base/src/infraestructure/repository"
	"microservice_gokit_base/src/mocks"
)

func TestDomainOrderService(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "order",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var (
		orderRepository = mocks.NewMockIOrderRepository(mockCtrl)
		uuidGen         = mocks.NewMockIUUIDGenerator(mockCtrl)
		dateGen         = mocks.NewMockIDateGenerator(mockCtrl)
		orderService    = &OrderService{
			repository: orderRepository,
			uuid:       uuidGen,
			date:       dateGen,
			logger:     logger,
		}
		ctx   = context.TODO()
		order = model.Order{
			ID:           "1",
			Status:       "Pending",
			RestaurantID: "EL MAGIO",
			OrderItems: []model.OrderItem{
				{},
			},
		}
		orderEmpty    = model.Order{}
		numberOfItems = int64(2)
		errorCount    = int64(-1)
		page          = int64(1)
		size          = int64(50)
		listOrders    = []*model.Order{&order, &order}
		mockError     = errors.New("errors")
	)

	t.Run("orderService.Create",
		func(t *testing.T) {
			t.Run("WHEN everything is ok SHOULD return an no null id",
				func(t *testing.T) {
					gomock.InOrder(
						uuidGen.EXPECT().GenerateID().Return(order.ID).Times(1),
						dateGen.EXPECT().NowTimestamp().Return(int64(0)).Times(1),
						orderRepository.EXPECT().CreateOrder(
							ctx,
							//gomock.AssignableToTypeOf(order)).Return(order.ID, nil).Times(1)
							order).Return(order.ID, nil).Times(1),
					)

					id, err := orderService.Create(ctx, order)
					assert.NilError(t, err)
					assert.Assert(t, id != "")
				})
			t.Run("WHEN an error happend SHOULD return an error",
				func(t *testing.T) {
					orderRepository.EXPECT().CreateOrder(
						ctx,
						gomock.AssignableToTypeOf(order)).Return("", errors.New("error on the mongo repository")).Times(1)
					uuidGen.EXPECT().GenerateID().Return(order.ID).Times(1)
					dateGen.EXPECT().NowTimestamp().Return(int64(0)).Times(1)
					id, err := orderService.Create(ctx, order)
					assert.Assert(t, id == "")
					assert.ErrorContains(t, err, "error on the mongo repository")
				})
		})

	t.Run("orderService.GetPage",
		func(t *testing.T) {
			t.Run("WHEN everything is ok SHOULD return an array of objects",
				func(t *testing.T) {
					orderRepository.EXPECT().GetPage(
						ctx,
						page,
						size).Return(listOrders, nil).Times(1)

					list, err := orderService.GetPage(ctx, page, size)
					assert.NilError(t, err)
					assert.Assert(t, len(list) == 2)
				})
			t.Run("WHEN an error happend on the repository SHOULD return an error",
				func(t *testing.T) {
					orderRepository.EXPECT().GetPage(
						ctx,
						page,
						size).Return(nil, repository.ErrMongoRepository).Times(1)

					list, err := orderService.GetPage(ctx, page, size)
					assert.Assert(t, err == repository.ErrMongoRepository)
					assert.Assert(t, list == nil)
				})
			t.Run("WHEN an error happend on the cursor reading SHOULD retrrun an error",
				func(t *testing.T) {
					orderRepository.EXPECT().GetPage(
						ctx,
						page,
						size).Return(nil, repository.ErrMongoRepository).Times(1)

					list, err := orderService.GetPage(ctx, page, size)
					assert.Assert(t, err == repository.ErrMongoRepository)
					assert.Assert(t, list == nil)
				})
		})

	t.Run("orderService.Count",
		func(t *testing.T) {
			t.Run("WHEN everything is ok SHOULD return an no null id",
				func(t *testing.T) {
					gomock.InOrder(
						orderRepository.EXPECT().Count(
							ctx).Return(numberOfItems, nil).Times(1),
					)

					count, err := orderService.Count(ctx)
					assert.NilError(t, err)
					assert.Assert(t, count == numberOfItems)
				})
			t.Run("WHEN an error happend SHOULD return an error",
				func(t *testing.T) {
					orderRepository.EXPECT().Count(
						ctx).Return(errorCount, repository.ErrMongoRepository).Times(1)

					count, err := orderService.Count(ctx)
					assert.Assert(t, err == repository.ErrMongoRepository)
					assert.Assert(t, count == errorCount)
				})
		})

	t.Run("orderService.GetByID",
		func(t *testing.T) {

			t.Run("WHEN everything is ok SHOULD return an objects",
				func(t *testing.T) {
					gomock.InOrder(
						orderRepository.EXPECT().GetOrderByID(
							ctx,
							order.ID).Return(order, nil).Times(1),
					)

					cli, err := orderService.GetByID(ctx, order.ID)
					assert.NilError(t, err)
					assert.Assert(t, cli.ID == order.ID)
				})

			t.Run("WHEN an error happend SHOULD return an error",
				func(t *testing.T) {
					gomock.InOrder(
						orderRepository.EXPECT().GetOrderByID(
							ctx,
							order.ID).Return(orderEmpty, mockError).Times(1),
					)

					_, err := orderService.GetByID(ctx, order.ID)
					assert.Assert(t, err == mockError)
				})
		})

	t.Run("orderService.GetAll",
		func(t *testing.T) {

			t.Run("WHEN everything is ok SHOULD return an objects",
				func(t *testing.T) {
					gomock.InOrder(
						orderRepository.EXPECT().GetAll(
							ctx).Return(listOrders, nil).Times(1),
					)

					orders, err := orderService.GetAll(ctx)
					assert.NilError(t, err)
					assert.Assert(t, len(orders) == len(listOrders))
				})

			t.Run("WHEN ana error happend SHOULD an error",
				func(t *testing.T) {
					gomock.InOrder(
						orderRepository.EXPECT().GetAll(
							ctx).Return(nil, mockError).Times(1),
					)

					_, err := orderService.GetAll(ctx)
					assert.Assert(t, err == mockError)
				})
		})

	t.Run("orderService.ChangeStatus",
		func(t *testing.T) {
			partialUpdate := "disabled"

			t.Run("WHEN everything is ok SHOULD return a number of correct changes status",
				func(t *testing.T) {
					gomock.InOrder(
						orderRepository.EXPECT().ChangeOrderStatus(
							ctx,
							order.ID,
							partialUpdate).Return(int64(1), nil).Times(1),
					)
					statusCount, err := orderService.ChangeStatus(ctx, order.ID, partialUpdate)
					assert.NilError(t, err)
					assert.Assert(t, statusCount == 1)
				})
			t.Run("WHEN an error happend on the repository SHOULD return an error",
				func(t *testing.T) {
					gomock.InOrder(
						orderRepository.EXPECT().ChangeOrderStatus(
							ctx,
							order.ID,
							partialUpdate).Return(int64(0), repository.ErrMongoRepository).Times(1),
					)

					statusCount, err := orderService.ChangeStatus(ctx, order.ID, partialUpdate)
					assert.Assert(t, err == repository.ErrMongoRepository)
					assert.Assert(t, statusCount == 0)
				})
		})
}

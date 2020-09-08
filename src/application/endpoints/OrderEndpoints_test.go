package endpoints

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"

	"microservice_gokit_base/src/domain/model"
	"microservice_gokit_base/src/mocks"
)

func TestApplicationOrderEndpoints(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var (
		orderServiceDomain = mocks.NewMockIOrderService(mockCtrl)
		ctx                = context.TODO()
		orderEndpoints     = MakeOrderEndpoints(orderServiceDomain)
		order              = model.Order{
			ID: "001",
		}
		listOrder = []*model.Order{
			&order,
		}
		mockError = errors.New("errors")
		page      = int64(10)
		size      = int64(10)
		count     = int64(10)
	)

	t.Run("orderEndpoints.Create",
		func(t *testing.T) {
			t.Run("WHEN everything is ok SHOULD return an correct response",
				func(t *testing.T) {

					gomock.InOrder(
						orderServiceDomain.EXPECT().Create(
							ctx,
							order).Return(order.ID, nil).Times(1),
					)

					req := CreateRequest{
						Order: order,
					}

					ok, err := orderEndpoints.CreateEndpoint()(ctx, req)
					assert.NilError(t, err)
					assert.Assert(t, ok.(CreateResponse).ID == order.ID)
				})

			t.Run("WHEN an error happend SHOULD return an error response",
				func(t *testing.T) {

					gomock.InOrder(
						orderServiceDomain.EXPECT().Create(
							ctx,
							order).Return("", mockError).Times(1),
					)

					req := CreateRequest{
						Order: order,
					}

					ok, err := orderEndpoints.CreateEndpoint()(ctx, req)
					assert.NilError(t, err)
					assert.Assert(t, ok.(CreateResponse).Err == mockError)

				})
		})

	t.Run("orderEndpoints.GetByIDEndpoint",
		func(t *testing.T) {
			t.Run("WHEN everything is ok SHOULD return an correct response",
				func(t *testing.T) {

					gomock.InOrder(
						orderServiceDomain.EXPECT().GetByID(
							ctx,
							order.ID).Return(order, nil).Times(1),
					)

					req := GetByIDRequest{
						ID: order.ID,
					}

					ok, err := orderEndpoints.GetByIDEndpoint()(ctx, req)
					assert.NilError(t, err)
					assert.Assert(t, ok.(GetByIDResponse).Order.ID == order.ID)
				})

			t.Run("WHEN an error happend SHOULD return an error response",
				func(t *testing.T) {

					gomock.InOrder(
						orderServiceDomain.EXPECT().GetByID(
							ctx,
							order.ID).Return(order, mockError).Times(1),
					)

					req := GetByIDRequest{
						ID: order.ID,
					}

					ok, err := orderEndpoints.GetByIDEndpoint()(ctx, req)
					assert.NilError(t, err)
					assert.Assert(t, ok.(GetByIDResponse).Err == mockError)
				})
		})

	t.Run("orderEndpoints.GetAllEndpoint",
		func(t *testing.T) {
			t.Run("WHEN everything is ok and have page and size SHOULD return an correct response of paged data",
				func(t *testing.T) {

					gomock.InOrder(
						orderServiceDomain.EXPECT().GetPage(
							ctx, page, size).Return(listOrder, nil).Times(1),
					)

					req := GetAllRequest{
						Page: page,
						Size: size,
					}

					ok, err := orderEndpoints.GetAllEndpoint()(ctx, req)
					assert.NilError(t, err)
					assert.Assert(t, len(ok.(GetlAllResponse).Orders) == len(listOrder))
				})

			t.Run("WHEN everything is ok and have no page and size SHOULD return an correct response of all data",
				func(t *testing.T) {

					gomock.InOrder(
						orderServiceDomain.EXPECT().GetAll(
							ctx).Return(listOrder, nil).Times(1),
					)

					req := GetAllRequest{}

					ok, err := orderEndpoints.GetAllEndpoint()(ctx, req)
					assert.NilError(t, err)
					assert.Assert(t, len(ok.(GetlAllResponse).Orders) == len(listOrder))
				})

			t.Run("WHEN an error happend SHOULD return an error response",
				func(t *testing.T) {

					gomock.InOrder(
						orderServiceDomain.EXPECT().GetAll(
							ctx).Return(nil, mockError).Times(1),
					)

					req := GetAllRequest{}

					ok, err := orderEndpoints.GetAllEndpoint()(ctx, req)

					assert.NilError(t, err)
					assert.Assert(t, ok.(GetlAllResponse).Err == mockError)
				})
		})

	t.Run("orderEndpoints.CountEndpoint",
		func(t *testing.T) {
			t.Run("WHEN everything is ok SHOULD return an correct response",
				func(t *testing.T) {

					gomock.InOrder(
						orderServiceDomain.EXPECT().Count(
							ctx).Return(count, nil).Times(1),
					)

					req := CountRequest{}

					ok, err := orderEndpoints.CountEndpoint()(ctx, req)
					assert.NilError(t, err)
					assert.Assert(t, ok.(CountResponse).Count == count)
				})

			t.Run("WHEN an error happend SHOULD return an error response",
				func(t *testing.T) {

					gomock.InOrder(
						orderServiceDomain.EXPECT().Count(
							ctx).Return(int64(0), mockError).Times(1),
					)

					req := CountRequest{}
					res := CountResponse{
						Count: 0,
						Err:   mockError,
					}

					ok, err := orderEndpoints.CountEndpoint()(ctx, req)
					assert.NilError(t, err)
					assert.Assert(t, ok == res)
				})
		})
}

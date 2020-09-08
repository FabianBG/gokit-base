package endpoints

import (
	"context"
	"os"

	"microservice_gokit_base/src/domain/model"
	"microservice_gokit_base/src/domain/service"

	"github.com/go-kit/kit/endpoint"
)

// IOrderEndpoints holds all Go kit endpoints for the Order service.
type IOrderEndpoints interface {
	CreateEndpoint() endpoint.Endpoint
	GetByIDEndpoint() endpoint.Endpoint
	GetAllEndpoint() endpoint.Endpoint
	ChangeStatusEndpoint() endpoint.Endpoint
	CountEndpoint() endpoint.Endpoint
	SaludoEndpoint() endpoint.Endpoint
}

// OrderEndpoints Struct to instanciate endpoints
type OrderEndpoints struct {
	orderDomainService service.IOrderService
}

// MakeOrderEndpoints initializes all Go kit endpoints for the Order service.
func MakeOrderEndpoints(s service.IOrderService) IOrderEndpoints {
	return &OrderEndpoints{
		orderDomainService: s,
	}
}

// CreateRequest holds the request parameters for the Create method.
type CreateRequest struct {
	Order model.Order
}

// CreateResponse holds the response values for the Create method.
type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

// Failed implements endpoint.Failer.
func (r CreateResponse) Failed() error { return r.Err }

// CreateEndpoint Service to expose domain logic
func (s *OrderEndpoints) CreateEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.orderDomainService.Create(ctx, req.Order)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

// GetByIDRequest holds the request parameters for the GetByID method.
type GetByIDRequest struct {
	ID string
}

// GetByIDResponse holds the response values for the GetByID method.
type GetByIDResponse struct {
	Order model.Order `json:"result"`
	Err   error       `json:"error,omitempty"`
}

// Failed implements endpoint.Failer.
func (r GetByIDResponse) Failed() error { return r.Err }

//GetByIDEndpoint Service to expose domain logic
func (s *OrderEndpoints) GetByIDEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		orderRes, err := s.orderDomainService.GetByID(ctx, req.ID)
		return GetByIDResponse{Order: orderRes, Err: err}, nil
	}
}

// GetAllRequest holds the request parameters for the GetAll method.
type GetAllRequest struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
}

// GetlAllResponse holds the response values for the GetAll method.
type GetlAllResponse struct {
	Orders []*model.Order `json:"result"`
	Err    error          `json:"error,omitempty"`
}

// Failed implements endpoint.Failer.
func (r GetlAllResponse) Failed() error { return r.Err }

// GetAllEndpoint Service to expose domain logic
func (s *OrderEndpoints) GetAllEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		var orders []*model.Order
		req := request.(GetAllRequest)
		if req.Size > 0 {
			orders, err = s.orderDomainService.GetPage(ctx, req.Page, req.Size)
		} else {
			orders, err = s.orderDomainService.GetAll(ctx)
		}
		if orders == nil {
			orders = make([]*model.Order, 0)
		}
		return GetlAllResponse{Orders: orders, Err: err}, nil
	}
}

// CountRequest holds the request parameters for the GetPage method.
type CountRequest struct {
}

// CountResponse holds the response values for the GetPage method.
type CountResponse struct {
	Count int64 `json:"result"`
	Err   error `json:"error,omitempty"`
}

// Failed implements endpoint.Failer.
func (r CountResponse) Failed() error { return r.Err }

// CountEndpoint Service to expose domain logic
func (s *OrderEndpoints) CountEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		count, err := s.orderDomainService.Count(ctx)
		return CountResponse{Count: count, Err: err}, nil
	}
}

// ChangeStatusRequest holds the request parameters for the ChangeStatus method.
type ChangeStatusRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// ChangeStatusResponse holds the response values for the ChangeStatus method.
type ChangeStatusResponse struct {
	Updated int64 `json:"updated"`
	Err     error `json:"error,omitempty"`
}

// Failed implements endpoint.Failer.
func (r ChangeStatusResponse) Failed() error { return r.Err }

// ChangeStatusEndpoint Service to expose domain logic
func (s *OrderEndpoints) ChangeStatusEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeStatusRequest)
		changed, err := s.orderDomainService.ChangeStatus(ctx, req.ID, req.Status)
		return ChangeStatusResponse{Updated: changed, Err: err}, nil
	}
}

// SaludoEndpoint funcion de pruebas
func (s *OrderEndpoints) SaludoEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		name, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		return "Hostname ( " + name + " ) saludos desde quito.", nil
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = GetlAllResponse{}
	_ endpoint.Failer = GetByIDResponse{}
	_ endpoint.Failer = ChangeStatusResponse{}
	_ endpoint.Failer = CreateResponse{}
	_ endpoint.Failer = CountResponse{}
)

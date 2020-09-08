package repository

import (
	"context"
	"errors"
	"os"

	"github.com/go-kit/kit/log/level"

	"go.mongodb.org/mongo-driver/bson"

	"microservice_gokit_base/config"
	"microservice_gokit_base/src/domain/model"
	domainRepo "microservice_gokit_base/src/domain/repository"

	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// ErrMongoRepository general error from mongo repo
	ErrMongoRepository = errors.New("error on the mongo repository")
	// ErrConnectionMongoRepository when the conection fails
	ErrConnectionMongoRepository = errors.New("error connecting on the mongo database")
	// ErrNotFoundMongoRepository when the query returns no results
	ErrNotFoundMongoRepository = errors.New("error no results found")
)

const (
	collection  = "order"
	idField     = "_id"
	statusField = "status"
)

type repositoryMongo struct {
	collection *mongo.Collection
	logger     log.Logger
}

// GetConnectionMongo geenrate a conection to mongo
func GetConnectionMongo(ctx context.Context, logger log.Logger) *mongo.Collection {
	config := config.Instance()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		level.Error(logger).Log("err", err)
	} else {
		level.Info(logger).Log("msg", "mongo connection OK")
	}
	return client.Database(config.MongoDB).Collection(collection)
}

// NewOrderMongoRepository returns a concrete repository backe by mem array
func NewOrderMongoRepository(collection *mongo.Collection, logger log.Logger) (domainRepo.IOrderRepository, error) {
	return &repositoryMongo{
		collection: collection,
		logger:     log.With(logger, "rep", "mongo"),
	}, nil
}

// CreateOrder inserts a new order and its order items into db
func (repo *repositoryMongo) CreateOrder(ctx context.Context, order model.Order) (string, error) {

	insertResult, err := repo.collection.InsertOne(ctx, order)
	if err != nil {
		level.Error(repo.logger).Log("err", err)
	}
	return insertResult.InsertedID.(string), nil
}

// ChangeOrderStatus changes the order status
func (repo *repositoryMongo) ChangeOrderStatus(ctx context.Context, orderID string, status string) (int64, error) {
	filter := bson.D{bson.E{Key: idField, Value: orderID}}
	update := bson.D{
		{"$set", bson.D{
			{statusField, status},
		}},
		{"$currentDate", bson.D{
			{"lastModified", true},
		}},
	}
	updateResult, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		level.Error(repo.logger).Log("err", err)
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

// GetOrderByID query the order by given id
func (repo *repositoryMongo) GetOrderByID(ctx context.Context, id string) (model.Order, error) {
	filter := bson.D{bson.E{Key: idField, Value: id}}
	options := options.FindOne()
	findResult := repo.collection.FindOne(ctx, filter, options)
	var result model.Order
	err := findResult.Decode(&result)
	if err != nil {
		level.Debug(repo.logger).Log("msg", err)
		return result, ErrNotFoundMongoRepository
	}
	return result, err
}

// GetAll query all orders
func (repo *repositoryMongo) GetAll(ctx context.Context) ([]*model.Order, error) {

	var results []*model.Order
	findOptions := options.Find()
	cur, err := repo.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		level.Error(repo.logger).Log("err", err)
		return nil, ErrMongoRepository
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result model.Order
		err := cur.Decode(&result)
		if err != nil {
			level.Debug(repo.logger).Log("msg", err)
		}
		results = append(results, &result)
	}
	if err := cur.Err(); err != nil {
		level.Debug(repo.logger).Log("msg", err)
		return nil, ErrNotFoundMongoRepository
	}
	return results, nil
}

// GetPage query orders by a page
func (repo *repositoryMongo) GetPage(ctx context.Context, page int64, size int64) ([]*model.Order, error) {

	var results []*model.Order
	findOptions := options.Find()
	findOptions.SetLimit(size)
	findOptions.SetSkip(page * size)
	cur, err := repo.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		level.Error(repo.logger).Log("err", err)
		return nil, ErrMongoRepository
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result model.Order
		err := cur.Decode(&result)
		if err != nil {
			level.Debug(repo.logger).Log("msg", err)
		}
		results = append(results, &result)
	}
	if err := cur.Err(); err != nil {
		level.Debug(repo.logger).Log("msg", err)
		return nil, ErrNotFoundMongoRepository
	}
	return results, nil
}

// Count get the count of documents
func (repo *repositoryMongo) Count(ctx context.Context) (int64, error) {
	countOptions := options.Count()
	counted, err := repo.collection.CountDocuments(ctx, bson.D{}, countOptions)
	if err != nil {
		level.Error(repo.logger).Log("err", err)
		return -1, ErrMongoRepository
	}
	return counted, nil
}

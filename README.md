# MICROSERVICE PROJECT TEMPLATE

Project base for a simple microservice on golang, based on go kit, DDD.

## Project structure

```
src/
    services/
    domain/
    infraestructure/
```

## Usage

* Load envs coping .env.example to .env and customize the params, then `source .nev`

Run:

`go run main.go`

Build:

`go build`

## Testing
For testing the testing files are named by the prefix _test:

`go test microservice_gokit_base/src/...`

Testing coverage:

`go test ./... -cover` 

`go test ./... -coverprofile cover.out; go tool cover -func cover.out`

Mocking interfaces:


`mockgen -destination=src/mocks/mock_IUUIDGenerator.go -package=mocks microservice_gokit_base/src/domain/utils IUUIDGenerator`

`mockgen -destination=src/mocks/mock_IDateGenerator.go -package=mocks microservice_gokit_base/src/domain/utils IDateGenerator`

`mockgen -destination=src/mocks/mock_IOrderRepository.go -package=mocks microservice_gokit_base/src/domain/repository IOrderRepository`

## Validator

Refer to [https://godoc.org/gopkg.in/validator.v2]

## Environment

Copy env.example file and add to enviroment to run the app


## More info

* GO KIT - [HOME](http://gokit.io)
* DDD - [INFO](https://apiumhub.com/tech-blog-barcelona/introduction-domain-driven-design)
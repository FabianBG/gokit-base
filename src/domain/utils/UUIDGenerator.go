package utils

import "github.com/gofrs/uuid"

// IUUIDGenerator describes the uuid generator util.
type IUUIDGenerator interface {
	GenerateID() string
}

// UUIDGenerator describes implementation of generator
type UUIDGenerator struct {
}

// NewUUIDGenerator function to instanciate the generator
func NewUUIDGenerator() IUUIDGenerator {
	return &UUIDGenerator{}
}

// GenerateID function to generate an uuid
func (g *UUIDGenerator) GenerateID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

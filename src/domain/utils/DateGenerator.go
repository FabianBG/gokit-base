package utils

import "time"

// IDateGenerator describes the date generator util.
type IDateGenerator interface {
	NowTimestamp() int64
}

// DateGenerator describes implementation of generator
type DateGenerator struct {
}

// NewDateGenerator function to instanciate the generator
func NewDateGenerator() IDateGenerator {
	return &DateGenerator{}
}

// NowTimestamp function to get current timestamp
func (g *DateGenerator) NowTimestamp() int64 {
	return time.Now().Unix()
}

package repository

import (
	"microservice_gokit_base/src/domain/model"

	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gotest.tools/assert"
)

func TestOrderRepository(t *testing.T) {
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

	// Create DB
	var dbMemory = &DbMemory{
		Data: []model.Order{},
	}

	t.Run("NewOrderRepositoryMem",
		func(t *testing.T) {
			repo, err := NewOrderRepositoryMem(dbMemory, logger)
			assert.NilError(t, err)
			assert.Assert(t, repo != nil)
		})
}

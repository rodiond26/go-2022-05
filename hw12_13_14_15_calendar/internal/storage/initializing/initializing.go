package initializng

import (
	"context"

	storage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage"
	memoryStorage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/memory"
	sqlStorage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/sql"
)

func NewStorage(ctx context.Context, typo string, dsn string) (s storage.Storage, err error) {
	switch typo {
	case "sql":
		storageSQL := sqlStorage.New()
		pool, err := storageSQL.Connect(ctx, dsn)
		if err != nil {
			return nil, err
		}
		storageSQL.PgxPool = pool
		s = storageSQL
	case "memory":
	default:
		s = memoryStorage.New()
	}
	return s, nil
}

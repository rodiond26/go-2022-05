package init_storage

import (
	storage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/memory"
	sqlStorage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/sql"
)

func NewStorage( /*ctx context.Context,*/ typo string, dsn string) (storage.Storage, error) {
	var s storage.Storage
	switch typo {
	case "sql":
		storageSQL := sqlStorage.New()
		// pool, err := sqlStorage.Connect(ctx, dsn)
		// if err != nil {
		// 	return nil, err
		// }
		// storageSQL.PgxPool = pool
		s = storageSQL
	default:
		storageMemory := memorystorage.New()
		s = storageMemory
	}
	return s, nil
}

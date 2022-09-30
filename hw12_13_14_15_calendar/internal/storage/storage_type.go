package storage

// import (
// 	"context"

// 	memorystorage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/memory"
// 	sqlstorage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/sql"
// 	"golang.org/x/mod/sumdb/storage"
// )

// const (
// 	sql    = "sql"
// 	memory = "memory"
// )

// func StorageType(ctx context.Context, typo string, dsn string) (s storage.Storage, err error) {
// 	switch typo {
// 	case sql:
// 		storageSQL := sqlstorage.New()
// 		pool, err := sqlstorage.Connect(ctx, dsn)
// 		if err != nil {
// 			return nil, err
// 		}
// 		storageSQL.Pool = pool
// 		db = storageSQL
// 	case memory:
// 		storageSQL := sqlstorage.New()
// 		pool, err := sqlstorage.Connect(ctx, dsn)
// 		if err != nil {
// 			return nil, err
// 		}
// 		storageSQL.Pool = pool
// 		db = storageSQL
// 	default:
// 		storageMemory := memorystorage.New()
// 		db = storageMemory
// 	}
// 	return db, nil
// }

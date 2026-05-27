package store

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	TestStoreMethod() error
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{
		db: db,
	}
}

func (s *store) TestStoreMethod() error {
	fmt.Println("Hello from store!")
	return nil
}

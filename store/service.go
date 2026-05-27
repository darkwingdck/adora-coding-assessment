package store

import (
	"database/sql"
	"fmt"
)

type Store interface {
	TestStoreMethod() error
}

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &store{
		db: db,
	}
}

func (s *store) TestStoreMethod() error {
	fmt.Println("Hello from store!")
	return nil
}

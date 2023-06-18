package dao

import (
	"context"
	"database/sql"

	"gitea.slauson.io/slausonio/go-utils/siodao"
)

var ctx = context.Background()

type PostDao struct {
	db *sql.DB
}

func NewPostDao() *PostDao {
	return &PostDao{
		db: siodao.DatabaseConnection(),
	}
}

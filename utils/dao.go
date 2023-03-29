package utils

import (
	"gitea.slauson.io/slausonio/go-libs/dao"
	"github.com/uptrace/bun"
)

func InitDB() *bun.DB {
	return dao.DatabaseConnection()
}

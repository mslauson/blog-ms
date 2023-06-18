package dao

import (
	"context"
	"database/sql"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/siodao"
)

var ctx = context.Background()

type Post struct {
	ID           siogeneric.SioNullInt64
	Title        siogeneric.SioNullString
	Body         siogeneric.SioNullString
	PostedDate   siogeneric.SioNullTime
	UpdatedDate  siogeneric.SioNullTime
	DeletionDate siogeneric.SioNullTime
	SoftDeleted  siogeneric.SioNullBool
	CreatedByID  siogeneric.SioNullInt64
	UpdatedByID  siogeneric.SioNullInt64
}

type Comment struct {
	ID           siogeneric.SioNullInt64
	Content      siogeneric.SioNullString
	CommentDate  siogeneric.SioNullTime
	SoftDeleted  siogeneric.SioNullBool
	DeletionDate siogeneric.SioNullTime
	PostID       siogeneric.SioNullInt64
	UserID       siogeneric.SioNullInt64
}



type PostDao struct {
	db *sql.DB
}

func NewPostDao() *PostDao {
	return &PostDao{
		db: siodao.DatabaseConnection(),
	}
}

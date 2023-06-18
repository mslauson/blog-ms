package dao

import (
	"context"
	"database/sql"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/siodao"
)

var ctx = context.Background()

type PostDao struct {
	db *sql.DB
}

type PostDaoInterface interface {
	CreatePost(post *siogeneric.BlogPost) error
	ExistsPost(post *siogeneric.BlogPost) (bool, error)
	ExistsByIDPost(ID int64) (bool, error)
	GetAllPosts() (*[]siogeneric.BlogPost, error)
	UpdatePost(post *siogeneric.BlogPost) error
	AddCommentToPost(post *siogeneric.BlogPost, comment *siogeneric.BlogComment) error
	UpdateCommentInPost(post *siogeneric.BlogPost, comment *siogeneric.BlogComment) error
	SoftDeletePost(post *siogeneric.BlogPost) error
	SoftDeleteCommentOnPost(post *siogeneric.BlogPost, comment *siogeneric.BlogComment) error
}

func NewPostDao() *PostDao {
	return &PostDao{
		db: siodao.DatabaseConnection(),
	}
}

package dao

import (
	"context"
	"database/sql"

	"entgo.io/ent/entc/integration/edgefield/ent/post"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/siodao"
)

var ctx = context.Background()

type PDao struct {
	db *sql.DB
}

type PostDao interface {
	CreatePost(post *siogeneric.BlogPost) error
	Exists(post *siogeneric.BlogPost) (bool, error)
	ExistsByID(ID int64) (bool, error)
	GetAllPosts() (*[]siogeneric.BlogPost, error)
	UpdatePost(post *siogeneric.BlogPost) error
	AddComment(post *siogeneric.BlogPost, comment *siogeneric.BlogComment) error
	UpdateComment(post *siogeneric.BlogPost, comment *siogeneric.BlogComment) error
	SoftDeletePost(post *siogeneric.BlogPost) error
	SoftDeleteComment(post *siogeneric.BlogPost, comment *siogeneric.BlogComment) error
}

func NewPostDao() *PDao {
	return &PDao{
		db: siodao.DatabaseConnection(),
	}
}

func (pd *PDao) CreatePost(post *siogeneric.BlogPost) error {
	sql := `INSERT INTO blog_posts (title, body, posted_date, created_by_id, soft_deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := pd.db.QueryRowContext(
		ctx,
		sql,
		post.Title,
		post.Body,
		post.PostedDate,
		post.CreatedByID,
		false,
	).Scan(&post.ID)
	return err
}

func (pd *PDao) Exists(post *siogeneric.BlogPost) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM blog_posts WHERE title = $1 created_by_id = $2)`
	var exists bool
	err := pd.db.QueryRowContext(ctx, sql, post.Title, post.CreatedByID).
		Scan(&exists)
	return exists, err
}

func (pd *PDao) ExistsByID(ID int64) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM blog_posts WHERE id = $1)`
	var exists bool
	err := pd.db.QueryRowContext(ctx, sql, post.ID).Scan(&exists)
	return exists, err
}

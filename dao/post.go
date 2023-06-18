package dao

import (
	"context"
	"database/sql"

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
	err := pd.db.QueryRowContext(ctx, sql, ID).Scan(&exists)
	return exists, err
}

func (pd *PDao) GetAllPosts() (*[]siogeneric.BlogPost, error) {
	sql := `SELECT 
			p.id, p.title, p.body, p.posted_date, p.updated_date, p.deletion_date, p.soft_deleted,
			c.id, c.content, c.comment_date, c.soft_deleted, c.deletion_date
		FROM blog.post p
		LEFT JOIN blog.comment c ON p.id = c.post_id`

	rows, err := pd.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []siogeneric.BlogPost
	for rows.Next() {
		var post siogeneric.BlogPost
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Body,
			&post.PostedDate,
			&post.CreatedByID,
			&post.SoftDeleted,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func (pd *PDao) scanPost(row *sql.Row) (*siogeneric.BlogPost, error) {
	post := new(siogeneric.BlogPost)
	if err := row.Scan(&post.ID, &post.Title, &post.Body, &post.PostedDate,
		&post.UpdatedDate, &post.DeletionDate, &post.SoftDeleted); err != nil {
		return nil, err
	}
}

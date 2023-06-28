package dao

import (
	"context"
	"database/sql"

	"gitea.slauson.io/slausonio/go-types/sioblog"
	"gitea.slauson.io/slausonio/go-utils/siodao"
)

var ctx = context.Background()

type BDao struct {
	db *sql.DB
}

//go:generate mockery --name BlogDao
type BlogDao interface {
	CreatePost(post *sioblog.BlogPost) error
	PostExists(title string, createdByID int64) (bool, error)
	PostExistsByID(ID int64) (bool, error)
	CommentExistsByID(ID int64) (bool, error)
	GetPostByID(ID int64) (*sioblog.BlogPost, error)
	GetCommentByID(ID int64) (*sioblog.BlogComment, error)
	GetAllPosts() (*[]*sioblog.BlogPost, error)
	GetAllCommentsByPostID(postID int64) (*[]*sioblog.BlogComment, error)
	UpdatePost(post *sioblog.BlogPost) error
	AddComment(comment *sioblog.BlogComment) error
	UpdateComment(comment *sioblog.BlogComment) error
	SoftDeletePost(post *sioblog.BlogPost) error
	SoftDeleteComment(comment *sioblog.BlogComment) error
}

func NewBlogDao() *BDao {
	return &BDao{
		db: siodao.DatabaseConnection(),
	}
}

func (bd *BDao) CreatePost(post *sioblog.BlogPost) error {
	query := `INSERT INTO post (title, body, posted_date, created_by_id, soft_deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := bd.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Body,
		post.PostedDate,
		post.CreatedByID,
		false,
	).Scan(&post.ID)
	return err
}

func (bd *BDao) PostExists(title string, createdByID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM post WHERE title = $1 AND created_by_id = $2 AND soft_deleted = false)`
	var exists bool
	err := bd.db.QueryRowContext(ctx, query, title, createdByID).
		Scan(&exists)
	return exists, err
}

func (bd *BDao) PostExistsByID(ID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM post WHERE id = $1 AND soft_deleted = false)`
	var exists bool
	err := bd.db.QueryRowContext(ctx, query, ID).Scan(&exists)
	return exists, err
}

func (bd *BDao) CommentExistsByID(ID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM comment WHERE id = $1)`
	var exists bool
	err := bd.db.QueryRowContext(ctx, query, ID).Scan(&exists)
	return exists, err
}

func (bd *BDao) GetPostByID(ID int64) (*sioblog.BlogPost, error) {
	query := `SELECT id, title, body, created_by_id, updated_by_id, posted_date, updated_date, deletion_date, soft_deleted FROM post WHERE id = $1 AND soft_deleted = false`
	rows, err := bd.db.QueryContext(ctx, query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		post, err := bd.scanPost(rows)
		if err != nil {
			return nil, err
		}

		comments, err := bd.GetAllCommentsByPostID(ID)
		if err != nil {
			return nil, err
		}

		post.Comments = comments
		return post, nil
	}

	return nil, sql.ErrNoRows
}

func (bd *BDao) GetCommentByID(ID int64) (*sioblog.BlogComment, error) {
	query := `SELECT id, content, comment_date, user_id, post_id, soft_deleted, deletion_date FROM comment WHERE id = $1 AND soft_deleted = false`
	rows, err := bd.db.QueryContext(ctx, query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return bd.scanComment(rows)
	}

	return nil, sql.ErrNoRows
}

func (bd *BDao) GetAllPosts() (*[]*sioblog.BlogPost, error) {
	query := `SELECT id, title, body, created_by_id, updated_by_id, posted_date, updated_date, deletion_date, soft_deleted FROM post WHERE soft_deleted = false`
	rows, err := bd.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*sioblog.BlogPost, 0)
	for rows.Next() {
		post, err := bd.scanPost(rows)
		if err != nil {
			return nil, err
		}
		comments, err := bd.GetAllCommentsByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}

	return &posts, nil
}

func (bd *BDao) GetAllCommentsByPostID(postID int64) (*[]*sioblog.BlogComment, error) {
	query := `SELECT id, content, comment_date, user_id, post_id, soft_deleted, deletion_date FROM comment WHERE post_id = $1 AND soft_deleted = false`
	rows, err := bd.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return bd.scanComments(rows)
}

func (bd *BDao) UpdatePost(post *sioblog.BlogPost) error {
	query := `UPDATE post
		SET 
		title = COALESCE($1, title),
		body = COALESCE($2,body),
		updated_date = $3,
		updated_by_id = $4
		WHERE id = $5
	`
	if _, err := bd.db.ExecContext(ctx, query, post.Title, post.Body, post.UpdatedDate, post.UpdatedByID, post.ID); err != nil {
		return err
	}
	return nil
}

func (bd *BDao) AddComment(comment *sioblog.BlogComment) error {
	query := `INSERT INTO comment (content, comment_date, post_id, user_id, soft_deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := bd.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.CommentDate,
		comment.PostID,
		comment.UserID,
		false,
	).Scan(&comment.ID)
	return err
}

func (bd *BDao) UpdateComment(comment *sioblog.BlogComment) error {
	query := `UPDATE comment
		SET 
		content = COALESCE($1, content),
		updated_date = COALESCE($2, updated_date)
		WHERE id = $3
	`
	if _, err := bd.db.ExecContext(ctx, query, comment.Content, comment.UpdatedDate, comment.ID); err != nil {
		return err
	}
	return nil
}

func (bd *BDao) SoftDeletePost(post *sioblog.BlogPost) error {
	query := `UPDATE post 
	SET
	soft_deleted = $1,
	deletion_date = $2
	where id = $3
	`

	if _, err := bd.db.ExecContext(ctx, query, true, post.DeletionDate, post.ID); err != nil {
		return err
	}
	return nil
}

func (bd *BDao) SoftDeleteComment(comment *sioblog.BlogComment) error {
	query := `UPDATE comment
		SET 
		soft_deleted = $1,
		deletion_date = $2
		WHERE id = $3
	`
	if _, err := bd.db.ExecContext(ctx, query, true, comment.DeletionDate, comment.ID); err != nil {
		return err
	}
	return nil
}

func (bd *BDao) scanPost(rows *sql.Rows) (*sioblog.BlogPost, error) {
	post := &sioblog.BlogPost{}
	err := rows.Scan(
		&post.ID,
		&post.Title,
		&post.Body,
		&post.CreatedByID,
		&post.UpdatedByID,
		&post.PostedDate,
		&post.UpdatedDate,
		&post.DeletionDate,
		&post.SoftDeleted,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (bd *BDao) scanComment(rows *sql.Rows) (*sioblog.BlogComment, error) {
	comment := &sioblog.BlogComment{}
	err := rows.Scan(
		&comment.ID,
		&comment.Content,
		&comment.CommentDate,
		&comment.UserID,
		&comment.PostID,
		&comment.SoftDeleted,
		&comment.DeletionDate,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (bd *BDao) scanComments(rows *sql.Rows) (*[]*sioblog.BlogComment, error) {
	comments := &[]*sioblog.BlogComment{}
	for rows.Next() {
		comment, err := bd.scanComment(rows)
		if err != nil {
			return nil, err
		}

		*comments = append(*comments, comment)
	}

	return comments, nil
}

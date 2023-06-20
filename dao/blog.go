package dao

import (
	"context"
	"database/sql"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/siodao"
)

var ctx = context.Background()

type BDao struct {
	db *sql.DB
}

//go:generate mockery --name BlogDao
type BlogDao interface {
	CreatePost(post *siogeneric.BlogPost) error
	PostExists(title string, createdByID int64) (bool, error)
	PostExistsByID(ID int64) (bool, error)
	CommentExistsByID(ID int64) (bool, error)
	GetPostByID(ID int64) (*siogeneric.BlogPost, error)
	GetCommentByID(ID int64) (*siogeneric.BlogComment, error)
	GetAllPosts() (*[]*siogeneric.BlogPost, error)
	GetAllCommentsByPostID(postID int64) (*[]*siogeneric.BlogComment, error)
	UpdatePost(post *siogeneric.BlogPost) error
	AddComment(comment *siogeneric.BlogComment) error
	UpdateComment(comment *siogeneric.BlogComment) error
	SoftDeletePost(post *siogeneric.BlogPost) error
	SoftDeleteComment(comment *siogeneric.BlogComment) error
}

func NewBlogDao() *BDao {
	return &BDao{
		db: siodao.DatabaseConnection(),
	}
}

func (bd *BDao) CreatePost(post *siogeneric.BlogPost) error {
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
	query := `SELECT EXISTS(SELECT 1 FROM post WHERE title = $1 created_by_id = $2)`
	var exists bool
	err := bd.db.QueryRowContext(ctx, query, title, createdByID).
		Scan(&exists)
	return exists, err
}

func (bd *BDao) PostExistsByID(ID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM post WHERE id = $1)`
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

func (bd *BDao) GetPostByID(ID int64) (*siogeneric.BlogPost, error) {
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

func (bd *BDao) GetCommentByID(ID int64) (*siogeneric.BlogComment, error) {
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

func (bd *BDao) GetAllPosts() (*[]*siogeneric.BlogPost, error) {
	query := `SELECT id, title, body, created_by_id, updated_by_id, posted_date, updated_date, deletion_date, soft_deleted FROM post WHERE soft_deleted = false`
	rows, err := bd.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*siogeneric.BlogPost, 0)
	for rows.Next() {
		post, err := bd.scanPost(rows)
		if err != nil {
			return nil, err
		}
		comments, err := bd.GetAllCommentsByPostID(post.ID.Int64)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}

	return &posts, nil
}

func (bd *BDao) GetAllCommentsByPostID(postID int64) (*[]*siogeneric.BlogComment, error) {
	query := `SELECT id, content, comment_date, user_id, post_id, soft_deleted, deletion_date FROM comment WHERE post_id = $1 AND soft_deleted = false`
	rows, err := bd.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return bd.scanComments(rows)
}

func (bd *BDao) UpdatePost(post *siogeneric.BlogPost) error {
	query := `UPDATE post
		SET 
		title = COALESCE($1, title),
		body = COALESCE($2,body),
		updated_date = COALESCE($3, updated_date)
		updated_by_id = COALESCE($4, updated_by_id)
	`
	if _, err := bd.db.ExecContext(ctx, query, post.Title, post.Body, post.UpdatedDate, post.UpdatedByID); err != nil {
		return err
	}
	return nil
}

func (bd *BDao) AddComment(comment *siogeneric.BlogComment) error {
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

func (bd *BDao) UpdateComment(comment *siogeneric.BlogComment) error {
	query := `UPDATE comment
		SET 
		content = COALESCE($1, content),
		updated_date = COALESCE($2, updated_date)
	`
	if _, err := bd.db.ExecContext(ctx, query, comment.Content, comment.UpdatedDate); err != nil {
		return err
	}
	return nil
}

func (bd *BDao) SoftDeletePost(post *siogeneric.BlogPost) error {
	query := `UPDATE post 
	SET
	soft_deleted = $1,
	delation_date = $2,
	where id = $3
	`

	if _, err := bd.db.ExecContext(ctx, query, true, post.DeletionDate, post.ID); err != nil {
		return err
	}
	return nil
}

func (bd *BDao) SoftDeleteComment(comment *siogeneric.BlogComment) error {
	query := `UPDATE comment
		SET 
		soft_deleted = $1,
		deletion_date = $2,
		WHERE id = $3
	`
	if _, err := bd.db.ExecContext(ctx, query, true, comment.DeletionDate, comment.ID); err != nil {
		return err
	}
	return nil
}

func (bd *BDao) scanPost(rows *sql.Rows) (*siogeneric.BlogPost, error) {
	post := &siogeneric.BlogPost{}
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

func (bd *BDao) scanComment(rows *sql.Rows) (*siogeneric.BlogComment, error) {
	comment := &siogeneric.BlogComment{}
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

func (bd *BDao) scanComments(rows *sql.Rows) (*[]*siogeneric.BlogComment, error) {
	comments := &[]*siogeneric.BlogComment{}
	for rows.Next() {
		comment, err := bd.scanComment(rows)
		if err != nil {
			return nil, err
		}

		*comments = append(*comments, comment)
	}

	return comments, nil
}

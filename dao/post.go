package dao

import (
	"context"
	"database/sql"
	"fmt"

	"gitea.slauson.io/blog/post-ms/constants"
	"gitea.slauson.io/slausonio/go-types/sioblog"
	"gitea.slauson.io/slausonio/go-utils/siodao"
)

var ctx = context.Background()

type PDao struct {
	db *sql.DB
}

//go:generate mockery --name PostDao
type PostDao interface {
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
	SoftDeletePost(post *sioblog.BlogPost) (sql.Result, error)
	SoftDeleteComment(comment *sioblog.BlogComment) (sql.Result, error)
}

func NewBlogDao() *PDao {
	return &PDao{
		db: siodao.DatabaseConnection(),
	}
}

func (pd *PDao) CreatePost(post *sioblog.BlogPost) error {
	query := `INSERT INTO post (title, body, posted_date, created_by_id, soft_deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := pd.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Body,
		post.PostedDate,
		post.CreatedByID,
		post.SoftDeleted,
	).Scan(&post.ID)
	return err
}

func (pd *PDao) PostExists(title string, createdByID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM post WHERE title = $1 AND created_by_id = $2 AND soft_deleted = false)`
	var exists bool
	err := pd.db.QueryRowContext(ctx, query, title, createdByID).
		Scan(&exists)
	return exists, err
}

func (pd *PDao) PostExistsByID(ID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM post WHERE id = $1 AND soft_deleted = false)`
	var exists bool
	err := pd.db.QueryRowContext(ctx, query, ID).Scan(&exists)
	return exists, err
}

func (pd *PDao) CommentExistsByID(ID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM comment WHERE id = $1)`
	var exists bool
	err := pd.db.QueryRowContext(ctx, query, ID).Scan(&exists)
	return exists, err
}

func (pd *PDao) GetPostByID(ID int64) (*sioblog.BlogPost, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM post WHERE id = $1 AND soft_deleted = false`,
		constants.SELECT_ITEMS_POST,
	)

	rows, err := pd.db.QueryContext(ctx, query, ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		post := &sioblog.BlogPost{}

		err := pd.scanPost(rows, post)
		if err != nil {
			return nil, err
		}

		comments, err := pd.GetAllCommentsByPostID(ID)
		if err != nil {
			return nil, err
		}

		post.Comments = comments
		return post, nil
	}

	return nil, sql.ErrNoRows
}

func (pd *PDao) GetCommentByID(ID int64) (*sioblog.BlogComment, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM comment WHERE id = $1 AND soft_deleted = false`,
		constants.SELECT_ITEMS_COMMENT,
	)

	rows, err := pd.db.QueryContext(ctx, query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		comment := &sioblog.BlogComment{}

		err := pd.scanComment(rows, comment)
		if err != nil {
			return nil, err
		}

		return comment, nil
	}

	return nil, sql.ErrNoRows
}

func (pd *PDao) GetAllPosts() (*[]*sioblog.BlogPost, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM post WHERE soft_deleted = false`,
		constants.SELECT_ITEMS_POST,
	)

	rows, err := pd.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	posts := make([]*sioblog.BlogPost, 0)

	for rows.Next() {
		post := &sioblog.BlogPost{}

		err := pd.scanPost(rows, post)
		if err != nil {
			return nil, err
		}

		comments, err := pd.GetAllCommentsByPostID(post.ID)
		if err != nil {
			return nil, err
		}

		post.Comments = comments
		posts = append(posts, post)
	}

	return &posts, nil
}

func (pd *PDao) GetAllCommentsByPostID(postID int64) (*[]*sioblog.BlogComment, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM comment WHERE post_id = $1 AND soft_deleted = false`,
		constants.SELECT_ITEMS_COMMENT,
	)

	rows, err := pd.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return pd.scanComments(rows)
}

func (pd *PDao) UpdatePost(post *sioblog.BlogPost) error {
	query := fmt.Sprintf(`UPDATE post
		SET 
		title = COALESCE($1, title),
		body = COALESCE($2,body),
		updated_date = $3,
		updated_by_id = $4
		WHERE id = $5
		returning %s`, constants.SELECT_ITEMS_POST)

	rows, err := pd.db.QueryContext(
		ctx,
		query,
		post.Title,
		post.Body,
		post.UpdatedDate,
		post.UpdatedByID,
		post.ID,
	)
	if err != nil {
		return err
	}

	defer rows.Close()

	return pd.scanPost(rows, post)
}

func (pd *PDao) AddComment(comment *sioblog.BlogComment) error {
	query := `INSERT INTO comment (content, comment_date, post_id, user_id, soft_deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := pd.db.QueryRowContext(
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

func (pd *PDao) UpdateComment(comment *sioblog.BlogComment) error {
	query := fmt.Sprintf(`UPDATE comment
		SET 
		content = COALESCE($1, content),
		updated_date = COALESCE($2, updated_date)
		WHERE id = $3 returning %s`, constants.SELECT_ITEMS_COMMENT)

	if _, err := pd.db.ExecContext(ctx, query, comment.Content, comment.UpdatedDate, comment.ID); err != nil {
		return err
	}

	return nil
}

func (pd *PDao) SoftDeletePost(post *sioblog.BlogPost) (sql.Result, error) {
	query := `UPDATE post 
	SET
	soft_deleted = $1,
	deletion_date = $2
	where id = $3`

	return pd.db.ExecContext(ctx, query, true, post.DeletionDate, post.ID)
}

func (pd *PDao) SoftDeleteComment(comment *sioblog.BlogComment) (sql.Result, error) {
	query := `UPDATE comment
		SET 
		soft_deleted = $1,
		deletion_date = $2
		WHERE id = $3`
	return pd.db.ExecContext(ctx, query, true, comment.DeletionDate, comment.ID)
}

func (pd *PDao) scanPost(rows *sql.Rows, post *sioblog.BlogPost) error {
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
		return err
	}

	return nil
}

func (pd *PDao) scanComment(rows *sql.Rows, comment *sioblog.BlogComment) error {
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
		return err
	}

	return nil
}

func (pd *PDao) scanComments(rows *sql.Rows) (*[]*sioblog.BlogComment, error) {
	comments := &[]*sioblog.BlogComment{}
	for rows.Next() {
		comment := &sioblog.BlogComment{}

		err := pd.scanComment(rows, comment)
		if err != nil {
			return nil, err
		}

		*comments = append(*comments, comment)
	}

	return comments, nil
}

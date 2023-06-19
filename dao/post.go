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
	GetAllCommentsByPostID(postID int64) (*[]siogeneric.BlogComment, error)
	UpdatePost(post *siogeneric.BlogPost) error
	AddComment(comment *siogeneric.BlogComment) error
	UpdateComment(comment *siogeneric.BlogComment) error
	SoftDeletePost(post *siogeneric.BlogPost) error
	SoftDeleteComment(comment *siogeneric.BlogComment) error
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

func (pd *PDao) GetPostByID(ID int64) (*siogeneric.BlogPost, error) {
	sql := `SELECT id, title, body, posted_date, updated_date, deletion_date, soft_deleted FROM blog_posts WHERE id = $1 AND soft_deleted = false`
	rows, err := pd.db.QueryContext(ctx, sql, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post *siogeneric.BlogPost
	for rows.Next() {
		post, err := pd.scanPost(rows)
		if err != nil {
			return nil, err
		}

		comments, err := pd.GetAllCommentsByPostID(ID)
		if err != nil {
			return nil, err
		}

		post.Comments = comments
	}
	return post, nil
}

func (pd *PDao) GetAllPosts() (*[]siogeneric.BlogPost, error) {
	sql := `SELECT id, title, body, posted_date, updated_date, deletion_date, soft_deleted FROM blog_posts WHERE soft_deleted = false`
	rows, err := pd.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]siogeneric.BlogPost, 0)
	for rows.Next() {
		post, err := pd.scanPost(rows)
		if err != nil {
			return nil, err
		}
		comments, err := pd.GetAllCommentsByPostID(post.ID.Int64)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, *post)
	}
	return &posts, nil
}

func (pd *PDao) GetAllCommentsByPostID(postID int64) (*[]siogeneric.BlogComment, error) {
	sql := `SELECT id, content, comment_date, soft_deleted, deletion_date FROM blog_comments WHERE post_id = $1 AND soft_deleted = false`
	rows, err := pd.db.QueryContext(ctx, sql, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pd.scanComment(rows)
}

func (pd *PDao) UpdatePost(post *siogeneric.BlogPost) error {
	sql := `UPDATE blog
		SET 
		title = COALESCE($1, title),
		body = COALESCE($2,body),
		updated_date = COALESCE($3, updated_date)
		updated_by_id = COALESCE($4, updated_by_id)
	`
	if _, err := pd.db.ExecContext(ctx, sql, post.Title, post.Body, post.UpdatedDate, post.UpdatedByID); err != nil {
		return err
	}
	return nil
}

func (pd *PDao) AddComment(comment *siogeneric.BlogComment) error {
	sql := `INSERT INTO blog_comments (content, comment_date, post_id, user_id, soft_deleted) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := pd.db.QueryRowContext(
		ctx,
		sql,
		comment.Content,
		comment.CommentDate,
		comment.PostID,
		comment.UserID,
		false,
	).Scan(&comment.ID)
	return err
}

func (pd *PDao) UpdateComment(comment *siogeneric.BlogComment) error {
	sql := `UPDATE blog_comments
		SET 
		content = COALESCE($1, content),
		updated_date = COALESCE($2, updated_date)
	`
	if _, err := pd.db.ExecContext(ctx, sql, comment.Content, comment.UpdatedDate); err != nil {
		return err
	}
	return nil
}

func (pd *PDao) SoftDeletePost(post *siogeneric.BlogPost) error {
	sql := `UPDATE blog 
	SET
	soft_deleted = $1,
	delation_date = $2,
	where id = $3
	`

	if _, err := pd.db.ExecContext(ctx, sql, true, post.DeletionDate, post.ID); err != nil {
		return err
	}
	return nil
}

func (pd *PDao) SoftDeleteComment(comment *siogeneric.BlogComment) error {
	sql := `UPDATE comment
		SET 
		soft_deleted = $1,
		deletion_date = $2,
		WHERE id = $3
	`
	if _, err := pd.db.ExecContext(ctx, sql, true, comment.DeletionDate, comment.ID); err != nil {
		return err
	}
	return nil
}

func (pd *PDao) scanPost(rows *sql.Rows) (*siogeneric.BlogPost, error) {
	post := &siogeneric.BlogPost{}
	err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.PostedDate,
		&post.UpdatedDate, &post.DeletionDate, &post.SoftDeleted)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pd *PDao) scanComment(rows *sql.Rows) (*[]siogeneric.BlogComment, error) {
	comments := &[]siogeneric.BlogComment{}
	for rows.Next() {
		comment := siogeneric.BlogComment{}
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.CommentDate,
			&comment.SoftDeleted,
			&comment.DeletionDate,
		)
		if err != nil {
			return nil, err
		}
		*comments = append(*comments, comment)
	}

	return comments, nil
}

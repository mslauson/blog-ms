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

func (pd *PDao) GetPostByID(ID int64) (*siogeneric.BlogPost, error) {
	sql := `SELECT id, title, body, posted_date, updated_date, deletion_date, soft_deleted FROM blog_posts WHERE id = $1 AND soft_deleted = false`
	row := pd.db.QueryRowContext(ctx, sql, ID)
	post, err := pd.scanPost(row)
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

// func (pd *PDao) GetAllPosts() (*[]siogeneric.BlogPost, error) {
// 	sql := `SELECT
// 			p.id, p.title, p.body, p.posted_date, p.updated_date, p.deletion_date, p.soft_deleted,
// 			c.id, c.content, c.comment_date, c.soft_deleted, c.deletion_date where p.soft_deleted = false
// 		FROM blog.post p
// 		LEFT JOIN blog.comment c ON p.id = c.post_id`
//
// 	rows, err := pd.db.QueryContext(ctx, sql)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	posts := make(map[int64]*siogeneric.BlogPost)
// 	for rows.Next() {
// 		comment := &siogeneric.BlogComment{}
// 		post := &siogeneric.BlogPost{}
// 		if err := pd.scanPost(rows, post); err != nil {
// 			return nil, err
// 		}
// 		if _, exists := posts[post.ID.Int64]; !exists {
// 			posts[post.ID.Int64] = post
// 		}
// 		if err := pd.scanComment(rows, comment); err != nil {
// 			return nil, err
// 		}
// 		posts[post.ID.Int64].Comments = append(posts[post.ID.Int64].Comments, comment)
// 	}
//
// 	// Convert map values to slice
// 	allPosts := make([]BlogPost, 0, len(posts))
// 	for _, post := range posts {
// 		allPosts = append(allPosts, *post)
// 	}
//
// 	return allPosts, nil}

func (pd *PDao) GetAllCommentsByPostID(postID int64) (*[]siogeneric.BlogComment, error) {
	sql := `SELECT id, content, comment_date, soft_deleted, deletion_date FROM blog_comments WHERE post_id = $1 AND soft_deleted = false`
	rows, err := pd.db.QueryContext(ctx, sql, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := make([]siogeneric.BlogComment, 0)
	return &comments, nil
}

func (pd *PDao) scanPost(row *sql.Row) (*siogeneric.BlogPost, error) {
	post := &siogeneric.BlogPost{}
	err := row.Scan(&post.ID, &post.Title, &post.Body, &post.PostedDate,
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

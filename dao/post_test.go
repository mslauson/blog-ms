package dao

import (
	"testing"
	"time"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// TODO: Add Returns
func TestPostDao(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	// Define common data
	post := &siogeneric.BlogPost{
		ID:          siogeneric.NewSioNullInt64(1),
		Title:       siogeneric.NewSioNullString("Test title"),
		Body:        siogeneric.NewSioNullString("Test body"),
		PostedDate:  siogeneric.NewSioNullTime(time.Now()),
		CreatedByID: siogeneric.NewSioNullInt64(1),
		SoftDeleted: siogeneric.NewSioNullBool(false),
	}

	comment := &siogeneric.BlogComment{
		ID:          siogeneric.NewSioNullInt64(1),
		Content:     siogeneric.NewSioNullString("Test comment"),
		CommentDate: siogeneric.NewSioNullTime(time.Now()),
		PostID:      siogeneric.NewSioNullInt64(1),
		UserID:      siogeneric.NewSioNullInt64(1),
		SoftDeleted: siogeneric.NewSioNullBool(false),
	}

	t.Run("CreatePost", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`INSERT INTO post`).
			WithArgs(post.Title.String, post.Body.String, post.PostedDate.Time, post.CreatedByID.Int64, false).
			WillReturnRows(rows)
		err = pd.CreatePost(post)
		require.NoError(t, err)
		require.Equal(t, int64(1), post.ID.Int64)
	})

	t.Run("PostExists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs(post.Title.String, post.CreatedByID.Int64).
			WillReturnRows(rows)
		exists, err := pd.PostExists(post.Title.String, post.CreatedByID.Int64)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("PostExistsByID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs(post.ID.Int64).
			WillReturnRows(rows)
		exists, err := pd.PostExistsByID(post.ID.Int64)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("CommentExistsByID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs(comment.ID.Int64).
			WillReturnRows(rows)
		exists, err := pd.CommentExistsByID(comment.ID.Int64)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("GetPostByID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "body", "created_by_id", "updated_by_id", "posted_date", "updated_date", "deletion_date", "soft_deleted"}).
			AddRow(post.ID.Int64, post.Title.String, post.Body.String, post.CreatedByID, post.UpdatedByID, post.PostedDate.Time, post.UpdatedDate, post.DeletionDate, post.SoftDeleted.Bool)
		mock.ExpectQuery(`SELECT id, title, body, created_by_id, updated_by_id,posted_date, updated_date, deletion_date, soft_deleted FROM post WHERE id = $1 AND soft_deleted = false`).
			WithArgs(post.ID.Int64).
			WillReturnRows(rows)
		returnedPost, err := pd.GetPostByID(post.ID.Int64)
		require.NoError(t, err)
		require.Equal(t, post, returnedPost)
	})

	t.Run("GetCommentByID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "content", "comment_date", "post_id", "user_id", "soft_deleted"}).
			AddRow(comment.ID.Int64, comment.Content.String, comment.CommentDate.Time, comment.PostID.Int64, comment.UserID.Int64, comment.SoftDeleted.Bool)
		mock.ExpectQuery(`SELECT id, content, comment_date, post_id, user_id, soft_deleted FROM comment WHERE id = $1`).
			WithArgs(comment.ID.Int64).
			WillReturnRows(rows)
		returnedComment, err := pd.GetCommentByID(comment.ID.Int64)
		require.NoError(t, err)
		require.Equal(t, comment, returnedComment)
	})

	t.Run("GetAllPosts", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "body", "created_by_id", "updated_by_id", "posted_date", "updated_date", "deletion_date", "soft_deleted"}).
			AddRow(post.ID.Int64, post.Title.String, post.Body.String, post.CreatedByID, post.UpdatedByID, post.PostedDate.Time, post.UpdatedDate, post.DeletionDate, post.SoftDeleted.Bool)
		mock.ExpectQuery(`SELECT id, title, body, "created_by_id", "updated_by_id", posted_date, updated_date, deletion_date, soft_deleted FROM post WHERE soft_deleted = false`).
			WillReturnRows(rows)
		posts, err := pd.GetAllPosts()
		require.NoError(t, err)
		require.Equal(t, []*siogeneric.BlogPost{post}, *posts)
	})

	t.Run("GetAllCommentsByPostID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "content", "comment_date", "soft_deleted", "deletion_date"}).
			AddRow(comment.ID.Int64, comment.Content.String, comment.CommentDate.Time, comment.SoftDeleted.Bool, comment.DeletionDate.Time)
		mock.ExpectQuery(`SELECT id, content, comment_date, soft_deleted, deletion_date FROM comment WHERE post_id = $1 AND soft_deleted = false`).
			WithArgs(post.ID.Int64).
			WillReturnRows(rows)
		comments, err := pd.GetAllCommentsByPostID(post.ID.Int64)
		require.NoError(t, err)
		require.Equal(t, []*siogeneric.BlogComment{comment}, *comments)
	})

	t.Run("UpdatePost", func(t *testing.T) {
		mock.ExpectExec(`UPDATE blog`).
			WithArgs(post.Title.String, post.Body.String, post.UpdatedDate.Time, post.UpdatedByID.Int64).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := pd.UpdatePost(post)
		require.NoError(t, err)
	})

	t.Run("AddComment", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`INSERT INTO comment`).
			WithArgs(comment.Content.String, comment.CommentDate.Time, comment.PostID.Int64, comment.UserID.Int64, false).
			WillReturnRows(rows)
		err = pd.AddComment(comment)
		require.NoError(t, err)
		require.Equal(t, int64(1), comment.ID.Int64)
	})

	t.Run("UpdateComment", func(t *testing.T) {
		mock.ExpectExec(`UPDATE comment`).
			WithArgs(comment.Content.String, comment.UpdatedDate.Time).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := pd.UpdateComment(comment)
		require.NoError(t, err)
	})

	t.Run("SoftDeletePost", func(t *testing.T) {
		mock.ExpectExec(`UPDATE blog`).
			WithArgs(true, post.DeletionDate.Time, post.ID.Int64).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := pd.SoftDeletePost(post)
		require.NoError(t, err)
	})

	t.Run("SoftDeleteComment", func(t *testing.T) {
		mock.ExpectExec(`UPDATE comment`).
			WithArgs(true, comment.DeletionDate.Time, comment.ID.Int64).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := pd.SoftDeleteComment(comment)
		require.NoError(t, err)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

package dao

import (
	"testing"
	"time"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// Define common data
var post = &siogeneric.BlogPost{
	ID:           siogeneric.NewSioNullInt64(1),
	Title:        siogeneric.NewSioNullString("Test title"),
	Body:         siogeneric.NewSioNullString("Test body"),
	PostedDate:   siogeneric.NewSioNullTime(time.Now()),
	CreatedByID:  siogeneric.NewSioNullInt64(1),
	UpdatedByID:  siogeneric.NewSioNullInt64(1),
	SoftDeleted:  siogeneric.NewSioNullBool(false),
	DeletionDate: siogeneric.NewSioNullTime(time.Now()),
	UpdatedDate:  siogeneric.NewSioNullTime(time.Now()),
}

var comment = &siogeneric.BlogComment{
	ID:           siogeneric.NewSioNullInt64(1),
	Content:      siogeneric.NewSioNullString("Test comment"),
	CommentDate:  siogeneric.NewSioNullTime(time.Now()),
	PostID:       siogeneric.NewSioNullInt64(1),
	UserID:       siogeneric.NewSioNullInt64(1),
	SoftDeleted:  siogeneric.NewSioNullBool(false),
	DeletionDate: siogeneric.NewSioNullTime(time.Now()),
	UpdatedDate:  siogeneric.NewSioNullTime(time.Now()),
}

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(`INSERT INTO post`).
		WithArgs(post.Title.String, post.Body.String, post.PostedDate.Time, post.CreatedByID.Int64, false).
		WillReturnRows(rows)
	err = pd.CreatePost(post)
	require.NoError(t, err)
	require.Equal(t, int64(1), post.ID.Int64)
}

func TestPostExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(post.Title.String, post.CreatedByID.Int64).
		WillReturnRows(rows)
	exists, err := pd.PostExists(post.Title.String, post.CreatedByID.Int64)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestPostExistsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(`SELECT EXISTS`).
		WillReturnRows(rows)
	exists, err := pd.PostExistsByID(post.ID.Int64)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestCommentExistsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(`SELECT EXISTS`).
		WillReturnRows(rows)
	exists, err := pd.CommentExistsByID(comment.ID.Int64)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestGetPostByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	rows := sqlmock.NewRows([]string{
		"id",
		"title",
		"body",
		"created_by_id",
		"updated_by_id",
		"posted_date",
		"updated_date",
		"deletion_date",
		"soft_deleted",
	}).
		AddRow(post.ID.Int64, post.Title.String, post.Body.String, post.CreatedByID, post.UpdatedByID, post.PostedDate.Time, post.UpdatedDate, post.DeletionDate, post.SoftDeleted.Bool)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)
	returnedPost, err := pd.GetPostByID(post.ID.Int64)
	require.NoError(t, err)
	require.Equal(t, post, returnedPost)
}

func TestGetCommentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	rows := sqlmock.NewRows([]string{
		"id",
		"content",
		"comment_date",
		"post_id",
		"user_id",
		"soft_deleted",
		"deletion_date",
	}).
		AddRow(
			comment.ID.Int64,
			comment.Content.String,
			comment.CommentDate.Time,
			comment.PostID.Int64,
			comment.UserID.Int64,
			comment.SoftDeleted.Bool,
			comment.DeletionDate,
		)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)
	returnedComment, err := pd.GetCommentByID(comment.ID.Int64)
	require.NoError(t, err)
	require.Equal(t, comment, returnedComment)
}

func TestGetAllPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	rows := sqlmock.NewRows([]string{"id", "title", "body", "created_by_id", "updated_by_id", "posted_date", "updated_date", "deletion_date", "soft_deleted"}).
		AddRow(post.ID.Int64, post.Title.String, post.Body.String, post.CreatedByID, post.UpdatedByID, post.PostedDate.Time, post.UpdatedDate, post.DeletionDate, post.SoftDeleted.Bool)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)
	posts, err := pd.GetAllPosts()
	require.NoError(t, err)
	require.Equal(t, []*siogeneric.BlogPost{post}, *posts)
}

func TestGetAllCommentsByPostID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	rows := sqlmock.NewRows([]string{
		"id",
		"content",
		"comment_date",
		"post_id",
		"user_id",
		"soft_deleted",
		"deletion_date",
	}).
		AddRow(
			comment.ID.Int64,
			comment.Content.String,
			comment.CommentDate.Time,
			comment.PostID.Int64,
			comment.UserID.Int64,
			comment.SoftDeleted.Bool,
			comment.DeletionDate,
		)

	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)
	comments, err := pd.GetAllCommentsByPostID(post.ID.Int64)
	require.NoError(t, err)
	require.Equal(t, []*siogeneric.BlogComment{comment}, *comments)
}

func TestUpdatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(post.Title.String, post.Body.String, post.UpdatedDate.Time, post.UpdatedByID.Int64).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.UpdatePost(post)
	require.NoError(t, err)
}

func TestAddComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(`INSERT INTO comment`).
		WithArgs(comment.Content.String, comment.CommentDate.Time, comment.PostID.Int64, comment.UserID.Int64, false).
		WillReturnRows(rows)
	err = pd.AddComment(comment)
	require.NoError(t, err)
	require.Equal(t, int64(1), comment.ID.Int64)
}

func TestUpdateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(comment.Content.String, comment.UpdatedDate.Time).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.UpdateComment(comment)
	require.NoError(t, err)
}

func TestSoftDeletePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(true, post.DeletionDate.Time, post.ID.Int64).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.SoftDeletePost(post)
	require.NoError(t, err)
}

func TestSoftDeleteComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &PDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(true, comment.DeletionDate.Time, comment.ID.Int64).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.SoftDeleteComment(comment)
	require.NoError(t, err)
}

// require.NoError(t, mock.ExpectationsWereMet())

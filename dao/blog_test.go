package dao

import (
	"errors"
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

	pd := &BDao{db: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(`INSERT INTO post`).
		WithArgs(post.Title.String, post.Body.String, post.PostedDate.Time, post.CreatedByID.Int64, false).
		WillReturnRows(rows)
	err = pd.CreatePost(post)
	require.NoError(t, err)
	require.Equal(t, int64(1), post.ID.Int64)
}

func TestCreatePost_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`INSERT INTO post`).
		WithArgs(post.Title.String, post.Body.String, post.PostedDate.Time, post.CreatedByID.Int64, false).
		WillReturnError(errors.New("test error"))
	err = pd.CreatePost(post)
	require.Error(t, err)
}

func TestPostExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(post.Title.String, post.CreatedByID.Int64).
		WillReturnRows(rows)
	exists, err := pd.PostExists(post.Title.String, post.CreatedByID.Int64)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestPostExists_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(post.Title.String, post.CreatedByID.Int64).
		WillReturnError(errors.New("test error"))
	exists, err := pd.PostExists(post.Title.String, post.CreatedByID.Int64)
	require.Error(t, err)
	require.False(t, exists)
}

func TestPostExistsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(`SELECT EXISTS`).
		WillReturnRows(rows)
	exists, err := pd.PostExistsByID(post.ID.Int64)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestPostExistsByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`SELECT EXISTS`).
		WillReturnError(errors.New("test error"))
	exists, err := pd.PostExistsByID(post.ID.Int64)
	require.Error(t, err)
	require.False(t, exists)
}

func TestCommentExistsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(`SELECT EXISTS`).
		WillReturnRows(rows)
	exists, err := pd.CommentExistsByID(comment.ID.Int64)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestCommentExistsByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`SELECT EXISTS`).
		WillReturnError(errors.New("test error"))
	exists, err := pd.CommentExistsByID(comment.ID.Int64)
	require.Error(t, err)
	require.False(t, exists)
}

func TestGetPostByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

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

	cRows := sqlmock.NewRows([]string{
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
		WillReturnRows(cRows)

	returnedPost, err := pd.GetPostByID(post.ID.Int64)
	require.NoError(t, err)
	require.NotNil(t, returnedPost)
}

func TestGetPostByID_ErrorPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("test error"))

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("test error"))

	returnedPost, err := pd.GetPostByID(post.ID.Int64)
	require.Error(t, err)
	require.Nil(t, returnedPost)
}

func TestGetPostByID_ErrorComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

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

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("test error"))

	returnedPost, err := pd.GetPostByID(post.ID.Int64)
	require.Error(t, err)
	require.Nil(t, returnedPost)
}

func TestGetCommentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

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
	require.NotNil(t, returnedComment)
}

func TestGetCommentByID_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("test error"))
	returnedComment, err := pd.GetCommentByID(comment.ID.Int64)
	require.Error(t, err)
	require.Nil(t, returnedComment)
}

func TestGetAllPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

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
		AddRow(post.ID.Int64,
			post.Title.String,
			post.Body.String,
			post.CreatedByID,
			post.UpdatedByID,
			post.PostedDate.Time,
			post.UpdatedDate,
			post.DeletionDate,
			post.SoftDeleted.Bool)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)

	cRows := sqlmock.NewRows([]string{
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
		WillReturnRows(cRows)

	posts, err := pd.GetAllPosts()
	require.NoError(t, err)
	require.NotNil(t, posts)
}

func TestGetAllPosts_ErrPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("test error"))

	posts, err := pd.GetAllPosts()
	require.Error(t, err)
	require.Nil(t, posts)
}

func TestGetAllPosts_ErrComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

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
		AddRow(post.ID.Int64,
			post.Title.String,
			post.Body.String,
			post.CreatedByID,
			post.UpdatedByID,
			post.PostedDate.Time,
			post.UpdatedDate,
			post.DeletionDate,
			post.SoftDeleted.Bool)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("test error"))

	posts, err := pd.GetAllPosts()
	require.Error(t, err)
	require.Nil(t, posts)
}

func TestGetAllCommentsByPostID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

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
	require.NotNil(t, comments)
}

func TestGetAllCommentsByPostID_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("test error"))

	comments, err := pd.GetAllCommentsByPostID(post.ID.Int64)
	require.Error(t, err)
	require.Nil(t, comments)
}

func TestUpdatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(post.Title.String, post.Body.String, post.UpdatedDate.Time, post.UpdatedByID.Int64).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.UpdatePost(post)
	require.NoError(t, err)
}

func TestUpdatePost_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(post.Title.String, post.Body.String, post.UpdatedDate.Time, post.UpdatedByID.Int64).
		WillReturnError(errors.New("test error"))
	err = pd.UpdatePost(post)
	require.Error(t, err)
}

func TestAddComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(`INSERT INTO comment`).
		WithArgs(comment.Content.String, comment.CommentDate.Time, comment.PostID.Int64, comment.UserID.Int64, false).
		WillReturnRows(rows)
	err = pd.AddComment(comment)
	require.NoError(t, err)
	require.Equal(t, int64(1), comment.ID.Int64)
}

func TestAddComment_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`INSERT INTO comment`).
		WithArgs(comment.Content.String, comment.CommentDate.Time, comment.PostID.Int64, comment.UserID.Int64, false).
		WillReturnError(errors.New("test error"))

	err = pd.AddComment(comment)
	require.Error(t, err)
}

func TestUpdateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(comment.Content.String, comment.UpdatedDate.Time).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.UpdateComment(comment)
	require.NoError(t, err)
}

func TestUpdateComment_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(comment.Content.String, comment.UpdatedDate.Time).
		WillReturnError(errors.New("test error"))

	err = pd.UpdateComment(comment)
	require.Error(t, err)
}

func TestSoftDeletePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(true, post.DeletionDate.Time, post.ID.Int64).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.SoftDeletePost(post)
	require.NoError(t, err)
}

func TestSoftDeletePost_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(true, post.DeletionDate.Time, post.ID.Int64).
		WillReturnError(errors.New("test error"))

	err = pd.SoftDeletePost(post)
	require.Error(t, err)
}

func TestSoftDeleteComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(true, comment.DeletionDate.Time, comment.ID.Int64).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.SoftDeleteComment(comment)
	require.NoError(t, err)
}

func TestSoftDeleteComment_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(true, comment.DeletionDate.Time, comment.ID.Int64).
		WillReturnError(errors.New("test error"))

	err = pd.SoftDeleteComment(comment)
	require.Error(t, err)
}

// require.NoError(t, mock.ExpectationsWereMet())

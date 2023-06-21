package dao

import (
	"errors"
	"testing"

	"gitea.slauson.io/blog/blog-ms/testing/mockdata"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(`INSERT INTO post`).
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.PostedDate.Time, mockdata.PostEntity.CreatedByID.Int64, false).
		WillReturnRows(rows)
	err = pd.CreatePost(mockdata.PostEntity)
	require.NoError(t, err)
	require.Equal(t, int64(1), mockdata.PostEntity.ID.Int64)
}

func TestCreatePost_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`INSERT INTO post`).
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.PostedDate.Time, mockdata.PostEntity.CreatedByID.Int64, false).
		WillReturnError(errors.New("test error"))
	err = pd.CreatePost(mockdata.PostEntity)
	require.Error(t, err)
}

func TestPostExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.CreatedByID.Int64).
		WillReturnRows(rows)
	exists, err := pd.PostExists(
		mockdata.PostEntity.Title.String,
		mockdata.PostEntity.CreatedByID.Int64,
	)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestPostExists_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.CreatedByID.Int64).
		WillReturnError(errors.New("test error"))
	exists, err := pd.PostExists(
		mockdata.PostEntity.Title.String,
		mockdata.PostEntity.CreatedByID.Int64,
	)
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
	exists, err := pd.PostExistsByID(mockdata.PostEntity.ID.Int64)
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
	exists, err := pd.PostExistsByID(mockdata.PostEntity.ID.Int64)
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
	exists, err := pd.CommentExistsByID(mockdata.CommentEntity.ID.Int64)
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
	exists, err := pd.CommentExistsByID(mockdata.CommentEntity.ID.Int64)
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
		AddRow(mockdata.PostEntity.ID.Int64, mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.CreatedByID, mockdata.PostEntity.UpdatedByID, mockdata.PostEntity.PostedDate.Time, mockdata.PostEntity.UpdatedDate, mockdata.PostEntity.DeletionDate, mockdata.PostEntity.SoftDeleted.Bool)
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
			mockdata.CommentEntity.ID.Int64,
			mockdata.CommentEntity.Content.String,
			mockdata.CommentEntity.CommentDate.Time,
			mockdata.CommentEntity.PostID.Int64,
			mockdata.CommentEntity.UserID.Int64,
			mockdata.CommentEntity.SoftDeleted.Bool,
			mockdata.CommentEntity.DeletionDate,
		)

	mock.ExpectQuery(`SELECT`).
		WillReturnRows(cRows)

	returnedPost, err := pd.GetPostByID(mockdata.PostEntity.ID.Int64)
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

	returnedPost, err := pd.GetPostByID(mockdata.PostEntity.ID.Int64)
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
		AddRow(mockdata.PostEntity.ID.Int64, mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.CreatedByID, mockdata.PostEntity.UpdatedByID, mockdata.PostEntity.PostedDate.Time, mockdata.PostEntity.UpdatedDate, mockdata.PostEntity.DeletionDate, mockdata.PostEntity.SoftDeleted.Bool)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("test error"))

	returnedPost, err := pd.GetPostByID(mockdata.PostEntity.ID.Int64)
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
			mockdata.CommentEntity.ID.Int64,
			mockdata.CommentEntity.Content.String,
			mockdata.CommentEntity.CommentDate.Time,
			mockdata.CommentEntity.PostID.Int64,
			mockdata.CommentEntity.UserID.Int64,
			mockdata.CommentEntity.SoftDeleted.Bool,
			mockdata.CommentEntity.DeletionDate,
		)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)
	returnedComment, err := pd.GetCommentByID(mockdata.CommentEntity.ID.Int64)
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
	returnedComment, err := pd.GetCommentByID(mockdata.CommentEntity.ID.Int64)
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
		AddRow(mockdata.PostEntity.ID.Int64,
			mockdata.PostEntity.Title.String,
			mockdata.PostEntity.Body.String,
			mockdata.PostEntity.CreatedByID,
			mockdata.PostEntity.UpdatedByID,
			mockdata.PostEntity.PostedDate.Time,
			mockdata.PostEntity.UpdatedDate,
			mockdata.PostEntity.DeletionDate,
			mockdata.PostEntity.SoftDeleted.Bool)
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
			mockdata.CommentEntity.ID.Int64,
			mockdata.CommentEntity.Content.String,
			mockdata.CommentEntity.CommentDate.Time,
			mockdata.CommentEntity.PostID.Int64,
			mockdata.CommentEntity.UserID.Int64,
			mockdata.CommentEntity.SoftDeleted.Bool,
			mockdata.CommentEntity.DeletionDate,
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
		AddRow(mockdata.PostEntity.ID.Int64,
			mockdata.PostEntity.Title.String,
			mockdata.PostEntity.Body.String,
			mockdata.PostEntity.CreatedByID,
			mockdata.PostEntity.UpdatedByID,
			mockdata.PostEntity.PostedDate.Time,
			mockdata.PostEntity.UpdatedDate,
			mockdata.PostEntity.DeletionDate,
			mockdata.PostEntity.SoftDeleted.Bool)
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
			mockdata.CommentEntity.ID.Int64,
			mockdata.CommentEntity.Content.String,
			mockdata.CommentEntity.CommentDate.Time,
			mockdata.CommentEntity.PostID.Int64,
			mockdata.CommentEntity.UserID.Int64,
			mockdata.CommentEntity.SoftDeleted.Bool,
			mockdata.CommentEntity.DeletionDate,
		)

	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)
	comments, err := pd.GetAllCommentsByPostID(mockdata.PostEntity.ID.Int64)
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

	comments, err := pd.GetAllCommentsByPostID(mockdata.PostEntity.ID.Int64)
	require.Error(t, err)
	require.Nil(t, comments)
}

func TestUpdatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.UpdatedDate.Time, mockdata.PostEntity.UpdatedByID.Int64).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.UpdatePost(mockdata.PostEntity)
	require.NoError(t, err)
}

func TestUpdatePost_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.UpdatedDate.Time, mockdata.PostEntity.UpdatedByID.Int64).
		WillReturnError(errors.New("test error"))
	err = pd.UpdatePost(mockdata.PostEntity)
	require.Error(t, err)
}

func TestAddComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(`INSERT INTO comment`).
		WithArgs(mockdata.CommentEntity.Content.String, mockdata.CommentEntity.CommentDate.Time, mockdata.CommentEntity.PostID.Int64, mockdata.CommentEntity.UserID.Int64, false).
		WillReturnRows(rows)
	err = pd.AddComment(mockdata.CommentEntity)
	require.NoError(t, err)
	require.Equal(t, int64(1), mockdata.CommentEntity.ID.Int64)
}

func TestAddComment_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`INSERT INTO comment`).
		WithArgs(mockdata.CommentEntity.Content.String, mockdata.CommentEntity.CommentDate.Time, mockdata.CommentEntity.PostID.Int64, mockdata.CommentEntity.UserID.Int64, false).
		WillReturnError(errors.New("test error"))

	err = pd.AddComment(mockdata.CommentEntity)
	require.Error(t, err)
}

func TestUpdateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(mockdata.CommentEntity.Content.String, mockdata.CommentEntity.UpdatedDate.Time).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.UpdateComment(mockdata.CommentEntity)
	require.NoError(t, err)
}

func TestUpdateComment_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(mockdata.CommentEntity.Content.String, mockdata.CommentEntity.UpdatedDate.Time).
		WillReturnError(errors.New("test error"))

	err = pd.UpdateComment(mockdata.CommentEntity)
	require.Error(t, err)
}

func TestSoftDeletePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(true, mockdata.PostEntity.DeletionDate.Time, mockdata.PostEntity.ID.Int64).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.SoftDeletePost(mockdata.PostEntity)
	require.NoError(t, err)
}

func TestSoftDeletePost_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(true, mockdata.PostEntity.DeletionDate.Time, mockdata.PostEntity.ID.Int64).
		WillReturnError(errors.New("test error"))

	err = pd.SoftDeletePost(mockdata.PostEntity)
	require.Error(t, err)
}

func TestSoftDeleteComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(true, mockdata.CommentEntity.DeletionDate.Time, mockdata.CommentEntity.ID.Int64).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = pd.SoftDeleteComment(mockdata.CommentEntity)
	require.NoError(t, err)
}

func TestSoftDeleteComment_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE comment`).
		WithArgs(true, mockdata.CommentEntity.DeletionDate.Time, mockdata.CommentEntity.ID.Int64).
		WillReturnError(errors.New("test error"))

	err = pd.SoftDeleteComment(mockdata.CommentEntity)
	require.Error(t, err)
}

// require.NoError(t, mock.ExpectationsWereMet())

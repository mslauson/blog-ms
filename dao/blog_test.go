package dao

import (
	"errors"
	"testing"

	"gitea.slauson.io/blog/post-ms/integration/mockdata"
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
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.PostedDate, mockdata.PostEntity.CreatedByID, false).
		WillReturnRows(rows)
	err = pd.CreatePost(mockdata.PostEntity)
	require.NoError(t, err)
	require.Equal(t, int64(1), mockdata.PostEntity.ID)
}

func TestCreatePost_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`INSERT INTO post`).
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.PostedDate, mockdata.PostEntity.CreatedByID, false).
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
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.CreatedByID).
		WillReturnRows(rows)
	exists, err := pd.PostExists(
		mockdata.PostEntity.Title.String,
		mockdata.PostEntity.CreatedByID,
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
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.CreatedByID).
		WillReturnError(errors.New("test error"))
	exists, err := pd.PostExists(
		mockdata.PostEntity.Title.String,
		mockdata.PostEntity.CreatedByID,
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
	exists, err := pd.PostExistsByID(mockdata.PostEntity.ID)
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
	exists, err := pd.PostExistsByID(mockdata.PostEntity.ID)
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
	exists, err := pd.CommentExistsByID(mockdata.CommentEntity.ID)
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
	exists, err := pd.CommentExistsByID(mockdata.CommentEntity.ID)
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
		AddRow(mockdata.PostEntity.ID, mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.CreatedByID, mockdata.PostEntity.UpdatedByID, mockdata.PostEntity.PostedDate, mockdata.PostEntity.UpdatedDate, mockdata.PostEntity.DeletionDate, mockdata.PostEntity.SoftDeleted)
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
			mockdata.CommentEntity.ID,
			mockdata.CommentEntity.Content.String,
			mockdata.CommentEntity.CommentDate,
			mockdata.CommentEntity.PostID,
			mockdata.CommentEntity.UserID,
			mockdata.CommentEntity.SoftDeleted,
			mockdata.CommentEntity.DeletionDate,
		)

	mock.ExpectQuery(`SELECT`).
		WillReturnRows(cRows)

	returnedPost, err := pd.GetPostByID(mockdata.PostEntity.ID)
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

	returnedPost, err := pd.GetPostByID(mockdata.PostEntity.ID)
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
		AddRow(mockdata.PostEntity.ID, mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.CreatedByID, mockdata.PostEntity.UpdatedByID, mockdata.PostEntity.PostedDate, mockdata.PostEntity.UpdatedDate, mockdata.PostEntity.DeletionDate, mockdata.PostEntity.SoftDeleted)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("test error"))

	returnedPost, err := pd.GetPostByID(mockdata.PostEntity.ID)
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
			mockdata.CommentEntity.ID,
			mockdata.CommentEntity.Content.String,
			mockdata.CommentEntity.CommentDate,
			mockdata.CommentEntity.PostID,
			mockdata.CommentEntity.UserID,
			mockdata.CommentEntity.SoftDeleted,
			mockdata.CommentEntity.DeletionDate,
		)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)
	returnedComment, err := pd.GetCommentByID(mockdata.CommentEntity.ID)
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
	returnedComment, err := pd.GetCommentByID(mockdata.CommentEntity.ID)
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
		AddRow(mockdata.PostEntity.ID,
			mockdata.PostEntity.Title.String,
			mockdata.PostEntity.Body.String,
			mockdata.PostEntity.CreatedByID,
			mockdata.PostEntity.UpdatedByID,
			mockdata.PostEntity.PostedDate,
			mockdata.PostEntity.UpdatedDate,
			mockdata.PostEntity.DeletionDate,
			mockdata.PostEntity.SoftDeleted)
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
			mockdata.CommentEntity.ID,
			mockdata.CommentEntity.Content.String,
			mockdata.CommentEntity.CommentDate,
			mockdata.CommentEntity.PostID,
			mockdata.CommentEntity.UserID,
			mockdata.CommentEntity.SoftDeleted,
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
		AddRow(mockdata.PostEntity.ID,
			mockdata.PostEntity.Title.String,
			mockdata.PostEntity.Body.String,
			mockdata.PostEntity.CreatedByID,
			mockdata.PostEntity.UpdatedByID,
			mockdata.PostEntity.PostedDate,
			mockdata.PostEntity.UpdatedDate,
			mockdata.PostEntity.DeletionDate,
			mockdata.PostEntity.SoftDeleted)
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
			mockdata.CommentEntity.ID,
			mockdata.CommentEntity.Content.String,
			mockdata.CommentEntity.CommentDate,
			mockdata.CommentEntity.PostID,
			mockdata.CommentEntity.UserID,
			mockdata.CommentEntity.SoftDeleted,
			mockdata.CommentEntity.DeletionDate,
		)

	mock.ExpectQuery(`SELECT`).
		WillReturnRows(rows)
	comments, err := pd.GetAllCommentsByPostID(mockdata.PostEntity.ID)
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

	comments, err := pd.GetAllCommentsByPostID(mockdata.PostEntity.ID)
	require.Error(t, err)
	require.Nil(t, comments)
}

func TestUpdatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectExec(`UPDATE post`).
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.UpdatedDate, mockdata.PostEntity.UpdatedByID, mockdata.PostEntity.ID).
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
		WithArgs(mockdata.PostEntity.Title.String, mockdata.PostEntity.Body.String, mockdata.PostEntity.UpdatedDate, mockdata.PostEntity.UpdatedByID, mockdata.PostEntity.ID).
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
		WithArgs(mockdata.CommentEntity.Content.String, mockdata.CommentEntity.CommentDate, mockdata.CommentEntity.PostID, mockdata.CommentEntity.UserID, false).
		WillReturnRows(rows)
	err = pd.AddComment(mockdata.CommentEntity)
	require.NoError(t, err)
	require.Equal(t, int64(1), mockdata.CommentEntity.ID)
}

func TestAddComment_Err(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	pd := &BDao{db: db}

	mock.ExpectQuery(`INSERT INTO comment`).
		WithArgs(mockdata.CommentEntity.Content.String, mockdata.CommentEntity.CommentDate, mockdata.CommentEntity.PostID, mockdata.CommentEntity.UserID, false).
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
		WithArgs(mockdata.CommentEntity.Content.String, mockdata.CommentEntity.UpdatedDate, mockdata.CommentEntity.ID).
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
		WithArgs(mockdata.CommentEntity.Content.String, mockdata.CommentEntity.UpdatedDate, mockdata.CommentEntity.ID).
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
		WithArgs(true, mockdata.PostEntity.DeletionDate, mockdata.PostEntity.ID).
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
		WithArgs(true, mockdata.PostEntity.DeletionDate, mockdata.PostEntity.ID).
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
		WithArgs(true, mockdata.CommentEntity.DeletionDate, mockdata.CommentEntity.ID).
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
		WithArgs(true, mockdata.CommentEntity.DeletionDate, mockdata.CommentEntity.ID).
		WillReturnError(errors.New("test error"))

	err = pd.SoftDeleteComment(mockdata.CommentEntity)
	require.Error(t, err)
}

// require.NoError(t, mock.ExpectationsWereMet())

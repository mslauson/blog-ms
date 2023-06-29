package service

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitea.slauson.io/blog/post-ms/constants"
	"gitea.slauson.io/blog/post-ms/dao/mocks"
	"gitea.slauson.io/blog/post-ms/integration/mockdata"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
)

var (
	updateResultHappy = sqlmock.NewResult(1, 1)
	updateResultFail  = sqlmock.NewResult(0, 0)
)

func initEnv(t *testing.T) (*PostSvc, *mocks.PostDao) {
	mDao := mocks.NewPostDao(t)
	svc := &PostSvc{dao: mDao}
	return svc, mDao
}

func TestGetPost(t *testing.T) {
	bs, dao := initEnv(t)

	testID := int64(1)

	dao.On("GetPostByID", testID).Return(mockdata.PostEntity, nil)

	resp, err := bs.GetPost(testID)
	assert.NoError(t, err)
	assert.Equal(t, mockdata.PostEntity.Title.String, resp.Title)
	dao.AssertExpectations(t)
}

func TestGetPost_WithComments(t *testing.T) {
	bs, dao := initEnv(t)

	testID := int64(1)
	mockdata.PostEntity.Comments = mockdata.Comments

	dao.On("GetPostByID", testID).Return(mockdata.PostEntity, nil)

	resp, err := bs.GetPost(testID)
	assert.NoError(t, err)
	assert.Equal(t, mockdata.PostEntity.Title.String, resp.Title)
	dao.AssertExpectations(t)
}

func TestGetPost_NotFound(t *testing.T) {
	bs, dao := initEnv(t)

	testID := int64(1)
	mockdata.PostEntity.Comments = mockdata.Comments

	testError := sioerror.NewSioNotFoundError("post not found")

	dao.On("GetPostByID", testID).Return(nil, sql.ErrNoRows)

	resp, err := bs.GetPost(testID)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestGetAllPosts(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("GetAllPosts").Return(mockdata.Posts, nil)

	resp, err := bs.GetAllPosts()
	assert.NoError(t, err)
	assert.Equal(t, len(*mockdata.Posts), len(*resp))
	dao.AssertExpectations(t)
}

func TestGetAllPosts_NotFound(t *testing.T) {
	bs, dao := initEnv(t)
	testError := sioerror.NewSioNotFoundError("posts not found")

	dao.On("GetAllPosts").Return(nil, sql.ErrNoRows)

	resp, err := bs.GetAllPosts()
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestCreatePost(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("PostExists", mockdata.CreatePostRequest.Title, mockdata.CreatePostRequest.CreatedByID).
		Return(false, nil)
	dao.On("CreatePost", mock.AnythingOfType("*sioblog.BlogPost")).Return(nil)

	resp, err := bs.CreatePost(mockdata.CreatePostRequest)
	assert.NoError(t, err)
	assert.Equal(t, mockdata.PostEntity.Title.String, resp.Title)
	assert.Equal(t, mockdata.PostEntity.Body.String, resp.Body)
	dao.AssertExpectations(t)
}

func TestCreatePost_ErrAlreadyExists(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioBadRequestError("post already exists")

	dao.On("PostExists", mockdata.CreatePostRequest.Title, mockdata.CreatePostRequest.CreatedByID).
		Return(true, nil)

	resp, err := bs.CreatePost(mockdata.CreatePostRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestCreatePost_ErrExistsCheckUnexpectedErr(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioInternalServerError(
		"unexpected DB error: " + sql.ErrConnDone.Error(),
	)

	dao.On("PostExists", mockdata.CreatePostRequest.Title, mockdata.CreatePostRequest.CreatedByID).
		Return(false, sql.ErrConnDone)

	resp, err := bs.CreatePost(mockdata.CreatePostRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestCreatePost_CreateErr(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioInternalServerError(
		"unexpected DB error: " + sql.ErrConnDone.Error(),
	)

	dao.On("PostExists", mockdata.CreatePostRequest.Title, mockdata.CreatePostRequest.CreatedByID).
		Return(false, nil)
	dao.On("CreatePost", mock.AnythingOfType("*sioblog.BlogPost")).Return(sql.ErrConnDone)

	resp, err := bs.CreatePost(mockdata.CreatePostRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestAddComment(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("PostExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("AddComment", mock.AnythingOfType("*sioblog.BlogComment")).Return(nil)

	resp, err := bs.AddComment(mockdata.AddCommentRequest)
	assert.NoError(t, err)
	assert.Equal(t, mockdata.CommentEntity.Content.String, resp.Content)
	dao.AssertExpectations(t)
}

func TestAddComment_ErrNotFound(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioNotFoundError("no post found")

	dao.On("PostExistsByID", mock.AnythingOfType("int64")).Return(false, nil)

	resp, err := bs.AddComment(mockdata.AddCommentRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestAddComment_ErrExistsCheckUnexpectedErr(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioInternalServerError(
		"unexpected DB error: " + sql.ErrConnDone.Error(),
	)

	dao.On("PostExistsByID", mock.AnythingOfType("int64")).Return(false, sql.ErrConnDone)

	resp, err := bs.AddComment(mockdata.AddCommentRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestAddComment_CreateErr(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioInternalServerError(
		"unexpected DB error: " + sql.ErrConnDone.Error(),
	)

	dao.On("PostExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("AddComment", mock.AnythingOfType("*sioblog.BlogComment")).Return(sql.ErrConnDone)

	resp, err := bs.AddComment(mockdata.AddCommentRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestUpdatePost(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("UpdatePost", mock.AnythingOfType("*sioblog.BlogPost")).Return(updateResultHappy, nil)
	dao.On("GetPostByID", mock.AnythingOfType("int64")).Return(mockdata.PostEntity, nil)

	resp, err := bs.UpdatePost(mockdata.PostEntity.ID, mockdata.UpdatePostRequest)
	assert.NoError(t, err)
	assert.Equal(t, mockdata.PostEntity.Title.String, resp.Title)
	dao.AssertExpectations(t)
}

func TestUpdatePost_ErrNotFound(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioNotFoundError("no post found")

	dao.On("UpdatePost", mock.AnythingOfType("*sioblog.BlogPost")).Return(updateResultFail, nil)

	resp, err := bs.UpdatePost(mockdata.PostEntity.ID, mockdata.UpdatePostRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestUpdatePost_UpdateErr(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioInternalServerError(
		"unexpected DB error: " + sql.ErrConnDone.Error(),
	)

	dao.On("UpdatePost", mock.AnythingOfType("*sioblog.BlogPost")).
		Return(updateResultFail, sql.ErrConnDone)

	resp, err := bs.UpdatePost(mockdata.PostEntity.ID, mockdata.UpdatePostRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestUpdateComment(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("UpdateComment", mock.AnythingOfType("*sioblog.BlogComment")).
		Return(updateResultHappy, nil)
	dao.On("GetCommentByID", mock.AnythingOfType("int64")).Return(mockdata.CommentEntity, nil)

	resp, err := bs.UpdateComment(mockdata.CommentEntity.ID, mockdata.UpdateCommentRequest)
	assert.NoError(t, err)
	assert.Equal(t, mockdata.CommentEntity.Content.String, resp.Content)
	dao.AssertExpectations(t)
}

func TestUpdateComment_ErrNotFound(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioNotFoundError("no comment found")

	dao.On("UpdateComment", mock.AnythingOfType("*sioblog.BlogComment")).
		Return(updateResultFail, nil)

	resp, err := bs.UpdateComment(mockdata.CommentEntity.ID, mockdata.UpdateCommentRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestUpdateComment_UpdateErr(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioInternalServerError(
		"unexpected DB error: " + sql.ErrConnDone.Error(),
	)

	dao.On("UpdateComment", mock.AnythingOfType("*sioblog.BlogComment")).
		Return(updateResultFail, sql.ErrConnDone)

	resp, err := bs.UpdateComment(mockdata.CommentEntity.ID, mockdata.UpdateCommentRequest)
	assert.Error(t, err)
	assert.Equal(t, testError.Error(), err.Error())
	assert.Nil(t, resp)
	dao.AssertExpectations(t)
}

func TestSoftDeletePost(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("SoftDeletePost", mock.AnythingOfType("*sioblog.BlogPost")).
		Return(updateResultHappy, nil)

	resp, err := bs.SoftDeletePost(mockdata.PostEntity.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
	dao.AssertExpectations(t)
}

func TestSoftDeletePost_NotFound(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioNotFoundError(constants.NO_POST_FOUND)

	dao.On("SoftDeletePost", mock.AnythingOfType("*sioblog.BlogPost")).
		Return(updateResultFail, nil)

	resp, err := bs.SoftDeletePost(mockdata.PostEntity.ID)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, testError.Error(), err.Error())
	dao.AssertExpectations(t)
}

func TestSoftDeletePost_DeleteErr(t *testing.T) {
	bs, dao := initEnv(t)
	testError := sioerror.NewSioInternalServerError(
		"unexpected DB error: " + sql.ErrConnDone.Error(),
	)

	dao.On("SoftDeletePost", mock.AnythingOfType("*sioblog.BlogPost")).
		Return(updateResultFail, sql.ErrConnDone)

	resp, err := bs.SoftDeletePost(mockdata.PostEntity.ID)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, testError.Error(), err.Error())
	dao.AssertExpectations(t)
}

func TestSoftDeleteComment(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("SoftDeleteComment", mock.AnythingOfType("*sioblog.BlogComment")).
		Return(updateResultHappy, nil)

	resp, err := bs.SoftDeleteComment(mockdata.CommentEntity.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
	dao.AssertExpectations(t)
}

func TestSoftDeleteComment_NotFound(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioNotFoundError("no comment found")

	dao.On("SoftDeleteComment", mock.AnythingOfType("*sioblog.BlogComment")).
		Return(updateResultFail, nil)

	resp, err := bs.SoftDeleteComment(mockdata.CommentEntity.ID)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, testError.Error(), err.Error())
	dao.AssertExpectations(t)
}

func TestSoftDeleteComment_DeleteErr(t *testing.T) {
	bs, dao := initEnv(t)
	testError := sioerror.NewSioInternalServerError(
		"unexpected DB error: " + sql.ErrConnDone.Error(),
	)

	dao.On("SoftDeleteComment", mock.AnythingOfType("*sioblog.BlogComment")).
		Return(updateResultFail, sql.ErrConnDone)

	resp, err := bs.SoftDeleteComment(mockdata.CommentEntity.ID)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, testError.Error(), err.Error())
	dao.AssertExpectations(t)
}
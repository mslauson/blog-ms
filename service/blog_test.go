package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitea.slauson.io/blog/blog-ms/dao/mocks"
	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
)

var (
	comment = &siogeneric.BlogComment{
		ID:           siogeneric.NewSioNullInt64(1),
		Content:      siogeneric.NewSioNullString("test"),
		DeletionDate: siogeneric.NewSioNullTime(time.Now()),
	}
	post = &siogeneric.BlogPost{
		ID:    siogeneric.NewSioNullInt64(1),
		Title: siogeneric.NewSioNullString("test"),
		Body:  siogeneric.NewSioNullString("test"),
	}
	posts    = &[]*siogeneric.BlogPost{post, post}
	comments = &[]*siogeneric.BlogComment{comment, comment}
)

func initEnv(t *testing.T) (*BlogSvc, *mocks.BlogDao) {
	mDao := mocks.NewBlogDao(t)
	svc := &BlogSvc{dao: mDao}
	return svc, mDao
}

func TestGetPost(t *testing.T) {
	bs, dao := initEnv(t)

	testID := int64(1)

	dao.On("GetPostByID", testID).Return(post, nil)

	resp, err := bs.GetPost(testID)
	assert.NoError(t, err)
	assert.Equal(t, post.Title.String, resp.Title)
	dao.AssertExpectations(t)
}

func TestGetPost_WithComments(t *testing.T) {
	bs, dao := initEnv(t)

	testID := int64(1)
	post.Comments = comments

	dao.On("GetPostByID", testID).Return(post, nil)

	resp, err := bs.GetPost(testID)
	assert.NoError(t, err)
	assert.Equal(t, post.Title.String, resp.Title)
	dao.AssertExpectations(t)
}

func TestGetPost_NotFound(t *testing.T) {
	bs, dao := initEnv(t)

	testID := int64(1)
	post.Comments = comments

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

	dao.On("GetAllPosts").Return(posts, nil)

	resp, err := bs.GetAllPosts()
	assert.NoError(t, err)
	assert.Equal(t, len(*posts), len(*resp))
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

	req := &dto.CreatePostRequest{Title: "test", Body: "test", CreatedByID: 1}

	dao.On("PostExists", req.Title, req.CreatedByID).Return(false, nil)
	dao.On("CreatePost", mock.AnythingOfType("*siogeneric.BlogPost")).Return(nil)

	resp, err := bs.CreatePost(req)
	assert.NoError(t, err)
	assert.Equal(t, post.Title.String, resp.Title)
	assert.Equal(t, post.Body.String, resp.Body)
	dao.AssertExpectations(t)
}

func TestAddComment(t *testing.T) {
	bs, dao := initEnv(t)

	req := &dto.AddCommentRequest{PostID: 1, UserID: 1, Content: "test"}

	dao.On("AddComment", mock.AnythingOfType("*siogeneric.BlogComment")).Return(nil)

	resp, err := bs.AddComment(req)
	assert.NoError(t, err)
	assert.Equal(t, comment.Content.String, resp.Content)
	dao.AssertExpectations(t)
}

func TestUpdatePost(t *testing.T) {
	bs, dao := initEnv(t)

	req := &dto.UpdatePostRequest{Title: "test"}

	dao.On("PostExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("UpdatePost", mock.AnythingOfType("*siogeneric.BlogPost")).Return(nil)
	dao.On("GetPostByID", mock.AnythingOfType("int64")).Return(post, nil)

	resp, err := bs.UpdatePost(post.ID.Int64, req)
	assert.NoError(t, err)
	assert.Equal(t, post.Title.String, resp.Title)
	dao.AssertExpectations(t)
}

func TestUpdateComment(t *testing.T) {
	bs, dao := initEnv(t)

	req := &dto.UpdateCommentRequest{Content: "test"}

	dao.On("CommentExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("UpdateComment", mock.AnythingOfType("*siogeneric.BlogComment")).Return(nil)
	dao.On("GetCommentByID", mock.AnythingOfType("int64")).Return(comment, nil)

	resp, err := bs.UpdateComment(comment.ID.Int64, req)
	assert.NoError(t, err)
	assert.Equal(t, comment.Content.String, resp.Content)
	dao.AssertExpectations(t)
}

func TestSoftDeletePost(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("PostExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("SoftDeletePost", mock.AnythingOfType("*siogeneric.BlogPost")).Return(nil)

	resp, err := bs.SoftDeletePost(post.ID.Int64)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
	dao.AssertExpectations(t)
}

func TestSoftDeletePost_NotFound(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioNotFoundError("post not found")

	dao.On("PostExistsByID", mock.AnythingOfType("int64")).Return(false, sql.ErrNoRows)

	resp, err := bs.SoftDeletePost(post.ID.Int64)
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

	dao.On("PostExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("SoftDeletePost", mock.AnythingOfType("*siogeneric.BlogPost")).
		Return(sql.ErrConnDone)

	resp, err := bs.SoftDeletePost(post.ID.Int64)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, testError.Error(), err.Error())
	dao.AssertExpectations(t)
}

func TestSoftDeleteComment(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("CommentExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("SoftDeleteComment", mock.AnythingOfType("*siogeneric.BlogComment")).Return(nil)

	resp, err := bs.SoftDeleteComment(comment.ID.Int64)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
	dao.AssertExpectations(t)
}

func TestSoftDeleteComment_NotFound(t *testing.T) {
	bs, dao := initEnv(t)

	testError := sioerror.NewSioNotFoundError("comment not found")

	dao.On("CommentExistsByID", mock.AnythingOfType("int64")).Return(false, sql.ErrNoRows)

	resp, err := bs.SoftDeleteComment(comment.ID.Int64)
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

	dao.On("CommentExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("SoftDeleteComment", mock.AnythingOfType("*siogeneric.BlogComment")).
		Return(sql.ErrConnDone)

	resp, err := bs.SoftDeleteComment(comment.ID.Int64)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, testError.Error(), err.Error())
	dao.AssertExpectations(t)
}

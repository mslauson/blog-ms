package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitea.slauson.io/blog/blog-ms/dao/mocks"
	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
)

var (
	comment = &siogeneric.BlogComment{
		ID:           siogeneric.NewSioNullInt64(1),
		DeletionDate: siogeneric.NewSioNullTime(time.Now()),
	}
	post = &siogeneric.BlogPost{
		ID:    siogeneric.NewSioNullInt64(1),
		Title: siogeneric.NewSioNullString("test1"),
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

func TestGetAllPosts(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("GetAllPosts").Return(posts, nil)

	resp, err := bs.GetAllPosts()
	assert.NoError(t, err)
	assert.Equal(t, len(*posts), len(*resp))
	dao.AssertExpectations(t)
}

func TestCreatePost(t *testing.T) {
	bs, dao := initEnv(t)

	req := &dto.CreatePostRequest{Title: "test", CreatedByID: 1}

	dao.On("PostExists", req.Title, req.CreatedByID).Return(false, nil)
	dao.On("CreatePost", mock.AnythingOfType("*dao.Post")).Return(nil)

	resp, err := bs.CreatePost(req)
	assert.NoError(t, err)
	assert.Equal(t, post.Title, resp.Title)
	dao.AssertExpectations(t)
}

func TestAddComment(t *testing.T) {
	bs, dao := initEnv(t)

	req := &dto.AddCommentRequest{PostID: 1, UserID: 1, Content: "test"}

	dao.On("AddComment", mock.AnythingOfType("*siogeneric.BlogComment")).Return(nil)

	resp, err := bs.AddComment(req)
	assert.NoError(t, err)
	assert.Equal(t, comment.Content, resp.Content)
	dao.AssertExpectations(t)
}

func TestUpdatePost(t *testing.T) {
	bs, dao := initEnv(t)

	req := &dto.UpdatePostRequest{Title: "test"}

	dao.On("PostExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("UpdatePost", mock.AnythingOfType("*dao.Post")).Return(nil)
	dao.On("GetPostByID", mock.AnythingOfType("int64")).Return(post, nil)

	resp, err := bs.UpdatePost(post.ID.Int64, req)
	assert.NoError(t, err)
	assert.Equal(t, post.Title, resp.Title)
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
	assert.Equal(t, comment.Content, resp.Content)
	dao.AssertExpectations(t)
}

func TestSoftDeletePost(t *testing.T) {
	bs, dao := initEnv(t)

	dao.On("CommentExistsByID", mock.AnythingOfType("int64")).Return(true, nil)
	dao.On("SoftDeleteComment", mock.AnythingOfType("*siogeneric.BlogComment")).Return(nil)

	resp, err := bs.SoftDeletePost(comment.ID.Int64)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
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

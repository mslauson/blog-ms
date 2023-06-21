package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/blog/blog-ms/service/mocks"
	"gitea.slauson.io/blog/blog-ms/testing/mockdata"
	"gitea.slauson.io/slausonio/go-testing/siotest"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func initEnv() (*BlogHdlr, *mocks.BlogService) {
	mSvc := &mocks.BlogService{}
	hdlr := &BlogHdlr{
		svc: mSvc,
	}
	return hdlr, mSvc
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name   string
		req    *dto.CreatePostRequest
		res    *dto.PostResponse
		status int
	}{
		{
			name: "success",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				Body:        "Test Body",
				CreatedByID: 1,
			},
			status: http.StatusOK,
		},
		{
			name: "Bad Title",
			req: &dto.CreatePostRequest{
				Title:       mockdata.LongTitle,
				Body:        "Test Body",
				CreatedByID: 1,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "No Title",
			req: &dto.CreatePostRequest{
				Body:        "Test Body",
				CreatedByID: 1,
			},
			status: http.StatusBadRequest,
		},

		{
			name: "Bad Body",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				Body:        mockdata.LongBody,
				CreatedByID: 1,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "No Body",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				CreatedByID: 1,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Missing CreatedByID",
			req: &dto.CreatePostRequest{
				Title: "Title",
				Body:  "Test Body",
			},
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdlr, mSvc := initEnv()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			siotest.MockJson(c, tt.req, http.MethodPost)

			if tt.status == http.StatusOK {
				mSvc.On("CreatePost", mock.AnythingOfType("*dto.CreatePostRequest")).
					Return(tt.res, nil)
			}

			hdlr.CreatePost(c)
			if tt.status == http.StatusOK {
				mSvc.On("CreatePost", tt.req).Return(tt.res, nil)
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil errors: %v", c.Errors)
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
			}
		})
	}
}

func TestCreatePost_SvcErr(t *testing.T) {
	req := &dto.CreatePostRequest{
		Title:       "Title",
		Body:        "Test Body",
		CreatedByID: 1,
	}
	hdlr, mSvc := initEnv()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	siotest.MockJson(c, req, http.MethodPost)
	mSvc.On("CreatePost", mock.AnythingOfType("*dto.CreatePostRequest")).
		Return(nil, errors.New("error"))

	hdlr.CreatePost(c)
	assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
}

func TestUpdatePost(t *testing.T) {
	tests := []struct {
		name   string
		ID     string
		req    *dto.UpdatePostRequest
		res    *dto.PostResponse
		status int
	}{
		{
			name: "success Update Both",
			req: &dto.UpdatePostRequest{
				Title:       "Title",
				Body:        "Test Body",
				UpdatedByID: 1,
			},
			status: http.StatusOK,
			ID:     "1",
		},
		{
			name: "success Update Title",
			req: &dto.UpdatePostRequest{
				Title:       "Title",
				UpdatedByID: 1,
			},
			status: http.StatusOK,
			ID:     "1",
		},
		{
			name: "success Update Body",
			req: &dto.UpdatePostRequest{
				Body:        "Test Body",
				UpdatedByID: 1,
			},
			status: http.StatusOK,
			ID:     "1",
		},
		{
			name: "Bad ID",
			req: &dto.UpdatePostRequest{
				Body:        "Test Body",
				UpdatedByID: 1,
			},
			status: http.StatusBadRequest,
			ID:     "1sdfsdfe",
		},
		{
			name: "Bad Title",
			req: &dto.UpdatePostRequest{
				Title:       mockdata.LongTitle,
				Body:        "Test Body",
				UpdatedByID: 1,
			},
			status: http.StatusBadRequest,
			ID:     "1",
		},
		{
			name: "Bad Body",
			req: &dto.UpdatePostRequest{
				Title:       "Title",
				Body:        mockdata.LongBody,
				UpdatedByID: 1,
			},
			status: http.StatusBadRequest,
			ID:     "1",
		},
		{
			name: "Bad - No updates passed",
			req: &dto.UpdatePostRequest{
				UpdatedByID: 1,
			},
			status: http.StatusBadRequest,
			ID:     "1",
		},
		{
			name: "Missing UpdatedByID",
			req: &dto.UpdatePostRequest{
				Title: "Title",
				Body:  "Test Body",
			},
			status: http.StatusBadRequest,
			ID:     "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdlr, mSvc := initEnv()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			siotest.MockJson(c, tt.req, http.MethodPatch)
			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.ID}}
			if tt.status == http.StatusOK {
				mSvc.On("UpdatePost", mock.AnythingOfType("int64"), mock.AnythingOfType("*dto.UpdatePostRequest")).
					Return(tt.res, nil)
			}

			hdlr.UpdatePost(c)
			if tt.status == http.StatusOK {
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil errors: %v", c.Errors)
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
			}
		})
	}
}

func TestUpdatePost_SvcErr(t *testing.T) {
	req := &dto.UpdatePostRequest{
		Title:       "Title",
		Body:        "Test Body",
		UpdatedByID: 1,
	}
	hdlr, mSvc := initEnv()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	siotest.MockJson(c, req, http.MethodPatch)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	mSvc.On("UpdatePost", mock.AnythingOfType("int64"), mock.AnythingOfType("*dto.UpdatePostRequest")).
		Return(nil, errors.New("error"))

	hdlr.UpdatePost(c)
	assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
}

func TestAddComment(t *testing.T) {
	tests := []struct {
		name   string
		req    *dto.AddCommentRequest
		res    *dto.CommentResponse
		status int
	}{
		{
			name: "success",
			req: &dto.AddCommentRequest{
				Content: "Test Content",
				UserID:  1,
				PostID:  1,
			},
			status: http.StatusOK,
		},
		{
			name: "Bad Content",
			req: &dto.AddCommentRequest{
				Content: mockdata.LongComment,
				UserID:  1,
				PostID:  1,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Missing Content",
			req: &dto.AddCommentRequest{
				UserID: 1,
				PostID: 1,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Missing UserID",
			req: &dto.AddCommentRequest{
				Content: "Test Content",
				PostID:  1,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Missing PostID",
			req: &dto.AddCommentRequest{
				Content: "Test Content",
				UserID:  1,
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdlr, mSvc := initEnv()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			siotest.MockJson(c, tt.req, http.MethodPost)

			if tt.status == http.StatusOK {
				mSvc.On("AddComment", mock.AnythingOfType("*dto.AddCommentRequest")).
					Return(tt.res, nil)
			}

			hdlr.AddComment(c)
			if tt.status == http.StatusOK {
				mSvc.On("AddComment", tt.req).Return(tt.res, nil)
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil errors: %v", c.Errors)
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
			}
		})
	}
}

func TestAddComment_SvcErr(t *testing.T) {
	req := &dto.AddCommentRequest{
		Content: "Test Content",
		UserID:  1,
		PostID:  1,
	}

	hdlr, mSvc := initEnv()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	siotest.MockJson(c, req, http.MethodPost)
	mSvc.On("AddComment", mock.AnythingOfType("*dto.AddCommentRequest")).
		Return(nil, errors.New("error"))

	hdlr.AddComment(c)
	assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
}

func TestUpdateComment(t *testing.T) {
	tests := []struct {
		name   string
		ID     string
		req    *dto.UpdateCommentRequest
		res    *dto.CommentResponse
		status int
	}{
		{
			name: "success",
			req: &dto.UpdateCommentRequest{
				Content: "Test Content",
			},
			status: http.StatusOK,
			ID:     "1",
		},
		{
			name: "success Bad ID",
			req: &dto.UpdateCommentRequest{
				Content: "Test Content",
			},
			status: http.StatusBadRequest,
			ID:     "1asdf",
		},
		{
			name: "success Bad Content",
			req: &dto.UpdateCommentRequest{
				Content: mockdata.LongComment,
			},
			status: http.StatusBadRequest,
			ID:     "1",
		},
		{
			name:   "success No Content",
			req:    &dto.UpdateCommentRequest{},
			status: http.StatusBadRequest,
			ID:     "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdlr, mSvc := initEnv()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			siotest.MockJson(c, tt.req, http.MethodPatch)
			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.ID}}
			if tt.status == http.StatusOK {
				mSvc.On("UpdateComment", mock.AnythingOfType("int64"), mock.AnythingOfType("*dto.UpdateCommentRequest")).
					Return(tt.res, nil)
			}

			hdlr.UpdateComment(c)
			if tt.status == http.StatusOK {
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil errors: %v", c.Errors)
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
			}
		})
	}
}

func TestUpdateComment_SvcErr(t *testing.T) {
	req := &dto.UpdateCommentRequest{
		Content: "Title",
	}
	hdlr, mSvc := initEnv()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	siotest.MockJson(c, req, http.MethodPatch)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	mSvc.On("UpdateComment", mock.AnythingOfType("int64"), mock.AnythingOfType("*dto.UpdateCommentRequest")).
		Return(nil, errors.New("error"))

	hdlr.UpdateComment(c)
	assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		name    string
		request string
		status  int
		result  *dto.PostResponse
	}{
		{name: "happy", request: "123", status: http.StatusOK, result: mockdata.PostResponse},
		{name: "BadId", request: "1s23", status: http.StatusBadRequest, result: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdlr, mSvc := initEnv()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.request}}
			if tt.status == http.StatusOK {
				i, _ := strconv.ParseInt(tt.request, 10, 64)
				mSvc.On("GetPost", i).
					Return(tt.result, nil)
			}

			hdlr.GetPost(c)
			if tt.status == http.StatusOK {
				assert.Nilf(t, c.Errors, "c.Errors should be nil")
			} else {
				assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestGetPost_SvcErr(t *testing.T) {
	hdlr, mSvc := initEnv()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	i, _ := strconv.ParseInt("1", 10, 64)
	mSvc.On("GetPost", i).
		Return(nil, errors.New("error"))

	hdlr.GetPost(c)
	assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
}

func TestGetAllPosts(t *testing.T) {
	hdlr, mSvc := initEnv()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	mSvc.On("GetAllPosts").
		Return(mockdata.PostResponses, nil)

	hdlr.GetAllPosts(c)
	assert.Nilf(t, c.Errors, "c.Errors should be nil")
}

func TestGetAllPosts_SvcErr(t *testing.T) {
	hdlr, mSvc := initEnv()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	mSvc.On("GetAllPosts").
		Return(nil, errors.New("error"))

	hdlr.GetAllPosts(c)
	assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
}

func TestSoftDeletePost(t *testing.T) {
	tests := []struct {
		name    string
		request string
		status  int
		result  *siogeneric.SuccessResponse
	}{
		{
			name:    "happy",
			request: "123",
			status:  http.StatusOK,
			result:  mockdata.SuccessResponseSuccess,
		},
		{name: "BadId", request: "1s23", status: http.StatusBadRequest, result: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdlr, mSvc := initEnv()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.request}}
			if tt.status == http.StatusOK {
				i, _ := strconv.ParseInt(tt.request, 10, 64)
				mSvc.On("SoftDeletePost", i).
					Return(tt.result, nil)
			}

			hdlr.SoftDeletePost(c)
			if tt.status == http.StatusOK {
				assert.Nilf(t, c.Errors, "c.Errors should be nil")
			} else {
				assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestSoftDeletePost_SvcErr(t *testing.T) {
	hdlr, mSvc := initEnv()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	i, _ := strconv.ParseInt("1", 10, 64)
	mSvc.On("SoftDeletePost", i).
		Return(nil, errors.New("error"))

	hdlr.SoftDeletePost(c)
	assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
}

func TestSoftDeleteComment(t *testing.T) {
	tests := []struct {
		name    string
		request string
		status  int
		result  *siogeneric.SuccessResponse
	}{
		{
			name:    "happy",
			request: "123",
			status:  http.StatusOK,
			result:  mockdata.SuccessResponseSuccess,
		},
		{name: "BadId", request: "1s23", status: http.StatusBadRequest, result: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdlr, mSvc := initEnv()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.request}}
			if tt.status == http.StatusOK {
				i, _ := strconv.ParseInt(tt.request, 10, 64)
				mSvc.On("SoftDeleteComment", i).
					Return(tt.result, nil)
			}

			hdlr.SoftDeleteComment(c)
			if tt.status == http.StatusOK {
				assert.Nilf(t, c.Errors, "c.Errors should be nil")
			} else {
				assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestSoftDeleteComment_SvcErr(t *testing.T) {
	hdlr, mSvc := initEnv()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	i, _ := strconv.ParseInt("1", 10, 64)
	mSvc.On("SoftDeleteComment", i).
		Return(nil, errors.New("error"))

	hdlr.SoftDeleteComment(c)
	assert.NotNilf(t, c.Errors, "c.Errors shouldnt be nil")
}

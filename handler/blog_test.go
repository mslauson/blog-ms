package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/blog/blog-ms/service/mocks"
	"gitea.slauson.io/blog/blog-ms/testing/mockdata"
	"gitea.slauson.io/slausonio/go-testing/siotest"
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

package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/blog/blog-ms/service/mocks"
	"gitea.slauson.io/blog/blog-ms/testing/mockdata"
	"gitea.slauson.io/slausonio/go-testing/siotest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
			status: http.StatusOK,
		},
		{
			name: "No Title",
			req: &dto.CreatePostRequest{
				Body:        "Test Body",
				CreatedByID: 1,
			},
			status: http.StatusOK,
		},

		{
			name: "Bad Body",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				Body:        mockdata.LongBody,
				CreatedByID: 1,
			},
			status: http.StatusOK,
		},
		{
			name: "No Body",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				CreatedByID: 1,
			},
			status: http.StatusOK,
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

			mSvc.On("CreatePost", tt.req).Return(tt.res, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			siotest.MockJson(c, tt.req, http.MethodGet)
			hdlr.CreatePost(c)
			if tt.status == http.StatusOK {
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil errors: %v", c.Errors)
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
			}
		})
	}
}

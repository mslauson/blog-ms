package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/blog/blog-ms/service/mocks"
	"gitea.slauson.io/slausonio/go-testing/siotest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var(
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
				Title:       "18d6cd35762e7155d9a4862f36e127b2",
				Body:        "eb7c2688988d2f9380250b167de24360",
				CreatedByID: 1,
			},
			status: http.StatusOK,
		},
		{
			name: "Bad Title",
			req: &dto.CreatePostRequest{
				Title:       "18d6cd35762e7155d9a4862f36e127b2",
				Body:        "eb7c2688988d2f9380250b167de24360",
				CreatedByID: 1,
			},
			status: http.StatusOK,
		},
		{
			name: "Missing CreatedByID",
			req: &dto.CreatePostRequest{
				Title:       "18d6cd35762e7155d9a4862f36e127b2",
				Body:        "eb7c2688988d2f9380250b167de24360",
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
			assert.Equal(t, tt.status, w.Code)
			if tt.status == http.StatusOK {
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil errors: %v", c.Errors)
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil errors: %v", c.Errors)
			}
		})
	}
}

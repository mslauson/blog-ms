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
				Body:        "40df6f20a48a0af47e72eec08dc622a",
				CreatedByID: 1,
			},
			status: http.StatusOK,
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
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
			}
		})
	}
}

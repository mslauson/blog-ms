package handler

import (
	"gitea.slauson.io/blog/blog-ms/service"
	"github.com/gin-gonic/gin"
)

type PostHdlr struct {
	svc service.BlogService
}

type PostHandler interface {
	GetPost(c *gin.Context)
	GetAllPosts(c *gin.Context)
	CreatePost(c *gin.Context)
	AddComment(c *gin.Context)
	UpdatePost(c *gin.Context)
	UpdateComment(c *gin.Context)
	SoftDeletePost(c *gin.Context)
	SoftDeleteComment(c *gin.Context)
}

func NewPostHdlr() *PostHdlr {
	return &PostHdlr{
		svc: service.NewBlogSvc(),
	}
}

func (PostHdlr) GetPost(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (PostHdlr) GetAllPosts(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (PostHdlr) CreatePost(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (PostHdlr) AddComment(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (PostHdlr) UpdatePost(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (PostHdlr) UpdateComment(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (PostHdlr) SoftDeletePost(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (PostHdlr) SoftDeleteComment(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

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

// GET
// @Summary Get post by ID
// @Description Get post by ID
// @Tags post
// @Accept  json
// @Produce  json
// @Param id path string true "Post id"
// @Success 200 {object} dto.PostResponse
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/post/v1/:id [get]
func (PostHdlr) GetPost(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// GET
// @Summary Get all posts
// @Description Get all posts
// @Tags post
// @Accept  json
// @Produce  json
// @Success 200 {array} dto.PostResponse
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/post/v1/ [get]
func (PostHdlr) GetAllPosts(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// POST
// @Summary Create Post
// @Description Create Post
// @Tags post
// @Accept  json
// @Produce  json
// @Param post body dto.CreatePostRequest true "Create Post request"
// @Success 200 {object} dto.PostResponse
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/post/v1/ [post]
func (PostHdlr) CreatePost(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// POST
// @Summary Add Comment
// @Description Add Comment
// @Tags comment
// @Accept  json
// @Produce  json
// @Param comment body dto.AddCommentRequest true "Add Comment request"
// @Success 200 {object} dto.CommentResponse
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/post/v1/comment [post]
func (PostHdlr) AddComment(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// PATCH
// @Summary Update Post
// @Description Update Post
// @Tags post
// @Accept  json
// @Produce  json
// @Param id path string true "Post id"
// @Param post body dto.UpdatePostRequest true "Update Post request"
// @Success 200 {object} dto.PostResponse
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/post/v1/:id [patch]
func (PostHdlr) UpdatePost(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// PATCH
// @Summary Update Comment
// @Description Update Comment
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path string true "Comment id"
// @Param comment body dto.UpdateCommentRequest true "Update Comment request"
// @Success 200 {object} dto.CommentResponse
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/post/v1/comment/:id [patch]
func (PostHdlr) UpdateComment(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// DELETE
// @Summary Delete Post
// @Description Delete Post
// @Tags post
// @Accept  json
// @Produce  json
// @Param id path string true "Post id"
// @Success 200 {object} dto.PostResponse
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/post/v1/:id [delete]
func (PostHdlr) SoftDeletePost(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// DELETE
// @Summary Delete Comment
// @Description Delete Comment
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path string true "Comment id"
// @Success 200 {object} dto.CommentResponse
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/post/v1/comment/:id [delete]
func (PostHdlr) SoftDeleteComment(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

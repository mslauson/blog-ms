package handler

import (
	"net/http"

	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/blog/blog-ms/service"
	"gitea.slauson.io/blog/blog-ms/utils"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
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
func (ph *PostHdlr) GetPost(c *gin.Context) {
	id := c.Param("id")

	iId, err := sioUtils.ConvertToInt64(id)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}
	if result, err := ph.svc.GetPost(iId); result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		_ = c.Error(err)
	}
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
func (ph *PostHdlr) GetAllPosts(c *gin.Context) {
	// includeDeleted := c.DefaultQuery("includeDeleted", "false")
	// idBool, err := sioUtils.ToBool(includeDeleted)
	// if err != nil {
	// 	_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
	// 	return
	// }

	if result, err := ph.svc.GetAllPosts(); result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		_ = c.Error(err)
	}
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
func (ph *PostHdlr) CreatePost(c *gin.Context) {
	validations := utils.NewValidator()
	var request dto.CreatePostRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}

	if err := validations.ValidateCreatePostRequest(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if result, err := ph.svc.CreatePost(&request); result == nil {
		_ = c.Error(err)
	} else {
		c.JSON(http.StatusOK, result)
	}
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
func (ph *PostHdlr) AddComment(c *gin.Context) {
	validations := utils.NewValidator()
	var request dto.AddCommentRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}

	if err := validations.ValidateAddCommentRequest(&request); err != nil {
		_ = c.Error(err)
		return
	}

	if result, err := ph.svc.AddComment(&request); result == nil {
		_ = c.Error(err)
	} else {
		c.JSON(http.StatusOK, result)
	}
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
func (ph *PostHdlr) UpdatePost(c *gin.Context) {
	validations := utils.NewValidator()

	id := c.Param("id")

	var request dto.UpdatePostRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}

	if err := validations.ValidateUpdatePostRequest(&request); err != nil {
		_ = c.Error(err)
		return
	}

	iId, err := sioUtils.ConvertToInt64(id)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}

	if result, err := ph.svc.UpdatePost(iId, &request); result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		_ = c.Error(err)
	}
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
func (ph *PostHdlr) UpdateComment(c *gin.Context) {
	validations := utils.NewValidator()

	id := c.Param("id")

	var request dto.UpdateCommentRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}

	if err := validations.ValidateUpdateCommentRequest(&request); err != nil {
		_ = c.Error(err)
		return
	}

	iId, err := sioUtils.ConvertToInt64(id)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}

	if result, err := ph.svc.UpdateComment(iId, &request); result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		_ = c.Error(err)
	}
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
func (ph *PostHdlr) SoftDeletePost(c *gin.Context) {
	id := c.Param("id")

	iId, err := sioUtils.ConvertToInt64(id)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}
	if result, err := ph.svc.SoftDeletePost(iId); result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		_ = c.Error(err)
	}
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
func (ph *PostHdlr) SoftDeleteComment(c *gin.Context) {
	id := c.Param("id")

	iId, err := sioUtils.ConvertToInt64(id)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}
	if result, err := ph.svc.SoftDeleteComment(iId); result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		_ = c.Error(err)
	}
}

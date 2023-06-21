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

type BlogHdlr struct {
	svc service.BlogService
}

type BlogHandler interface {
	GetPost(c *gin.Context)
	GetAllPosts(c *gin.Context)
	CreatePost(c *gin.Context)
	AddComment(c *gin.Context)
	UpdatePost(c *gin.Context)
	UpdateComment(c *gin.Context)
	SoftDeletePost(c *gin.Context)
	SoftDeleteComment(c *gin.Context)
}

func NewBlogHdlr() *BlogHdlr {
	return &BlogHdlr{
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
// @Router /api/blog/v1/post/:id [get]
func (ph *BlogHdlr) GetPost(c *gin.Context) {
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
// @Router /api/blog/v1/post/ [get]
func (ph *BlogHdlr) GetAllPosts(c *gin.Context) {
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
// @Router /api/blog/v1/post/ [post]
func (ph *BlogHdlr) CreatePost(c *gin.Context) {
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

	if result, err := ph.svc.CreatePost(&request); err != nil {
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
// @Router /api/blog/v1/post/comment [post]
func (ph *BlogHdlr) AddComment(c *gin.Context) {
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

	if result, err := ph.svc.AddComment(&request); err != nil {
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
// @Router /api/blog/v1/post/:id [patch]
func (ph *BlogHdlr) UpdatePost(c *gin.Context) {
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

	if result, err := ph.svc.UpdatePost(iId, &request); err != nil {
		_ = c.Error(err)
	} else {
		c.JSON(http.StatusOK, result)
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
// @Router /api/blog/v1/post/comment/:id [patch]
func (ph *BlogHdlr) UpdateComment(c *gin.Context) {
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
// @Router /api/blog/v1/post/:id [delete]
func (ph *BlogHdlr) SoftDeletePost(c *gin.Context) {
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
// @Router /api/blog/v1/post/comment/:id [delete]
func (ph *BlogHdlr) SoftDeleteComment(c *gin.Context) {
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

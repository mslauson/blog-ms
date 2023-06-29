package service

import (
	"time"

	"gitea.slauson.io/blog/post-ms/constants"
	"gitea.slauson.io/blog/post-ms/dao"
	"gitea.slauson.io/blog/post-ms/dto"
	"gitea.slauson.io/slausonio/go-types/sioblog"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/siodao"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
)

type PostSvc struct {
	dao dao.PostDao
}

// PostService interface
//
//go:generate mockery --name PostService
type PostService interface {
	GetPost(id int64) (*dto.PostResponse, error)
	GetAllPosts() (*[]*dto.PostResponse, error)
	CreatePost(req *dto.CreatePostRequest) (*dto.PostResponse, error)
	AddComment(req *dto.AddCommentRequest) (*dto.CommentResponse, error)
	UpdatePost(ID int64, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	UpdateComment(ID int64, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	SoftDeletePost(ID int64) (*siogeneric.SuccessResponse, error)
	SoftDeleteComment(ID int64) (*siogeneric.SuccessResponse, error)
}

func NewPostSvc() *PostSvc {
	return &PostSvc{
		dao: dao.NewBlogDao(),
	}
}

func (ps *PostSvc) GetPost(id int64) (*dto.PostResponse, error) {
	post, err := ps.dao.GetPostByID(id)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	return buildPostResponse(post), nil
}

func (ps *PostSvc) GetAllPosts() (*[]*dto.PostResponse, error) {
	posts, err := ps.dao.GetAllPosts()
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.POSTS)
	}

	if posts == nil || len(*posts) == 0 {
		return nil, sioerror.NewSioNotFoundError(constants.NO_POSTS_FOUND)
	}

	return buildAllPostsResponse(posts), nil
}

func (ps *PostSvc) CreatePost(req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	if exists, err := ps.dao.PostExists(req.Title, req.CreatedByID); err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	} else if exists {
		return nil, sioerror.NewSioBadRequestError(constants.POST_EXISTS)
	}

	post := buildCreatePostEntity(req)
	if err := ps.dao.CreatePost(post); err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	return buildPostResponse(post), nil
}

func (ps *PostSvc) AddComment(req *dto.AddCommentRequest) (*dto.CommentResponse, error) {
	if err := ps.postExistsByID(req.PostID); err != nil {
		return nil, err
	}

	comment := buildAddCommentEntity(req)
	if err := ps.dao.AddComment(comment); err != nil {
		return nil, siodao.HandleDbErr(err, constants.COMMENT)
	}

	return buildCommentResponse(comment), nil
}

func (ps *PostSvc) UpdatePost(ID int64, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post := buildUpdatePostEntity(ID, req)

	result, err := ps.dao.UpdatePost(post)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	n, err := result.RowsAffected()
	if n == 0 {
		return nil, sioerror.NewSioNotFoundError(constants.NO_POST_FOUND)
	}

	post, err = ps.dao.GetPostByID(ID)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	return buildPostResponse(post), nil
}

func (ps *PostSvc) UpdateComment(
	ID int64, req *dto.UpdateCommentRequest,
) (*dto.CommentResponse, error) {
	comment := buildUpdateCommentEntity(ID, req)

	result, err := ps.dao.UpdateComment(comment)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.COMMENT)
	}

	n, err := result.RowsAffected()
	if n == 0 {
		return nil, sioerror.NewSioNotFoundError(constants.NO_COMMENT_FOUND)
	}

	comment, err = ps.dao.GetCommentByID(ID)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.COMMENT)
	}

	return buildCommentResponse(comment), nil
}

func (ps *PostSvc) SoftDeletePost(ID int64) (*siogeneric.SuccessResponse, error) {
	post := &sioblog.BlogPost{
		ID:           ID,
		DeletionDate: siodao.BuildNullTime(time.Now()),
	}

	result, err := ps.dao.SoftDeletePost(post)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	n, err := result.RowsAffected()
	if n == 0 {
		return nil, sioerror.NewSioNotFoundError(constants.NO_POST_FOUND)
	}

	return &siogeneric.SuccessResponse{Success: true}, nil
}

func (ps *PostSvc) SoftDeleteComment(ID int64) (*siogeneric.SuccessResponse, error) {
	comment := &sioblog.BlogComment{
		ID:           ID,
		DeletionDate: siodao.BuildNullTime(time.Now()),
	}

	result, err := ps.dao.SoftDeleteComment(comment)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.COMMENT)
	}

	n, err := result.RowsAffected()
	if n == 0 {
		return nil, sioerror.NewSioNotFoundError(constants.NO_COMMENT_FOUND)
	}

	return &siogeneric.SuccessResponse{Success: true}, nil
}

func (ps *PostSvc) postExistsByID(id int64) error {
	if exists, err := ps.dao.PostExistsByID(id); err != nil {
		return siodao.HandleDbErr(err, constants.POST)
	} else if !exists {
		return sioerror.NewSioNotFoundError(constants.NO_POST_FOUND)
	}

	return nil
}

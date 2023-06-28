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

type BlogSvc struct {
	dao dao.BlogDao
}

// BlogService interface
//
//go:generate mockery --name BlogService
type BlogService interface {
	GetPost(id int64) (*dto.PostResponse, error)
	GetAllPosts() (*[]*dto.PostResponse, error)
	CreatePost(req *dto.CreatePostRequest) (*dto.PostResponse, error)
	AddComment(req *dto.AddCommentRequest) (*dto.CommentResponse, error)
	UpdatePost(ID int64, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	UpdateComment(ID int64, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	SoftDeletePost(ID int64) (*siogeneric.SuccessResponse, error)
	SoftDeleteComment(ID int64) (*siogeneric.SuccessResponse, error)
}

func NewBlogSvc() *BlogSvc {
	return &BlogSvc{
		dao: dao.NewBlogDao(),
	}
}

func (bs *BlogSvc) GetPost(id int64) (*dto.PostResponse, error) {
	post, err := bs.dao.GetPostByID(id)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	return buildPostResponse(post), nil
}

func (bs *BlogSvc) GetAllPosts() (*[]*dto.PostResponse, error) {
	post, err := bs.dao.GetAllPosts()
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.POSTS)
	}

	return buildAllPostsResponse(post), nil
}

func (bs *BlogSvc) CreatePost(req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	if exists, err := bs.dao.PostExists(req.Title, req.CreatedByID); err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	} else if exists {
		return nil, sioerror.NewSioBadRequestError(constants.POST_EXISTS)
	}

	post := buildCreatePostEntity(req)
	if err := bs.dao.CreatePost(post); err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	return buildPostResponse(post), nil
}

func (bs *BlogSvc) AddComment(req *dto.AddCommentRequest) (*dto.CommentResponse, error) {
	if err := bs.postExistsByID(req.PostID); err != nil {
		return nil, err
	}

	comment := buildAddCommentEntity(req)
	if err := bs.dao.AddComment(comment); err != nil {
		return nil, siodao.HandleDbErr(err, constants.COMMENT)
	}

	return buildCommentResponse(comment), nil
}

func (bs *BlogSvc) UpdatePost(ID int64, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	if err := bs.postExistsByID(ID); err != nil {
		return nil, err
	}

	post := buildUpdatePostEntity(req)
	if err := bs.dao.UpdatePost(post); err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	post, err := bs.dao.GetPostByID(ID)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	return buildPostResponse(post), nil
}

func (bs *BlogSvc) UpdateComment(
	ID int64, req *dto.UpdateCommentRequest,
) (*dto.CommentResponse, error) {
	if err := bs.commentExistsByID(ID); err != nil {
		return nil, err
	}

	comment := buildUpdateCommentEntity(req)
	if err := bs.dao.UpdateComment(comment); err != nil {
		return nil, siodao.HandleDbErr(err, constants.COMMENT)
	}

	comment, err := bs.dao.GetCommentByID(ID)
	if err != nil {
		return nil, siodao.HandleDbErr(err, constants.COMMENT)
	}

	return buildCommentResponse(comment), nil
}

func (bs *BlogSvc) SoftDeletePost(ID int64) (*siogeneric.SuccessResponse, error) {
	if err := bs.postExistsByID(ID); err != nil {
		return nil, err
	}

	post := new(sioblog.BlogPost)
	post.DeletionDate = time.Now()

	if err := bs.dao.SoftDeletePost(post); err != nil {
		return nil, siodao.HandleDbErr(err, constants.POST)
	}

	return &siogeneric.SuccessResponse{Success: true}, nil
}

func (bs *BlogSvc) SoftDeleteComment(ID int64) (*siogeneric.SuccessResponse, error) {
	if err := bs.commentExistsByID(ID); err != nil {
		return nil, err
	}

	comment := new(sioblog.BlogComment)
	comment.DeletionDate = time.Now()

	if err := bs.dao.SoftDeleteComment(comment); err != nil {
		return nil, siodao.HandleDbErr(err, constants.COMMENT)
	}

	return &siogeneric.SuccessResponse{Success: true}, nil
}

func (bs *BlogSvc) postExistsByID(id int64) error {
	if exists, err := bs.dao.PostExistsByID(id); err != nil {
		return siodao.HandleDbErr(err, constants.POST)
	} else if !exists {
		return sioerror.NewSioNotFoundError(constants.NO_POST_FOUND)
	}

	return nil
}

func (bs *BlogSvc) commentExistsByID(id int64) error {
	if exists, err := bs.dao.CommentExistsByID(id); err != nil {
		return siodao.HandleDbErr(err, constants.COMMENT)
	} else if !exists {
		return sioerror.NewSioNotFoundError(constants.NO_COMMENT_FOUND)
	}

	return nil
}

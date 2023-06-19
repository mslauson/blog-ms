package service

import (
	"gitea.slauson.io/blog/blog-ms/dao"
	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/slausonio/go-utils/siodao"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
)

type BlogSvc struct {
	dao dao.PostDao
}

// BlogService interface
type BlogService interface {
	GetPost(id int64) (*dto.PostResponse, error)
	GetAllPosts() (*[]*dto.PostResponse, error)
	CreatePost(req *dto.CreatePostRequest) (*dto.PostResponse, error)
	AddComment(req *dto.AddCommentRequest) (*dto.CommentResponse, error)
	UpdatePost(req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	UpdateComment(req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	SoftDeletePost(id int64) (*dto.PostResponse, error)
	SoftDeleteComment(id int64) (*dto.CommentResponse, error)
}

// Newbs function
func NewBlogSvc() *BlogSvc {
	return &BlogSvc{
		dao: dao.NewPostDao(),
	}
}

func (bs *BlogSvc) GetPost(id int64) (*dto.PostResponse, error) {
	post, err := bs.dao.GetPostByID(id)
	if err != nil {
		return nil, siodao.HandleDbErr(err, POST)
	}

	return buildPostResponse(post), nil
}

func (bs *BlogSvc) GetAllPosts() (*[]*dto.PostResponse, error) {
	post, err := bs.dao.GetAllPosts()
	if err != nil {
		return nil, siodao.HandleDbErr(err, POST)
	}

	return buildAllPostsResponse(post), nil
}

func (bs *BlogSvc) CreatePost(req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	if exists, err := bs.dao.Exists(req.Title, req.CreatedByID); err != nil {
		return nil, siodao.HandleDbErr(err, POST)
	} else if exists {
		return nil, sioerror.NewSioBadRequestError(POST_EXISTS)
	}

	post := buildCreatePostEntity(req)
	if err := bs.dao.CreatePost(post); err != nil {
		return nil, siodao.HandleDbErr(err, POST)
	}

	return buildPostResponse(post), nil
}

func (bs *BlogSvc) AddComment(req *dto.AddCommentRequest) (*dto.CommentResponse, error) {
	comment := buildAddCommentEntity(req)
	if err := bs.dao.AddComment(comment); err != nil {
		return nil, siodao.HandleDbErr(err, COMMENT)
	}

	return buildCommentResponse(comment), nil
}

func (bs *BlogSvc) UpdatePost(req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	panic("not implemented") // TODO: Implement
}


func (bs *BlogSvc) UpdateComment(
	req *dto.UpdateCommentRequest,
) (*dto.CommentResponse, error) {
}

func (bs *BlogSvc) SoftDeletePost(id int64) (*dto.PostResponse, error) {
}
func (bs *BlogSvc) SoftDeleteComment(id int64) (*dto.CommentResponse, error) {
}

func (bs *BlogSvc) postExistsByID(id int64) (bool, error) {
	return bs.dao.ExistsByID(id)
}

package service

import (
	"gitea.slauson.io/blog/blog-ms/dao"
	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/slausonio/go-utils/siodao"
)

const (
	POST    = "post"
	COMMENT = "comment"
)

type BlogSvc struct {
	dao dao.PostDao
}

// BlogService interface
type BlogService interface {
	GetPost(id int64) (*dto.PostResponse, error)
	GetAllPosts() (*[]*dto.PostResponse, error)
	CreatePost(post *dto.CreatePostRequest) (*dto.PostResponse, error)
	AddComment(comment *dto.AddCommentRequest) (*dto.CommentResponse, error)
	UpdatePost(post *dto.UpdatePostRequest) (*dto.PostResponse, error)
	UpdateComment(comment *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	SoftDeletePost(id int64) (*dto.PostResponse, error)
	SoftnDeleteComment(id int64) (*dto.CommentResponse, error)
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

func (bs *BlogSvc) CreatePost(post *dto.CreatePostRequest) (*dto.PostResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (bs *BlogSvc) AddComment(comment *dto.AddCommentRequest) (*dto.CommentResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (bs *BlogSvc) UpdatePost(post *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (bs *BlogSvc) SoftDeletePost(id int64) (*dto.PostResponse, error) {
}

func (bs *BlogSvc) UpdateComment(
	comment *dto.UpdateCommentRequest,
) (*dto.CommentResponse, error) {
}

func (bs *BlogSvc) SoftenDeleteComment(id int64) (*dto.CommentResponse, error) {
}

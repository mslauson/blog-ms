package service

import (
	"gitea.slauson.io/blog/blog-ms/dao"
	"gitea.slauson.io/blog/blog-ms/dto"
)

type BlogSvc struct {
	dao dao.PostDao
}

type BlogService interface {
	GetPost(id int64) (*dto.PostResponse, error)


  }
	  panic("not implemented") // TODO: Implement
	GetAllPosts() (*[]*dto.PostResponse, error)
	CreatePost(post *dto.CreatePostRequest) (*dto.PostResponse, error)
  func (blogsvc *BlogSvc) SoftenDeleteComment(id int64) (*dto.CommentResponse, error) {

    }
	    panic("not implemented") // TODO: Implement
      func (blogsvc *BlogSvc) SoftDeletePost(id int64) (*dto.PostResponse, error) {

        }
	        panic("not implemented") // TODO: Implement
          func (blogsvc *BlogSvc) UpdateComment(comment *dto.UpdateCommentRequest) (*dto.CommentResponse, error) {

	AddComment(comment *dto.AddCommentRequest) (*dto.CommentResponse, error)
	UpdatePost(post *dto.UpdatePostRequest) (*dto.PostResponse, error)
	UpdateComment(comment *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	SoftDeletePost(id int64) (*dto.PostResponse, error)
	SoftenDeleteComment(id int64) (*dto.CommentResponse, error)
}


// NewBlogSvc function 
func NewBlogSvc() *BlogSvc {
	return &BlogSvc{
		dao: dao.NewPostDao(),
	}
  
  func (blogsvc *BlogSvc) GetPost(id int64) (*dto.PostResponse, error) {
	  panic("not implemented") // TODO: Implement
  }

  func (blogsvc *BlogSvc) GetAllPosts() (*[]*dto.PostResponse, error) {
	  panic("not implemented") // TODO: Implement
  }

  func (blogsvc *BlogSvc) CreatePost(post *dto.CreatePostRequest) (*dto.PostResponse, error) {
	  panic("not implemented") // TODO: Implement
  }

  func (blogsvc *BlogSvc) AddComment(comment *dto.AddCommentRequest) (*dto.CommentResponse, error) {
	  panic("not implemented") // TODO: Implement
  }

  func (blogsvc *BlogSvc) UpdatePost(post *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	  panic("not implemented") // TODO: Implement
  }
}

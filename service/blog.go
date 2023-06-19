package service

import (
	"gitea.slauson.io/blog/blog-ms/dao"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
)

type BlogSvc struct {
	dao dao.PostDao
}

type BlogService interface {
	CreatePost(post *siogeneric.BlogPost) error
}

func NewBlogSvc() *BlogSvc {
	return &BlogSvc{
		dao: dao.NewPostDao(),
	}
}

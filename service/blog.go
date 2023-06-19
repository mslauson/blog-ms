package service

import "gitea.slauson.io/blog/blog-ms/dao"

type BlogSvc struct {
	dao dao.PostDao
}

type BlogService interface{}

func NewBlogSvc() *BlogSvc {
	return &BlogSvc{
		dao: dao.NewPostDao(),
	}
}

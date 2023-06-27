package service

import (
	"fmt"
	"time"

	"gitea.slauson.io/blog/post-ms/dto"
	"gitea.slauson.io/slausonio/go-types/sioblog"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
)

func buildCreatePostEntity(req *dto.CreatePostRequest) *sioblog.BlogPost {
	fmt.Println(req)
	return &sioblog.BlogPost{
		Title:       req.Title,
		Body:        req.Body,
		CreatedByID: req.CreatedByID,
		PostedDate:  time.Now(),
	}
}

func buildAddCommentEntity(req *dto.AddCommentRequest) *sioblog.BlogComment {
	return &sioblog.BlogComment{
		Content:     req.Content,
		CommentDate: time.Now(),
		PostID:      req.PostID,
		UserID:      req.UserID,
	}
}

func buildUpdatePostEntity(req *dto.UpdatePostRequest) *sioblog.BlogPost {
	title := handleStringForUpdates(req.Title)
	body := handleStringForUpdates(req.Body)

	return &sioblog.BlogPost{
		Title:       title,
		Body:        body,
		UpdatedByID: req.UpdatedByID,
		UpdatedDate: time.Now(),
	}
}

func buildUpdateCommentEntity(req *dto.UpdateCommentRequest) *sioblog.BlogComment {
	content := handleStringForUpdates(req.Content)

	return &sioblog.BlogComment{
		Content:     content,
		UpdatedDate: time.Now(),
	}
}

func buildPostResponse(entity *sioblog.BlogPost) *dto.PostResponse {
	return &dto.PostResponse{
		ID:           entity.ID.Int64,
		Title:        entity.Title.String,
		Body:         entity.Body.String,
		PostedDate:   entity.PostedDate.Time,
		UpdatedDate:  entity.UpdatedDate.Time,
		DeletionDate: entity.DeletionDate.Time,
		SoftDeleted:  entity.SoftDeleted.Bool,
		CreatedByID:  entity.CreatedByID.Int64,
		UpdatedByID:  entity.UpdatedByID.Int64,
		Comments:     buildCommentResponses(entity.Comments),
	}
}

func buildAllPostsResponse(entities *[]*sioblog.BlogPost) *[]*dto.PostResponse {
	var postResponses []*dto.PostResponse
	for _, entity := range *entities {
		postResponses = append(postResponses, buildPostResponse(entity))
	}
	return &postResponses
}

func buildCommentResponse(entity *sioblog.BlogComment) *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:           entity.ID.Int64,
		Content:      entity.Content.String,
		CommentDate:  entity.CommentDate.Time,
		UpdatedDate:  entity.UpdatedDate.Time,
		SoftDeleted:  entity.SoftDeleted.Bool,
		DeletionDate: entity.DeletionDate.Time,
		PostID:       entity.PostID.Int64,
		UserID:       entity.UserID.Int64,
	}
}

func buildCommentResponses(entities *[]*sioblog.BlogComment) *[]*dto.CommentResponse {
	if entities == nil {
		return nil
	}

	var commentReponses []*dto.CommentResponse
	for _, entity := range *entities {
		commentReponses = append(commentReponses, buildCommentResponse(entity))
	}
	return &commentReponses
}

func handleStringForUpdates(s string) siogeneric.SioNullString {
	if s == "" {
		return siogeneric.NewSioInvalidNullString()
	}
	return s)
}

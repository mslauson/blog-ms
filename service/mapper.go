package service

import (
	"database/sql"
	"fmt"
	"time"

	"gitea.slauson.io/blog/post-ms/dto"
	"gitea.slauson.io/slausonio/go-types/sioblog"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
)

func buildCreatePostEntity(req *dto.CreatePostRequest) *sioblog.BlogPost {
	fmt.Println(req)
	return &sioblog.BlogPost{
		Title:       handleStringForUpdates(req.Title),
		Body:        handleStringForUpdates(req.Body),
		CreatedByID: req.CreatedByID,
		PostedDate:  time.Now(),
	}
}

func buildAddCommentEntity(req *dto.AddCommentRequest) *sioblog.BlogComment {
	return &sioblog.BlogComment{
		Content:     handleStringForUpdates(req.Content),
		CommentDate: time.Now(),
		PostID:      req.PostID,
		UserID:      req.UserID,
	}
}

func buildUpdatePostEntity(req *dto.UpdatePostRequest) *sioblog.BlogPost {
	return &sioblog.BlogPost{
		Title:       handleStringForUpdates(req.Title),
		Body:        handleStringForUpdates(req.Body),
		UpdatedByID: req.UpdatedByID,
		UpdatedDate: time.Now(),
	}
}

func buildUpdateCommentEntity(req *dto.UpdateCommentRequest) *sioblog.BlogComment {
	return &sioblog.BlogComment{
		Content:     handleStringForUpdates(req.Content),
		UpdatedDate: time.Now(),
	}
}

func buildPostResponse(entity *sioblog.BlogPost) *dto.PostResponse {
	return &dto.PostResponse{
		ID:           entity.ID,
		Title:        entity.Title.String,
		Body:         entity.Body.String,
		PostedDate:   entity.PostedDate,
		UpdatedDate:  entity.UpdatedDate,
		DeletionDate: entity.DeletionDate,
		SoftDeleted:  entity.SoftDeleted,
		CreatedByID:  entity.CreatedByID,
		UpdatedByID:  entity.UpdatedByID,
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
		ID:           entity.ID,
		Content:      entity.Content.String,
		CommentDate:  entity.CommentDate,
		UpdatedDate:  entity.UpdatedDate,
		SoftDeleted:  entity.SoftDeleted,
		DeletionDate: entity.DeletionDate,
		PostID:       entity.PostID,
		UserID:       entity.UserID,
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

func handleStringForUpdates(s string) sql.NullString {
	if s == "" {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		Valid:  true,
		String: s,
	}
}

package service

import (
	"fmt"
	"time"

	"gitea.slauson.io/blog/post-ms/dto"
	"gitea.slauson.io/slausonio/go-types/sioblog"
	"gitea.slauson.io/slausonio/go-utils/siodao"
)

func buildCreatePostEntity(req *dto.CreatePostRequest) *sioblog.BlogPost {
	fmt.Println(req)
	return &sioblog.BlogPost{
		Title:       siodao.BuildNullString(req.Title),
		Body:        siodao.BuildNullString(req.Body),
		CreatedByID: req.CreatedByID,
		PostedDate:  time.Now(),
	}
}

func buildAddCommentEntity(req *dto.AddCommentRequest) *sioblog.BlogComment {
	return &sioblog.BlogComment{
		Content:     siodao.BuildNullString(req.Content),
		CommentDate: time.Now(),
		PostID:      req.PostID,
		UserID:      req.UserID,
	}
}

func buildUpdatePostEntity(req *dto.UpdatePostRequest) *sioblog.BlogPost {
	return &sioblog.BlogPost{
		Title:       siodao.BuildNullString(req.Title),
		Body:        siodao.BuildNullString(req.Body),
		UpdatedByID: siodao.BuildNullInt64(req.UpdatedByID),
		UpdatedDate: siodao.BuildNullTime(time.Now()),
	}
}

func buildUpdateCommentEntity(req *dto.UpdateCommentRequest) *sioblog.BlogComment {
	return &sioblog.BlogComment{
		Content:     siodao.BuildNullString(req.Content),
		UpdatedDate: siodao.BuildNullTime(time.Now()),
	}
}

func buildPostResponse(entity *sioblog.BlogPost) *dto.PostResponse {
	return &dto.PostResponse{
		ID:           entity.ID,
		Title:        entity.Title.String,
		Body:         entity.Body.String,
		PostedDate:   entity.PostedDate,
		UpdatedDate:  entity.UpdatedDate.Time,
		DeletionDate: entity.DeletionDate.Time,
		SoftDeleted:  entity.SoftDeleted,
		CreatedByID:  entity.CreatedByID,
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
		ID:           entity.ID,
		Content:      entity.Content.String,
		CommentDate:  entity.CommentDate,
		UpdatedDate:  entity.UpdatedDate.Time,
		SoftDeleted:  entity.SoftDeleted,
		DeletionDate: entity.DeletionDate.Time,
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

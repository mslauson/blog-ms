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
		Title:       siogeneric.NewSioNullString(req.Title),
		Body:        siogeneric.NewSioNullString(req.Body),
		CreatedByID: siogeneric.NewSioNullInt64(req.CreatedByID),
		PostedDate:  siogeneric.NewSioNullTime(time.Now()),
	}
}

func buildAddCommentEntity(req *dto.AddCommentRequest) *sioblog.BlogComment {
	return &sioblog.BlogComment{
		Content:     siogeneric.NewSioNullString(req.Content),
		CommentDate: siogeneric.NewSioNullTime(time.Now()),
		PostID:      siogeneric.NewSioNullInt64(req.PostID),
		UserID:      siogeneric.NewSioNullInt64(req.UserID),
	}
}

func buildUpdatePostEntity(req *dto.UpdatePostRequest) *sioblog.BlogPost {
	return &sioblog.BlogPost{
		Title:       siogeneric.NewSioNullString(req.Title),
		Body:        siogeneric.NewSioNullString(req.Body),
		UpdatedByID: siogeneric.NewSioNullInt64(req.UpdatedByID),
		UpdatedDate: siogeneric.NewSioNullTime(time.Now()),
	}
}

func buildUpdateCommentEntity(req *dto.UpdateCommentRequest) *sioblog.BlogComment {
	return &sioblog.BlogComment{
		Content:     siogeneric.NewSioNullString(req.Content),
		UpdatedDate: siogeneric.NewSioNullTime(time.Now()),
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

package service

import (
	"time"

	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
)

func buildCreatePostEntity(req *dto.CreatePostRequest) *siogeneric.BlogPost {
	return &siogeneric.BlogPost{
		Title:       siogeneric.NewSioNullString(req.Title),
		Body:        siogeneric.NewSioNullString(req.Body),
		CreatedByID: siogeneric.NewSioNullInt64(req.CreatedByID),
		PostedDate:  siogeneric.NewSioNullTime(time.Now()),
	}
}

func buildAddCommentEntity(req *dto.AddCommentRequest) *siogeneric.BlogComment {
	return &siogeneric.BlogComment{
		Content:     siogeneric.NewSioNullString(req.Content),
		CommentDate: siogeneric.NewSioNullTime(time.Now()),
		PostID:      siogeneric.NewSioNullInt64(req.PostID),
		UserID:      siogeneric.NewSioNullInt64(req.UserID),
	}
}

func buildUpdatePostEntity(req *dto.UpdatePostRequest) *siogeneric.BlogPost {
	return &siogeneric.BlogPost{
		Title:       siogeneric.NewSioNullString(req.Title),
		Body:        siogeneric.NewSioNullString(req.Body),
		UpdatedByID: siogeneric.NewSioNullInt64(req.UpdatedByID),
		UpdatedDate: siogeneric.NewSioNullTime(time.Now()),
	}
}

func buildUpdateCommentEntity(req *dto.UpdateCommentRequest) *siogeneric.BlogComment {
	return &siogeneric.BlogComment{
		Content:     siogeneric.NewSioNullString(req.Content),
		UpdatedDate: siogeneric.NewSioNullTime(time.Now()),
	}
}

func buildPostResponse(entity *siogeneric.BlogPost) *dto.PostResponse {
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

func buildAllPostsResponse(entities *[]*siogeneric.BlogPost) *[]*dto.PostResponse {
	var postResponses []*dto.PostResponse
	for _, entity := range *entities {
		postResponses = append(postResponses, buildPostResponse(entity))
	}
	return &postResponses
}

func buildCommentResponses(comments *[]*siogeneric.BlogComment) *[]*dto.CommentResponse {
	var commentReponses []*dto.CommentResponse
	for _, comment := range *comments {
		cr := &dto.CommentResponse{
			ID:           comment.ID.Int64,
			Content:      comment.Content.String,
			CommentDate:  comment.CommentDate.Time,
			UpdatedDate:  comment.UpdatedDate.Time,
			SoftDeleted:  comment.SoftDeleted.Bool,
			DeletionDate: comment.DeletionDate.Time,
			PostID:       comment.PostID.Int64,
			UserID:       comment.UserID.Int64,
		}
		commentReponses = append(commentReponses, cr)
	}
	return &commentReponses
}

package service

import (
	"time"

	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
)

func mapCreatePostRequestToEntity(req *dto.CreatePostRequest) *siogeneric.BlogPost {
	return &siogeneric.BlogPost{
		Title:       siogeneric.NewSioNullString(req.Title),
		Body:        siogeneric.NewSioNullString(req.Body),
		CreatedByID: siogeneric.NewSioNullInt64(req.CreatedByID),
		PostedDate:  siogeneric.NewSioNullTime(time.Now()),
	}
}

func mapAddCommentRequestToEntity(req *dto.AddCommentRequest) *siogeneric.BlogComment {
	return &siogeneric.BlogComment{
		Content:     siogeneric.NewSioNullString(req.Content),
		CommentDate: siogeneric.NewSioNullTime(time.Now()),
		PostID:      siogeneric.NewSioNullInt64(req.PostID),
		UserID:      siogeneric.NewSioNullInt64(req.UserID),
	}
}

func mapUpdatePostRequestToEntity(req *dto.UpdatePostRequest) *siogeneric.BlogPost {
	return &siogeneric.BlogPost{
		Title:       siogeneric.NewSioNullString(req.Title),
		Body:        siogeneric.NewSioNullString(req.Body),
		UpdatedByID: siogeneric.NewSioNullInt64(req.UpdatedByID),
		UpdatedDate: siogeneric.NewSioNullTime(time.Now()),
	}
}

func mapUpdateCommentRequestToEntity(req *dto.UpdateCommentRequest) *siogeneric.BlogComment {
	return &siogeneric.BlogComment{
		Content:     siogeneric.NewSioNullString(req.Content),
		UpdatedDate: siogeneric.NewSioNullTime(time.Now()),
	}
}

func mapPostResponse(entity *siogeneric.BlogPost) *dto.PostResponse {
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
		Comments:     mapCommentResponses(entity.Comments),
	}
}

func mapCommentResponses(comments *[]*siogeneric.BlogComment) *[]*dto.CommentResponse {
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

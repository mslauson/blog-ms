package dto

// type BlogPost struct {
// 	ID           SioNullInt64
// 	Title        SioNullString
// 	Body         SioNullString
// 	PostedDate   SioNullTime
// 	UpdatedDate  SioNullTime
// 	DeletionDate SioNullTime
// 	SoftDeleted  SioNullBool
// 	CreatedByID  SioNullInt64
// 	UpdatedByID  SioNullInt64
// 	Comments     *[]BlogComment
// }
//
// type BlogComment struct {
// 	ID           SioNullInt64
// 	Content      SioNullString
// 	CommentDate  SioNullTime
// 	UpdatedDate  SioNullTime
// 	SoftDeleted  SioNullBool
// 	DeletionDate SioNullTime
// 	PostID       SioNullInt64
// 	UserID       SioNullInt64
// }

type CreatePostRequest struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	CreatedByID int64  `json:"createdById"`
}

type UpdatePostRequest struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	UpdatedByID int64  `json:"updatedById"`
}

type AddCommentRequest struct {
	Content string `json:"content"`
	PostID  int64  `json:"postId"`
	UserID  int64  `json:"userId"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

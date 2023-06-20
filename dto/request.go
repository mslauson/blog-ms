package dto

type CreatePostRequest struct {
	Title       string `json:"title"       binding:"required"`
	Body        string `json:"body"        binding:"required"`
	CreatedByID int64  `json:"createdById" binding:"required"`
}

type UpdatePostRequest struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	UpdatedByID int64  `json:"updatedById" binding:"required"`
}

type AddCommentRequest struct {
	Content string `json:"content" binding:"required"`
	PostID  int64  `json:"postId"  binding:"required"`
	UserID  int64  `json:"userId"  binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

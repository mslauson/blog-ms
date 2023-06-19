package dto

import "time"

type PostResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	PostedDate   time.Time `json:"postedDate"`
	UpdatedDate  time.Time `json:"updatedDate"`
	DeletionDate time.Time `json:"deletionDate"`
	SoftDeleted  bool      `json:"softDeleted"`
	CreatedByID  int64     `json:"createdById"`
	UpdatedByID  int64     `json:"updatedById"`
}

type CommentResponse struct {
	ID           int64     `json:"id"`
	Content      string    `json:"content"`
	CommentDate  time.Time `json:"commentDate"`
	SoftDeleted  bool      `json:"softDeleted"`
	DeletionDate time.Time `json:"deletionDate"`
	PostID       int64     `json:"postId"`
	UserID       int64     `json:"userId"`
}

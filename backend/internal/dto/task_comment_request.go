package dto

type CreateTaskCommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}

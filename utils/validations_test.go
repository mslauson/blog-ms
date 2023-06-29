package utils

import (
	"testing"

	"gitea.slauson.io/blog/post-ms/constants"
	"gitea.slauson.io/blog/post-ms/dto"
	"gitea.slauson.io/blog/post-ms/integration/mockdata"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateRequest(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		req     *dto.CreatePostRequest
		err     string
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     mockdata.CreatePostRequest,
			wantErr: false,
		},

		{
			name: "Bad Title",
			req: &dto.CreatePostRequest{
				Title:       mockdata.LongTitle,
				Body:        "Test Body",
				CreatedByID: 1,
			},
			err:     constants.TITLE_TOO_LONG,
			wantErr: true,
		},
		{
			name: "Bad Body",
			req: &dto.CreatePostRequest{
				Title:       "Title",
				Body:        mockdata.LongBody,
				CreatedByID: 1,
			},
			err:     constants.BODY_TOO_LONG,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreatePostRequest(tt.req)
			if tt.wantErr {
				assert.NotNilf(t, err, "expected error but got nil")
				assert.Equalf(
					t,
					tt.err,
					err.Error(),
					"expected error %s but got %s",
					tt.err,
					err.Error(),
				)
			} else {
				assert.Nilf(t, err, "expected no error but got error")
			}
		})
	}
}

func TestValidateUpdateRequest(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		req     *dto.UpdatePostRequest
		err     string
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     mockdata.UpdatePostRequest,
			wantErr: false,
		},
		{
			name: "Valid only Title",
			req: &dto.UpdatePostRequest{
				Title:       "abc",
				UpdatedByID: 1,
			},
			wantErr: false,
		},
		{
			name: "Valid only Body",
			req: &dto.UpdatePostRequest{
				Body:        "abc",
				UpdatedByID: 1,
			},
			wantErr: false,
		},
		{
			name: "Bad Title",
			req: &dto.UpdatePostRequest{
				Title:       mockdata.LongTitle,
				Body:        "Test Body",
				UpdatedByID: 1,
			},
			err:     constants.TITLE_TOO_LONG,
			wantErr: true,
		},
		{
			name: "Bad Body",
			req: &dto.UpdatePostRequest{
				Title:       "Title",
				Body:        mockdata.LongBody,
				UpdatedByID: 1,
			},
			err:     constants.BODY_TOO_LONG,
			wantErr: true,
		},
		{
			name: "No Update",
			req: &dto.UpdatePostRequest{
				UpdatedByID: 1,
			},
			err:     constants.POST_UPDATE_INVALID,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUpdatePostRequest(tt.req)
			if tt.wantErr {
				assert.NotNilf(t, err, "expected error but got nil")
				assert.Equalf(
					t,
					tt.err,
					err.Error(),
					"expected error %s but got %s",
					tt.err,
					err.Error(),
				)
			} else {
				assert.Nilf(t, err, "expected no error but got error")
			}
		})
	}
}

func TestValidateAddCommentRequest(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		req     *dto.AddCommentRequest
		err     string
		wantErr bool
	}{
		{
			name:    "happy",
			req:     mockdata.AddCommentRequest,
			wantErr: false,
		},
		{
			name: "comment too long",
			req: &dto.AddCommentRequest{
				Content: mockdata.LongComment,
			},
			wantErr: true,
			err:     constants.COMMENT_TOO_LONG,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAddCommentRequest(tt.req)
			if tt.wantErr {
				assert.NotNilf(t, err, "expected error but got nil")
				assert.Equalf(
					t,
					tt.err,
					err.Error(),
					"expected error %s but got %s",
					tt.err,
					err.Error(),
				)
			} else {
				assert.Nilf(t, err, "expected no error but got error")
			}
		})
	}
}

func TestValidateUpdateCommentRequest(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		req     *dto.UpdateCommentRequest
		err     string
		wantErr bool
	}{
		{
			name:    "happy",
			req:     mockdata.UpdateCommentRequest,
			wantErr: false,
		},
		{
			name: "comment too long",
			req: &dto.UpdateCommentRequest{
				Content: mockdata.LongComment,
			},
			wantErr: true,
			err:     constants.COMMENT_TOO_LONG,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUpdateCommentRequest(tt.req)
			if tt.wantErr {
				assert.NotNilf(t, err, "expected error but got nil")
				assert.Equalf(
					t,
					tt.err,
					err.Error(),
					"expected error %s but got %s",
					tt.err,
					err.Error(),
				)
			} else {
				assert.Nilf(t, err, "expected no error but got error")
			}
		})
	}
}

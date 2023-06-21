package utils

import (
	"fmt"

	"gitea.slauson.io/blog/blog-ms/constants"
	"gitea.slauson.io/blog/blog-ms/dto"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
)

type BlogValidation struct {
	enc       *sioUtils.EncryptionUtil
	validator *sioUtils.SioValidator
}

func NewValidator() *BlogValidation {
	return &BlogValidation{
		enc:       sioUtils.NewEncryptionUtil(),
		validator: sioUtils.NewValidator(),
	}
}

func (bv *BlogValidation) ValidateCreatePostRequest(req *dto.CreatePostRequest) error {
	if err := bv.validateTitle(req.Title); err != nil {
		return err
	}

	if err := bv.validateBody(req.Body); err != nil {
		return err
	}

	return nil
}

func (bv *BlogValidation) ValidateUpdatePostRequest(req *dto.UpdatePostRequest) error {
	if req.Title == "" && req.Body == "" {
		return sioerror.NewSioBadRequestError(constants.POST_UPDATE_INVALID)
	}

	if req.Title != "" {
		if err := bv.validateTitle(req.Title); err != nil {
			return err
		}
	}

	if req.Body != "" {
		if err := bv.validateBody(req.Body); err != nil {
			return err
		}
	}

	return nil
}

func (bv *BlogValidation) ValidateAddCommentRequest(req *dto.AddCommentRequest) error {
	content, err := bv.enc.Decrypt(req.Content)
	if err != nil {
		return sioerror.NewSioInternalServerError(sioUtils.DecryptionFailed)
	}

	return bv.validateCommentContent(content)
}

func (bv *BlogValidation) ValidateUpdateCommentRequest(req *dto.UpdateCommentRequest) error {
	return bv.validateCommentContent(req.Content)
}

func (bv *BlogValidation) validateTitle(title string) error {
	if len(title) > 100 {
		return sioerror.NewSioBadRequestError(constants.TITLE_TOO_LONG)
	}

	return nil
}

func (bv *BlogValidation) validateBody(body string) error {
	fmt.Println("body:", len(body))
	if len(body) > 100000 {
		return sioerror.NewSioBadRequestError(constants.BODY_TOO_LONG)
	}

	return nil
}

func (bv *BlogValidation) validateCommentContent(content string) error {
	if len(content) > 3000 {
		return sioerror.NewSioBadRequestError(constants.COMMENT_TOO_LONG)
	}

	return nil
}

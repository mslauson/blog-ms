package utils

import (
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
	decryptedRequest := new(dto.CreatePostRequest)
	*decryptedRequest = *req
	if err := bv.enc.DecryptInterface(decryptedRequest); err != nil {
		return sioerror.NewSioInternalServerError(sioUtils.DecryptionFailed)
	}

	if err := bv.validateTitle(decryptedRequest.Title); err != nil {
		return err
	}

	if err := bv.validateBody(decryptedRequest.Body); err != nil {
		return err
	}

	return nil
}

func (bv *BlogValidation) ValidateUpdatePostRequest(req *dto.UpdatePostRequest) error {
	decryptedRequest := new(dto.UpdatePostRequest)
	*decryptedRequest = *req
	if err := bv.enc.DecryptInterface(decryptedRequest); err != nil {
		return sioerror.NewSioInternalServerError(sioUtils.DecryptionFailed)
	}

	if req.Title == "" && req.Body == "" {
		return sioerror.NewSioBadRequestError(POST_UPDATE_INVALID)
	}

	if req.Title != "" {
		if err := bv.validateTitle(decryptedRequest.Title); err != nil {
			return err
		}
	}

	if req.Body != "" {
		if err := bv.validateBody(decryptedRequest.Body); err != nil {
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
	content, err := bv.enc.Decrypt(req.Content)
	if err != nil {
		return sioerror.NewSioInternalServerError(sioUtils.DecryptionFailed)
	}

	return bv.validateCommentContent(content)
}

func (bv *BlogValidation) validateTitle(title string) error {
	if len(title) > 100 {
		return sioerror.NewSioBadRequestError(TITLE_TOO_LONG)
	}

	return nil
}

func (bv *BlogValidation) validateBody(body string) error {
	if len(body) > 100000 {
		return sioerror.NewSioBadRequestError(BODY_TOO_LONG)
	}

	return nil
}

func (bv *BlogValidation) validateCommentContent(content string) error {
	if len(content) > 100000 {
		return sioerror.NewSioBadRequestError(COMMENT_TOO_LONG)
	}

	return nil
}

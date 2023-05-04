package controller

import (
	"net/http"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioerror"

	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/gin-gonic/gin"

	"gitea.slauson.io/slausonio/iam-ms/service"
	"gitea.slauson.io/slausonio/iam-ms/utils"
)

type UserController struct {
	s service.IamUserService
}

//go:generate mockery --name IamUserController
type IamUserController interface {
	ListUsers(c *gin.Context)
	GetUserById(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdatePassword(c *gin.Context)
	UpdateEmail(c *gin.Context)
	UpdatePhone(c *gin.Context)
	DeleteUser(c *gin.Context)
}

func NewUserController() *UserController {
	return &UserController{
		s: service.NewUserService(),
	}
}

func (uc *UserController) ListUsers(c *gin.Context) {
	result, e := uc.s.ListUsers()

	if e != nil {
		_ = c.Error(e)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (uc *UserController) GetUserById(c *gin.Context) {
	id := c.Param("id")
	response, e := uc.s.GetUserByID(id)

	if e != nil {
		_ = c.Error(e)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (uc *UserController) CreateUser(c *gin.Context) {
	validations := utils.NewIamValidations()
	request := new(siogeneric.AwCreateUserRequest)
	err := sioUtils.DecryptAndHandle(request, c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = validations.ValidateCreateUserRequest(request)
	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}

	result, e := uc.s.CreateUser(request)

	if e != nil {
		_ = c.Error(e)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (uc *UserController) UpdatePassword(c *gin.Context) {
	validations := utils.NewIamValidations()
	id := c.Param("id")
	request := new(siogeneric.UpdatePasswordRequest)

	err := sioUtils.DecryptAndHandle(request, c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = validations.ValidateUpdatePasswordRequest(request)

	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}

	result, e := uc.s.UpdatePassword(id, request)

	if e != nil {
		_ = c.Error(e)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (uc *UserController) UpdateEmail(c *gin.Context) {
	validations := utils.NewIamValidations()
	id := c.Param("id")
	request := new(siogeneric.UpdateEmailRequest)
	err := sioUtils.DecryptAndHandle(request, c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = validations.ValidateUpdateEmailRequest(request)

	if err != nil {
		_ = c.Error(sioerror.NewSioBadRequestError(err.Error()))
		return
	}

	result, e := uc.s.UpdateEmail(id, request)
	if e != nil {
		_ = c.Error(e)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (uc *UserController) UpdatePhone(c *gin.Context) {
	validations := utils.NewIamValidations()
	id := c.Param("id")
	request := new(siogeneric.UpdatePhoneRequest)
	err := sioUtils.DecryptAndHandle(request, c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = validations.ValidateUpdatePhoneRequest(request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result, e := uc.s.UpdatePhone(id, request)
	if e != nil {
		_ = c.Error(e)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	response, err := uc.s.DeleteUser(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

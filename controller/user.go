package controller

import (
	"fmt"
	"net/http"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"gitea.slauson.io/slausonio/iam-ms/service"
	"gitea.slauson.io/slausonio/iam-ms/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	s   service.IamUserService
	enc *sioUtils.EncryptionUtil
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
	response := uc.s.ListUsers(c)
	if response != nil {
		c.JSON(http.StatusOK, response)
	}
}

func (uc *UserController) GetUserById(c *gin.Context) {
	id := c.Param("id")
	response := uc.s.GetUserByID(id, c)
	if response != nil {
		c.JSON(http.StatusOK, response)
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	validations := utils.NewIamValidations(c)
	request := new(sioModel.AwCreateUserRequest)
	err := c.BindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = uc.enc.DecryptInterface(request, false)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("decryption failed - unable to proceed"))
	}

	if !validations.ValidateCreateUserRequest(request) {
		return
	}

	result := uc.s.CreateUser(request, c)
	if result != nil {
		c.JSON(http.StatusOK, result)
		return
	}
}

func (uc *UserController) UpdatePassword(c *gin.Context) {
	validations := utils.NewIamValidations(c)
	id := c.Param("id")
	request := new(sioModel.UpdatePasswordRequest)

	err := c.BindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = uc.enc.DecryptInterface(request, false)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("decryption failed - unable to proceed"))
	}

	if !validations.ValidateUpdatePasswordRequest(request) {
		return
	}

	result := uc.s.UpdatePassword(id, request, c)
	if result != nil {
		c.JSON(http.StatusOK, result)
		return
	}
}

func (uc *UserController) UpdateEmail(c *gin.Context) {
	validations := utils.NewIamValidations(c)
	id := c.Param("id")
	request := new(sioModel.UpdateEmailRequest)

	err := c.BindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = uc.enc.DecryptInterface(request, false)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("decryption failed - unable to proceed"))
	}
	if !validations.ValidateUpdateEmailRequest(request) {
		return
	}
	result := uc.s.UpdateEmail(id, request, c)
	if result != nil {
		c.JSON(http.StatusOK, result)
		return
	}
}

func (uc *UserController) UpdatePhone(c *gin.Context) {
	validations := utils.NewIamValidations(c)
	id := c.Param("id")
	request := new(sioModel.UpdatePhoneRequest)

	err := c.BindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = uc.enc.DecryptInterface(request, false)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("decryption failed - unable to proceed"))
	}

	if !validations.ValidateUpdatePhoneRequest(request) {
		return
	}
	result := uc.s.UpdatePhone(id, request, c)
	if result != nil {
		c.JSON(http.StatusOK, result)
		return
	}
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	response := uc.s.DeleteUser(id, c)
	if response.Success {
		c.JSON(http.StatusOK, response)
	}
}

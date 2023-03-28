package controller

import (
	"net/http"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"gitea.slauson.io/slausonio/iam-ms/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	s service.IamUserService
}

type IamUserController interface {
	ListUsers(c *gin.Context)
	GetUserById(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdatePassword(c *gin.Context)
	UpdateEmail(c *gin.Context)
	UpdatePhone(c *gin.Context)
	DeleteUser(c *gin.Context)
}

func NewUserController(c *gin.Context) *UserController {
	return &UserController{
		s: service.NewUserService(c),
	}
}

func (uc *UserController) ListUsers(c *gin.Context) {
	response := uc.s.ListUsers()
	if response != nil {
		c.JSON(http.StatusOK, response)
	}
}

func (uc *UserController) GetUserById(c *gin.Context) {
	id := c.Param("id")
	newId, err := sioUtils.ConvertToInt64(id, c)
	if err != nil {
		return
	}
	response := uc.s.GetUserByID(id)
	if response != nil {
		c.JSON(http.StatusOK, response)
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	// TODO validations

	var request sioModel.AwCreateUserRequest
	err := c.BindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result := uc.s.CreateUser(&request)
	if result != nil {
		c.JSON(http.StatusOK, result)
		return
	}
}

func (uc *UserController) UpdatePassword(c *gin.Context) {
	// TODO validations
	id := c.Param("id")
	newId, err := sioUtils.ConvertToInt64(id, c)
	if err != nil {
		return
	}
	var request sioModel.UpdatePasswordRequest
	err = c.BindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result := uc.s.UpdatePassword(id, &request)
	if result != nil {
		c.JSON(http.StatusOK, result)
		return
	}
}

func (uc *UserController) UpdateEmail(c *gin.Context) {
	// TODO validations
	id := c.Param("id")
	newId, err := sioUtils.ConvertToInt64(id, c)
	if err != nil {
		return
	}
	var request sioModel.UpdateEmailRequest
	err = c.BindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result := uc.s.UpdateEmail(id, &request)
	if result != nil {
		c.JSON(http.StatusOK, result)
		return
	}
}

func (uc *UserController) UpdatePhone(c *gin.Context) {
	// TODO validations
	id := c.Param("id")
	newId, err := sioUtils.ConvertToInt64(id, c)
	if err != nil {
		return
	}
	var request sioModel.UpdatePhoneRequest
	err = c.BindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result := uc.s.UpdatePhone(id, &request)
	if result != nil {
		c.JSON(http.StatusOK, result)
		return
	}
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	newId, err := sioUtils.ConvertToInt64(id, c)
	if err != nil {
		return
	}
	response := uc.s.DeleteUser(id)
	if response.Success {
		c.JSON(http.StatusOK, response)
	}
}

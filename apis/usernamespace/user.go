package usernamespace

import (
	"net/http"
	"tp_go_gin_complex/models"
	"tp_go_gin_complex/models/user"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserSource struct {
}

// @Summary 获取用户列表信息
// @Tags user
// @Produce  json
// @Param   uid     path    int     true        "User ID"
// @Success 200 {object} user.User
// @Router /v1/api/user/{uid} [get]
func (s *UserSource) Get(c *gin.Context) {
	u := &user.User{}
	err := c.BindUri(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}
	err = u.Sync(models.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, &ResultResponse{Message: err.Error()})
		return
	}
	c.JSON(200, u)
}

func (s *UserSource) Put(c *gin.Context) {
	u := &user.User{}
	err := c.BindUri(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}
	uinput := &UserUpdateQuery{}
	err = c.ShouldBindBodyWith(uinput, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}
	err = u.Sync(models.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, &ResultResponse{Message: err.Error()})
		return
	}
	u.Name = uinput.Name
	err = u.Save(models.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ResultResponse{Message: err.Error()})
		return
	}
	c.JSON(200, &ResultResponse{Succeed: true})

}
func (s *UserSource) Delete(c *gin.Context) {
	u := &user.User{}
	err := c.BindUri(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}
	err = u.Sync(models.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, &ResultResponse{Message: err.Error()})
		return
	}
	err = u.Delete(models.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ResultResponse{Message: err.Error()})
		return
	}
	c.JSON(200, &ResultResponse{Succeed: true})
}

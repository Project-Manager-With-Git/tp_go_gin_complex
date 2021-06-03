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
// @Produce  application/json
// @Param   uid  path    integer   true  "User ID" minimum(1)
// @Success 200 {object} user.User
// @Failure 400 {string} ResultResponse "请求数据不符合要求"
// @Failure 404 {string} ResultResponse "未找到指定资源"
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

// @Summary 更新指定用户信息
// @Tags user
// @Produce  application/json
// @Param   uid  path    integer   true  "User ID" minimum(1)
// @Success 200 {object} user.User
// @Failure 400 {string} ResultResponse "请求数据不符合要求"
// @Failure 404 {string} ResultResponse "未找到指定资源"
// @Failure 500 {string} ResultResponse "服务器处理失败"
// @Router /v1/api/user/{uid} [put]
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

// @Summary 删除指定用户
// @Tags user
// @Produce  application/json
// @Param   uid  path    integer   true  "User ID" minimum(1)
// @Success 200 {object} user.User
// @Failure 400 {string} ResultResponse "请求数据不符合要求"
// @Failure 404 {string} ResultResponse "未找到指定资源"
// @Failure 500 {string} ResultResponse "服务器处理失败"
// @Router /v1/api/user/{uid} [delete]
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

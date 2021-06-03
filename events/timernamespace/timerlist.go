package timernamespace
"github.com/gin-contrib/sse"


type TimerListSource struct {
}



// @Summary 创建新用户
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param name body UserCreateQuery true "用户名信息"
// @Success 200 {object} user.User "用户信息"
// @Failure 400 {string} ResultResponse "请求数据不符合要求"
// @Failure 500 {string} ResultResponse "服务器处理失败"
// @Router /user [post]
func (s *TimerListSource) Post(c *gin.Context) {
	// 请求参数校验
	uinput := &UserCreateQuery{}
	err := c.ShouldBindBodyWith(uinput, binding.JSON)
	if err != nil {
		c.PureJSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}
	//创建用户
	u := &user.User{
		Name: uinput.Name,
	}
	err = u.Save(models.DB)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, &ResultResponse{Message: err.Error()})
		return
	}
	c.PureJSON(200, u)
}

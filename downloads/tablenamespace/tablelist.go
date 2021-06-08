package tablenamespace

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type TableListSource struct {
}

// @Summary 获取用户列表信息
// @Tags downloads table
// @Produce ext/csv
// @Success 200 {object} UserListResponse "表格响应信息,会提供表格下载"
// @Failure 404 {string} ResultResponse "请求不符合要求处理失败"
// @Router /download/table [get]
func (s *TableListSource) Get(c *gin.Context) {
	query := TableSearchQuery{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.PureJSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}
	table_name := ""
	if query.From == "" {
		if query.To == "" {
			table_name = fmt.Sprintf("table_%s.csv", time.Now().Format("2006-01-02"))
		} else {
			table_name = fmt.Sprintf("table_%s.csv", query.To)
		}
	} else {
		if query.To == "" {
			table_name = fmt.Sprintf("table_%s.csv", query.From)
		} else {
			if query.From == query.To {
				table_name = fmt.Sprintf("table_%s.csv", query.From)
			} else {
				table_name = fmt.Sprintf("attachment; filename=table_%s_%s.csv", query.From, query.To)
			}
		}
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", table_name))
	c.Writer.Flush()
	clientGone := c.Writer.CloseNotify()
	ch := make(chan string)
	go func() {
		ch <- "foo,bar\r\n"
		ch <- "a,1\r\n"
		ch <- "b,2\r\n"
		ch <- "c,3\r\n"
		ch <- "d,4\r\n"
		close(ch)
	}()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			{
				return false
			}
		case message, isopen := <-ch:
			{
				if isopen {
					c.Render(-1, render.String{Format: message})
					return true
				} else {
					return false
				}
			}
		}
	})
}

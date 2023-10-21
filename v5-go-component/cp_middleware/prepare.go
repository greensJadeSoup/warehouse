package cp_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
)

var sqlList = []string{
	"select ",
	"update ",
	"delete ",
	"drop ",
	"insert ",
	"where ",
	"show ",
	"count(",
	" or ",
	"-- ",
	"exec ",
	"#",
	"alert ",
	"modify ",
	"rename ",
	"union ",
	"char ",
	"declare ",
	"/*",
	"*/",
	"^",
}

func SQLCheck(str string) error {
	strNew := strings.ToLower(str)
	for _, v := range sqlList {
		if strings.Contains(strNew, v) {
			cp_log.Error("invalid_char:" + v)
			return errors.New("invalid_char:" + v)
		}
	}

	return nil
}

func Prepare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var err error

		if c.Request.Method == http.MethodGet {
			//sql防注入
			//if err = SQLCheck(string(c.Request.RequestURI)); err != nil {
			//	c.AbortWithStatusJSON(200, cp_obj.NewResponse().Err(err.Error()))
			//	return
			//}
		} else {
			if strings.Contains(c.Request.Header.Get("Content-Type"), "application/json") {
				body, err = ioutil.ReadAll(c.Request.Body)
				if err != nil {
					c.AbortWithStatusJSON(200, cp_obj.NewResponse().Err(err.Error()))
					return
				}

				//sql防注入
				if err = SQLCheck(string(body)); err != nil {
					c.AbortWithStatusJSON(200, cp_obj.NewResponse().Err(err.Error()))
					return
				}

				c.Set(gin.BodyBytesKey, body) //根据gin框架源码，进行预处理，这样后面gin框架的ShouldBindWith就不需要再进行一次readbody
			} else if c.Request.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
				cp_log.Info("specific_request_1")
			} else {
				cp_log.Info("specific_request_2"+ c.Request.Header.Get("Content-Type"))
			}
		}

		rid := cp_util.NewGuid()
		c.Set(cp_constant.REQUEST_ID, rid)
		cp_log.Info(fmt.Sprintf(`[================New Req==================] [method=%s] [uri=%s] [request_id=%s] [body=%s]`,
			c.Request.Method, c.Request.RequestURI, rid, string(body)))

		c.Next()
	}
}



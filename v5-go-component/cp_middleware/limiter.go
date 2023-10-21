package cp_middleware

import (
	"github.com/gin-gonic/gin"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_obj"
	"sync/atomic"
)

type Limiter struct {
	ServeMaxChan chan struct{} //最大并发请求
	Count        struct {
			     Serving    		int32  		//请求中
			     Served      		uint64 		//累计请求
			     Panic      		uint64 		//累计异常
			     Reject     		uint64 		//累计拒绝
			     ServedMax 			uint32 		//最高峰值并发
		     }
}

func NewLimiter(ServingLimit int) *Limiter {
	l := &Limiter{
		ServeMaxChan: make(chan struct{}, ServingLimit),
	}

	return l
}

// ParamNameTOLower 签名验证
func InLimiter(l *Limiter) gin.HandlerFunc {

	return func(c *gin.Context) {
		if c.Request.RequestURI != "/heartbeat" {
			select {
				case l.ServeMaxChan <- struct{}{}:
					l.serveAcquire()
					defer l.serveRelease()
				default:
					//最大请求数已满，拒绝请求
					l.serveReject()
					d := &cp_obj.Response {
						Code: cp_constant.RESPONSE_CODE_COMMON_ERROR,
						Message: "请求人数已满，请稍候再试",
					}

					c.AbortWithStatusJSON(200, d)
					return
			}
		}

		c.Next()

	}

}

//serveAcquire 获取请求权
func (l *Limiter) serveAcquire() {
	atomic.AddUint64(&l.Count.Served, 1)

	serving := uint32(atomic.AddInt32(&l.Count.Serving, 1))
	if serving > l.Count.ServedMax {
		l.Count.ServedMax = serving //记录最高峰
	}
}

//serveReject 拒绝请求的汇总
func (l *Limiter) serveReject() {
	atomic.AddUint64(&l.Count.Reject, 1)
}

//servePanic 请求异常的汇总
func (l *Limiter) servePanic() {
	atomic.AddUint64(&l.Count.Panic, 1)
}

//serveRelease 释放请求权
func (l *Limiter) serveRelease() {
	<-l.ServeMaxChan
	atomic.AddInt32(&l.Count.Serving, -1)
}

package cp_tracing

import (
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_dc"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_mq"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

var TracerProducer *Tracer

//Tracer 请求链对象
type Tracer struct {
	OnOff			bool
	ActionProducer    	cp_mq.IProducer
	RuntimeProducer    	cp_mq.IProducer
}

//NewTracer 创建新的请求链对象
func NewTracer(dcConf *cp_dc.DcConfig) *Tracer {
	if dcConf.TraceLog.OnOff == false {
		return &Tracer{OnOff: false}
	}

	kafkaConf, err := dcConf.GetMQ("kafka")
	if err != nil {
		panic(err)
	}

	ap, err := cp_mq.NewProducer(
		"kafka",
		fmt.Sprintf(`{"address":["%s"],"topics":"%s"}`,
			strings.Join(kafkaConf.Server, `","`),
			dcConf.TraceLog.ActionTopic),
	)
	if err != nil {
		panic(err)
	}

	rp, err := cp_mq.NewProducer(
		"kafka",
		fmt.Sprintf(`{"address":["%s"],"topics":"%s"}`,
			strings.Join(kafkaConf.Server, `","`),
			dcConf.TraceLog.RuntimeTopic),
	)

	if err != nil {
		panic(err)
	}

	TracerProducer = &Tracer{
		OnOff: true,
		ActionProducer: ap,
		RuntimeProducer: rp,
	}
	return TracerProducer
}

func NewTraceApi(ctx *gin.Context, body []byte, response *cp_obj.Response, startTime time.Time, level cp_constant.TracingLevel, errMsg string) *cp_obj.TraceAction {
	data := &cp_obj.TraceAction {
		LogLevel: level,
		RequestID: ctx.GetString(cp_constant.REQUEST_ID),
		IP: ctx.ClientIP(),
		ErrMsg: errMsg,
		StartTime: startTime.Format(time.StampMicro),
		TotalTime: (time.Now().UnixNano() - startTime.UnixNano()) / 1000,
		Uri: ctx.Request.RequestURI,
		Body: string(body),
		Response: response,
		SessionKey: ctx.GetString(cp_constant.HTTP_HEADER_SESSION_KEY),
		AppID: ctx.GetString(cp_constant.HTTP_HEADER_APPID),
		ChainID: ctx.GetString(cp_constant.HTTP_HEADER_CHAIN_ID),
		ChainLevel: ctx.GetString(cp_constant.HTTP_HEADER_CHAIN_LEVEL),
	}

	return data
}

func NewTraceRuntime(svrName string, level cp_constant.TracingLevel, errMsg string) *cp_obj.TraceRuntime {
	data := &cp_obj.TraceRuntime {
		SvrName: svrName,
		LogLevel: level,
		StartTime: time.Now().Format(time.StampMicro),
		RequestID: cp_util.NewGuid(),
		ErrMsg: errMsg,
	}

	return data
}

func (this *Tracer) PushAction(data *cp_obj.TraceAction) error {
	if this.OnOff == false {
		return nil
	}

	pb, err := cp_obj.Cjson.Marshal(data)
	if err != nil {
		return err
	}

	//err = this.ActionProducer.Publish(pb, "")
	//if err != nil {
	//	return err
	//}

	cp_log.Info("tracelog-action send successd: " + string(pb))

	return nil
}

func (this *Tracer) PushRuntime(data *cp_obj.TraceRuntime) error {
	if this.OnOff == false {
		return nil
	}

	pb, err := cp_obj.Cjson.Marshal(data)
	if err != nil {
		return err
	}

	err = this.RuntimeProducer.Publish(pb, "")
	if err != nil {
		return err
	}

	cp_log.Info("tracelog-runtime send successd: " + string(pb))

	return nil
}




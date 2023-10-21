package cp_app

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
)

//APIClient api客户端对象
type CallClient struct {
	svr		*Ins
	AccessToken	string
	ExpireTime	time.Time
}

func NewCallClient(svr *Ins) (*CallClient) {
	return &CallClient {
		svr:	svr,
	}
}

//微服务之间互相调用接口
func (this *CallClient) NewCall(solder *BaseController, callName, bodyOrQuery string) ([]byte, error) {
	var retryStr, callUrl, address string
	var req *http.Request
	var resp *http.Response
	var err error
	var respBody []byte

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3) //设置建立连接超时
				if err != nil {
					return nil, err
				}

				//c.SetDeadline(time.Now().Add(5 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	port := cp_api.SvrApiList[callName].Port
	uri := cp_api.SvrApiList[callName].UriQuery
	method := cp_api.SvrApiList[callName].Method

	if this.svr.DataCenter.Base.IsLocal {
		address = Instance.DataCenter.Base.IP
		callUrl = "http://" + address + uri
	} else {
		address = Instance.DataCenter.Base.IP + ":" + strconv.Itoa(port)
		callUrl = "http://" + address + uri
	}

	if method == http.MethodGet {
		req, _ = http.NewRequest(method, callUrl + "?" + bodyOrQuery, nil)
	} else {
		req, _ = http.NewRequest(method, callUrl, bytes.NewBufferString(bodyOrQuery))
		req.Header.Add("Content-Type", "application/json")
		cp_log.Info(fmt.Sprintf("[NewCall Req]:%s [Body]:%s", callUrl, bodyOrQuery))
	}

	if this.svr.DataCenter.Base.IsLocal { //本地调试和服务间调试不需要签名校验
		req.Header.Add(cp_constant.HTTP_HEADER_APPID, cp_constant.APPID_LOCAL)
	} else {
		req.Header.Add(cp_constant.HTTP_HEADER_APPID, cp_constant.APPID_SERVER)
	}

	req.Header.Add(cp_constant.HTTP_HEADER_CHAIN_LEVEL, solder.Ctx.GetString(cp_constant.HTTP_HEADER_CHAIN_LEVEL))
	req.Header.Add(cp_constant.HTTP_HEADER_CHAIN_ID, solder.Ctx.GetString(cp_constant.HTTP_HEADER_CHAIN_ID))
	req.Header.Add(cp_constant.HTTP_HEADER_SESSION_KEY, solder.Ctx.GetString(cp_constant.HTTP_HEADER_SESSION_KEY))

	if solder.Si != nil {
		si, err := cp_obj.Cjson.Marshal(solder.Si)
		if err != nil {
			cp_log.Error(err.Error())
		}
		req.Header.Add(cp_constant.HTTP_HEADER_SESSION_INFO, string(si))
	}

	for i, l := 1, 3; i <= l; i++ {
		resp, err = httpClient.Do(req)
		if err != nil {
			retryStr += fmt.Sprintf("第[%d]次请求[%s]接口失败: %s\r\n", i, callUrl, err.Error())
			time.Sleep(10 * time.Millisecond)
			continue
		} else {
			defer resp.Body.Close()
			respBody, err = ioutil.ReadAll(resp.Body)
			break
		}
	}

	if retryStr != "" {
		retryStr = this.svr.ServerName + "调用微服务" + callUrl + "请求失败: \r\n" + retryStr
		cp_log.Error(retryStr)
	}

	if err != nil {
		return nil, cp_error.NewSysError(retryStr)
	}

	cp_log.Info(fmt.Sprintf("[NewCall Resp]:%s [Body]:%s", callUrl, string(respBody)))
	return respBody, nil
}

package cp_microsrv

import (
	"fmt"
	"warehouse/v5-go-component/cp_error"
	"strconv"
	"strings"
)

//APIService API服务对象
type MicroService struct {
	MicroSvrMgr	*MicroSvrManager

	SID		string
	Type     	string
	Name     	string
	Version  	string

	IP       	string
	HTTPPort 	int
	HttpAddress 	string

	IsLocal     	bool
	IsInit      	bool

	StatusChange	bool
}

//Init 初始化Service的基本参数
func (mc *MicroService) Init(sid string) (error) {
	if mc.IsInit {
		return nil
	}

	mc.SID = sid
	msSlice := strings.Split(sid, "-")
	idsLen := len(msSlice)
	if idsLen == 0 {
		return cp_error.NewSysError("微服务sid错误")
	}

	if idsLen != 4 {  //长度错误
		return cp_error.NewSysError("微服务sid长度错误")
	} else {
		mc.Type = msSlice[0]
		mc.Name = msSlice[1]
		mc.IP = msSlice[2]
		mc.HTTPPort, _ = strconv.Atoi(msSlice[3])
	}

	if mc.IP == "127.0.0.1" {
		mc.IsLocal = true
	}

	if mc.MicroSvrMgr.isLocal == true {
		//mc.HttpAddress = fmt.Sprintf("https://%s", mc.IP) //证书过期了，暂时屏蔽
		mc.HttpAddress = fmt.Sprintf("http://%s", mc.IP)
	} else {
		mc.HttpAddress = fmt.Sprintf("http://%s", mc.IP)
	}
	mc.IsInit = true
	return nil
}


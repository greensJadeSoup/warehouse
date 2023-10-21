package model

import (
	"time"
	"warehouse/v5-go-component/cp_obj"
)

type SSOLoginMD struct {
	ID			uint64    	`json:"id" xorm:"id pk autoincr"`
	UserID			uint64   	`json:"user_id" xorm:"user_id"`
	Account			string		`json:"account" xorm:"account"`
	AccountType		string		`json:"account_type" xorm:"account_type"`
	SessionKey     		string    	`json:"session_key" xorm:"session_key"`
	DeviceType     		string    	`json:"device_type" xorm:"device_type"`		// web,wap,android,ios,pad,tv
	DeviceInfo  		string    	`json:"device_info" xorm:"device_info"` 	// 手机设备厂商，比如华为oppo
	LoginType      		string    	`json:"login_type" xorm:"login_type"`		// account:账号登录 email:邮箱 phone:手机 third:第三方 face:扫脸登陆
	Kick           		uint8     	`json:"kick" xorm:"kick"`			// 是否被踢
	DeviceID  		string    	`json:"device_id" xorm:"device_id"` 		// 设备ID，非必传
	AppVersion		string    	`json:"app_version" xorm:"app_version"`		// 软件版本
	LoginTime      		cp_obj.Datetime `json:"login_time" xorm:"login_time"`		// 最新登陆时间
	LogoutTime     		cp_obj.Datetime `json:"logout_time" xorm:"logout_time"`		// 登出时间
	LastActiveDate      	cp_obj.Datetime `json:"last_active_time" xorm:"last_active_time"` // 最新活跃日期
	LastActiveIP      	string 		`json:"last_active_ip" xorm:"last_active_ip"`	// 最后活跃日期
	ExpireTime     		cp_obj.Datetime `json:"expire_time" xorm:"expire_time"`		// 过期时间
	CreateTime		time.Time	`json:"create_time" xorm:"create_time created"`
	UpdateTime		time.Time	`json:"update_time" xorm:"update_time updated"`
}

func NewSsoLogin() *SSOLoginMD  {
	return &SSOLoginMD{}
}

// TableName 表名
func (m *SSOLoginMD) TableName() string {
	return "t_sso"
}

// DBConnectionName 数据库连接名
func (m *SSOLoginMD) DatabaseAlias() string {
	return "db_base"
}

package dav

import (
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_orm"
	"fmt"
	"strings"
)

//基本数据层
type SessionDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *SessionDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, &model.SSOLoginMD{})
}

func (this *SessionDAV) DBListOnlineModule(id uint64) ([]model.SSOLoginMD, error) {
	list := make([]model.SSOLoginMD, 0)

	searchSQL := fmt.Sprintf(`SELECT id, accountID, sessionKey, loginType, deviceType, deviceOS, deviceInfo, 
 		deviceID, appVersion, kick, loginTime, logoutTime, automaticLogin, 
		lastActiveDate, lastActiveIP, expireTime, created, modified FROM t_sso_login 
		WHERE accountID = %d and kick = 0 order by expireTime ASC`, id)

	cp_log.Debug(searchSQL)
	err := this.SQL(searchSQL).Find(&list)
	if err != nil {
		return nil, cp_error.NewSysError("[DBListOnlineModule]:" + err.Error())
	}

	return list, nil
}

func (this *SessionDAV) DBKick(kickList []string) error {
	md := &model.SSOLoginMD{}
	sql := fmt.Sprintf("Update %s Set kick=1, logout_time='%s' Where session_key in ('%s')", md.TableName(), cp_obj.NewDatetime().String(), strings.Join(kickList, "','"))

	cp_log.Debug(sql)
	_, err := this.Exec(sql)
	if err != nil {
		return cp_error.NewSysError("[DBKick]:" + err.Error())
	}
	return nil
}

func (this *SessionDAV) DBLogout(in *cbd.SessionReqCBD) (int64, error) {
	md := &model.SSOLoginMD{}
	sql := fmt.Sprintf("Update %s Set logout_time='%s',last_active_time='%s',last_active_ip='%s' Where session_key='%s'",
		md.TableName(), cp_obj.NewDatetime().String(), cp_obj.NewDatetime().String(), in.IP, in.SessionKey)

	cp_log.Debug(sql)
	res, err := this.Exec(sql)
	if err != nil {
		return 0, cp_error.NewSysError("[DBLogout]:" + err.Error())
	}

	return res.RowsAffected()
}

func (this *SessionDAV) DBInsertSession(s *cp_api.CheckSessionInfo) (int64, error) {
	md := model.NewSsoLogin()

	md.UserID = s.UserID
	md.Account = s.Account
	md.SessionKey = s.SessionKey
	md.AccountType = s.AccountType
	md.DeviceType = s.DeviceType
	md.DeviceInfo = s.DeviceInfo
	md.LoginType = s.LoginType
	md.DeviceID = ""
	md.AppVersion = "v1.0.0"
	md.Kick = 0
	md.LoginTime = s.LoginTime
	md.LastActiveDate = s.LastActiveDate
	md.LastActiveIP = s.LastActiveIP
	md.ExpireTime = s.ExpireTime

	execRow, err := this.Insert(md)
	if err != nil {
		return 0, cp_error.NewSysError("[DBInsertSession]Insert Session:" + err.Error())
	}

	return execRow, nil
}

func (this *SessionDAV) DBGetBySessionKey(sk string) (*model.SSOLoginMD, error) {
	md := model.NewSsoLogin()

	searchSQL := fmt.Sprintf(`select * from t_sso where session_key='%s'`, sk)

	cp_log.Debug(searchSQL)
	has, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[DBGetBySessionKey]:" + err.Error())
	} else if !has {
		return nil, nil
	}

	return md, nil
}

func (this *SessionDAV) DBUpdateLastActiveDate(in *cbd.SessionReqCBD) (int64, error) {
	md := &model.SSOLoginMD{}
	sql := fmt.Sprintf("Update %s Set last_active_time='%s',last_active_ip='%s' Where session_key='%s'",
		md.TableName(), cp_obj.NewDatetime().String(), in.IP, in.SessionKey)

	res, err := this.Exec(sql)
	if err != nil {
		return 0, cp_error.NewSysError("[DBUpdateLastActiveDate]:" + err.Error())
	}

	return res.RowsAffected()
}

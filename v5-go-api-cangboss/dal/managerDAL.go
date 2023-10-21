package dal

import (
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

//数据逻辑层

type ManagerDAL struct {
	dav.ManagerDAV
	Si *cp_api.CheckSessionInfo
}

func NewManagerDAL(si *cp_api.CheckSessionInfo) *ManagerDAL {
	return &ManagerDAL{Si: si}
}

func (this *ManagerDAL) AddManager(in *cbd.AddManagerReqCBD, warehouseRole string) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	salt := cp_util.ByteToString(cp_util.NewRandomBytes(6))

	md := &model.ManagerMD {
		Account: in.Account,
		Type: in.Type,
		VendorID: in.VendorID,
		RealName: in.RealName,
		Password: cp_util.Md5Encrypt(in.Password + salt),
		Salt: salt,
		Phone: in.Phone,
		Email: in.Email,
		WarehouseID: in.WarehouseID,
		AllowLogin: in.AllowLogin,
		WarehouseRole: warehouseRole,
		Note: in.Note,
	}

	return this.DBInsertAccount(md)
}

func (this *ManagerDAL) EditManager(in *cbd.EditManagerReqCBD, warehouseRole string) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.ManagerMD {
		ID: in.ManagerID,
		WarehouseID: in.WarehouseID,
		RealName: in.RealName,
		Phone: in.Phone,
		Email: in.Email,
		AllowLogin: in.AllowLogin,
		WarehouseRole: warehouseRole,
		Note: in.Note,
	}

	return this.DBUpdateManager(md)
}

func (this *ManagerDAL) GetModelByAccount(account string) (*model.ManagerMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByAccount(account)
}

func (this *ManagerDAL) ModifyPassword(in *cbd.ModifyPasswordReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	in.NewPassword = cp_util.Md5Encrypt(in.NewPassword + in.Salt)

	return this.DBModifyPassword(in)
}

func (this *ManagerDAL) ListManager(in *cbd.ListManagerReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListManager(in)
}

func (this *ManagerDAL) ResetPassword(in *cbd.ModifyPasswordReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	in.NewPassword = cp_util.Md5Encrypt(in.NewPassword + in.Salt)

	return this.DBModifyPassword(in)
}


func (this *ManagerDAL) GetModelByID(accountID uint64) (*model.ManagerMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(accountID)
}

func (this *ManagerDAL) DelManager(in *cbd.DelManagerReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBDelManager(in)
}

func (this *ManagerDAL) EditProfile(in *cbd.EditProfileReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := model.NewManager()
	md.ID = this.Si.UserID
	md.RealName = in.RealName
	md.Phone = in.Phone
	md.Email = in.Email
	md.CompanyName = in.CompanyName
	md.WechatNum = in.WechatNum

	_, err = this.UpdateProfile(md)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	err = this.Cache.Delete(cp_constant.REDIS_KEY_SESSIONKEY + this.Si.SessionKey)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	return nil
}

func (this *ManagerDAL) Test() error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	//cp_log.Info(this.Cache.LPUSH("lsj", "aaa"))
	//cp_log.Info(this.Cache.LPUSH("lsj", "bbb"))
	//cp_log.Info(this.Cache.RPUSH("lsj", "ccc"))
	//cp_log.Info(this.Cache.LLEN("lsj"))
	//cp_log.Info(this.Cache.LPOP("lsj"))
	//cp_log.Info(this.Cache.RPOP("lsj"))
	//cp_log.Info(this.Cache.LRANGE("lsj", 0, 100))
	//
	////cp_log.Info(this.Cache.SADD("ls", "cccc"))
	//cp_log.Info(this.Cache.SREM("ls", "aaa"))
	//cp_log.Info(this.Cache.SMEMBERS("ls"))
	//cp_log.Info(this.Cache.SCARD("ls"))
	//cp_log.Info(this.Cache.SPOP("ls"))

	return nil
}
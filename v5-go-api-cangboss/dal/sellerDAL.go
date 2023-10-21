package dal

import (
	"fmt"
	"strconv"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

//数据逻辑层

type SellerDAL struct {
	dav.SellerDAV
	Si *cp_api.CheckSessionInfo
}

func NewSellerDAL(si *cp_api.CheckSessionInfo) *SellerDAL {
	return &SellerDAL{Si: si}
}

func (this *SellerDAL) AddSeller(in *cbd.AddSellerReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	salt := cp_util.ByteToString(cp_util.NewRandomBytes(6))

	mdD, err := NewDiscountDAL(this.Si).GetDefaultByVendorID(in.VendorID)
	if err != nil {
		return cp_error.NewSysError(err)
	} else if mdD == nil {
		return cp_error.NewSysError("默认计价组不存在")
	}

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}

	md := &model.SellerMD {
		Account: in.Account,
		RealName: in.RealName,
		Password: cp_util.Md5Encrypt(in.Password + salt),
		Salt: salt,
		Phone: in.Phone,
		Email: in.Email,
		AllowLogin: in.AllowLogin,
		Note: in.Note,
	}

	err = this.DBInsertAccount(md)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	mdDs := &model.DiscountSellerMD {
		VendorID: in.VendorID,
		SellerID: md.ID,
		DiscountID: mdD.ID,
	}

	_, err = this.Insert(mdDs)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	err = NewVendorSellerDAL(this.Si).AddVendorSeller(&cbd.AddVendorSellerReqCBD{
		VendorID: in.VendorID,
		SellerID: md.ID,
	})
	if err != nil {
		this.Rollback()
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *SellerDAL) EditSeller(in *cbd.EditSellerReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	md := &model.SellerMD {
		ID: in.SellerID,
		RealName: in.RealName,
		Phone: in.Phone,
		Email: in.Email,
		AllowLogin: in.AllowLogin,
		Note: in.Note,
	}

	return this.DBUpdateSeller(md)
}

func (this *SellerDAL) ListSeller(in *cbd.ListSellerReqCBD) (*cp_orm.ModelList, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBListSeller(in)
}

func (this *SellerDAL) GetModelByAccount(account string) (*model.SellerMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByAccount(account)
}

func (this *SellerDAL) EditBalance(in *cbd.EditBalanceReqCBD) (err error) {
	err = this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	mdSeller, err := NewSellerDAL(this.Si).GetModelByID(in.SellerID)
	if err != nil {
		return err
	} else if mdSeller == nil {
		return cp_error.NewNormalError("账号不存在:" + strconv.FormatUint(in.SellerID, 10))
	}

	mdVs, err := NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, in.SellerID)
	if err != nil {
		return err
	} else if mdVs == nil {
		return cp_error.NewNormalError("无该账号访问权:" + strconv.FormatUint(in.SellerID, 10))
	}

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.DeferHandle(&err)

	if in.Type == constant.CALCULATE_TYPE_SUB {
		if mdVs.Balance < in.Num {
			return cp_error.NewNormalError("余额不足:" + strconv.FormatFloat(mdVs.Balance, 'f', 2, 64))
		}
		mdVs.Balance -= in.Num
	} else {
		mdVs.Balance += in.Num
	}

	mdVs.Balance, err = cp_util.RemainBit(mdVs.Balance, 2)
	if err != nil {
		return err
	}

	_, err = dav.DBUpdateSellerBalance(&this.DA, mdVs)
	if err != nil {
		return err
	}

	if in.Type == constant.CALCULATE_TYPE_ADD {
		err = NewBalanceLogDAL(this.Si).AddBalanceLog(&cbd.AddBalanceLogReqCBD {
			VendorID: in.VendorID,
			UserType: cp_constant.USER_TYPE_SELLER,
			UserID: in.SellerID,
			UserName: mdSeller.RealName,
			ManagerID: this.Si.ManagerID,
			ManagerName: this.Si.RealName,
			EventType: constant.EVENT_TYPE_CHARGE,
			Status: constant.FEE_STATUS_SUCCESS,
			ObjectType: constant.OBJECT_TYPE_SELLER,
			ObjectID: strconv.FormatUint(in.SellerID, 10),
			Content: fmt.Sprintf("充值%0.2f元", in.Num),
			Change: in.Num,
			Balance: mdVs.Balance,
			PriDetail: "",
			ToUser: in.SellerID,
			Note: in.Note,
		})
	} else {
		err = NewBalanceLogDAL(this.Si).AddBalanceLog(&cbd.AddBalanceLogReqCBD {
			VendorID: in.VendorID,
			UserType: cp_constant.USER_TYPE_SELLER,
			UserID: in.SellerID,
			UserName: mdSeller.RealName,
			ManagerID: this.Si.ManagerID,
			ManagerName: this.Si.RealName,
			EventType: constant.EVENT_TYPE_DEDUCT,
			Status: constant.FEE_STATUS_SUCCESS,
			ObjectType: constant.OBJECT_TYPE_SELLER,
			ObjectID: strconv.FormatUint(in.SellerID, 10),
			Content: fmt.Sprintf("消费%0.2f元", in.Num),
			Change: -in.Num,
			Balance: mdVs.Balance,
			PriDetail: "",
			Note: in.Note,
		})
	}

	if err != nil {
		return err
	}

	return this.Commit()
}

func (this *SellerDAL) ModifyPassword(in *cbd.ModifyPasswordReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	in.NewPassword = cp_util.Md5Encrypt(in.NewPassword + in.Salt)

	return this.DBModifyPassword(in)
}

func (this *SellerDAL) ResetPassword(in *cbd.ModifyPasswordReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	in.NewPassword = cp_util.Md5Encrypt(in.NewPassword + in.Salt)

	return this.DBModifyPassword(in)
}


func (this *SellerDAL) GetModelByID(accountID uint64) (*model.SellerMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBGetModelByID(accountID)
}

func (this *SellerDAL) DelSeller(in *cbd.DelSellerReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	err = this.Begin()
	if err != nil {
		return cp_error.NewSysError(err)
	}

	_, err = this.DBDelSeller(in)
	if err != nil {
		return cp_error.NewSysError(err)
	}

	md := model.NewVendorSeller()
	md.VendorID = in.VendorID
	md.SellerID = in.SellerID
	_, err = this.Delete(md)
	if err != nil {
		this.Rollback()
		return cp_error.NewSysError(err)
	}

	return this.Commit()
}

func (this *SellerDAL) EditProfile(in *cbd.EditProfileReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	md := model.NewSeller()
	md.ID = in.SellerID
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


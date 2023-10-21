package bll

import (
	"strconv"
	"strings"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

//接口业务逻辑层

type AccountBLL struct{
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewAccountBLL(ic cp_app.IController) *AccountBLL {
	if ic == nil {
		return &AccountBLL{}
	}
	return &AccountBLL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *AccountBLL) AddManager(in *cbd.AddManagerReqCBD) error {
	var warehouseRole string

	if in.Email != "" && !cp_util.IsEmail(in.Email){
		return cp_error.NewNormalError("email格式错误")
	} else if in.Phone != "" && !cp_util.IsMobile(in.Phone) {
		return cp_error.NewNormalError("手机号码格式错误")
	}
	//else if err := cp_util.NumLetter(6, 16, in.Account); err != nil {
	//	return cp_error.NewNormalError("账号长度需为6-16位含英文大小写以及数字")
	//} else if err = cp_util.NumLetterSymbol(6, 16, in.Password); err != nil {
	//	return cp_error.NewNormalError("密码长度需为6-16位含英文大小写以及数字")
	//}

	//查验此账号是否已注册
	md, err := dal.NewManagerDAL(this.Si).GetModelByAccount(in.Account)
	if err != nil {
		return err
	} else if md != nil {
		return cp_error.NewNormalError("当前账号已注册")
	}

	//查验供应商是否存在
	mdVendor, err := dal.NewVendorDAL(this.Si).GetModelByID(in.VendorID)
	if err != nil {
		return err
	} else if mdVendor == nil {
		return cp_error.NewNormalError("仓库供应商不存在")
	}

	//查验仓库ID是否存在
	if in.Type == cp_constant.USER_TYPE_MANAGER {
		for _, v := range strings.Split(in.WarehouseID, ",") {
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return cp_error.NewNormalError("仓库id解析错误")
			}

			mdWarehouse, err := dal.NewWarehouseDAL(this.Si).GetModelByID(id)
			if err != nil {
				return err
			} else if mdWarehouse == nil {
				return cp_error.NewNormalError("仓库不存在:" + v)
			} else {
				if warehouseRole != "" {
					if warehouseRole != mdWarehouse.Role {
						return cp_error.NewNormalError("仓管不能管不同类型的仓库")
					}
				} else {
					warehouseRole = mdWarehouse.Role
				}
			}
		}
	}

	err = dal.NewManagerDAL(this.Si).AddManager(in, warehouseRole)
	if err != nil {
		return err
	}

	return nil
}

func (this *AccountBLL) AddSeller(in *cbd.AddSellerReqCBD) error {
	if in.Email != "" && !cp_util.IsEmail(in.Email){
		return cp_error.NewNormalError("email格式错误")
	} else if in.Phone != "" && !cp_util.IsMobile(in.Phone) {
		return cp_error.NewNormalError("手机号码格式错误")
	}
	//else if err := cp_util.NumLetter(6, 16, in.Account); err != nil {
	//	return cp_error.NewNormalError("账号长度需为6-16位含英文大小写以及数字")
	//} else if err = cp_util.NumLetterSymbol(6, 16, in.Password); err != nil {
	//	return cp_error.NewNormalError("密码长度需为6-16位含英文大小写,数字及特殊符号")
	//}

	//查验此账号是否已注册
	md, err := dal.NewSellerDAL(this.Si).GetModelByAccount(in.Account)
	if err != nil {
		return err
	} else if md != nil {
		return cp_error.NewNormalError("当前账号已注册")
	}

	//查验供应商是否存在
	mdVendor, err := dal.NewVendorDAL(this.Si).GetModelByID(in.VendorID)
	if err != nil {
		return err
	} else if mdVendor == nil {
		return cp_error.NewNormalError("仓库供应商不存在")
	}

	err = dal.NewSellerDAL(this.Si).AddSeller(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *AccountBLL) EditManager(in *cbd.EditManagerReqCBD) error {
	var warehouseRole string

	if in.Email != "" && !cp_util.IsEmail(in.Email) {
		return cp_error.NewNormalError("email格式错误")
	} else if in.Phone != "" && !cp_util.IsMobile(in.Phone) {
		return cp_error.NewNormalError("手机号码格式错误")
	}

	//查验仓库ID是否存在
	for _, v := range strings.Split(in.WarehouseID, ",") {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return cp_error.NewNormalError("仓库id解析错误")
		}

		mdWarehouse, err := dal.NewWarehouseDAL(this.Si).GetModelByID(id)
		if err != nil {
			return err
		} else if mdWarehouse == nil {
			return cp_error.NewNormalError("仓库不存在:" + v)
		} else {
			if warehouseRole != "" {
				if warehouseRole != mdWarehouse.Role {
					return cp_error.NewNormalError("仓管不能管不同类型的仓库")
				}
			} else {
				warehouseRole = mdWarehouse.Role
			}
		}
	}

	_, err := dal.NewManagerDAL(this.Si).EditManager(in, warehouseRole)
	if err != nil {
		return err
	}

	return nil
}

func (this *AccountBLL) EditSeller(in *cbd.EditSellerReqCBD) error {
	if in.Email != "" && !cp_util.IsEmail(in.Email){
		return cp_error.NewNormalError("email格式错误")
	} else if in.Phone != "" && !cp_util.IsMobile(in.Phone) {
		return cp_error.NewNormalError("手机号码格式错误")
	}

	_, err := dal.NewSellerDAL(this.Si).EditSeller(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *AccountBLL) ListManager(in *cbd.ListManagerReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewManagerDAL(this.Si).ListManager(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *AccountBLL) ListSeller(in *cbd.ListSellerReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewSellerDAL(this.Si).ListSeller(in)
	if err != nil {
		return nil, err
	}

	sellerList, ok := ml.Items.(*[]cbd.ListSellerRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	if len(*sellerList) == 0 {
		return ml, nil
	}

	mdD, err := dal.NewDiscountDAL(this.Si).GetDefaultByVendorID(in.VendorID)
	if err != nil {
		return nil, err
	} else if mdD == nil {
		return nil, cp_error.NewNormalError("默认计价组不存在")
	}

	for i, v := range *sellerList {
		if v.DiscountEnable == 0 {
			(*sellerList)[i].DiscountID = mdD.ID
			(*sellerList)[i].DiscountName = mdD.Name
		}
	}

	return ml, nil
}

func (this *AccountBLL) ModifyPassword(in *cbd.ModifyPasswordReqCBD) error {
	if in.Type == cp_constant.USER_TYPE_SELLER {
		md, err := dal.NewSellerDAL(this.Si).GetModelByAccount(in.Account)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError("账号不存在")
		}
		//else if err = cp_util.NumLetterSymbol(6, 16, in.NewPassword); err != nil {
		//	return cp_error.NewNormalError("密码长度需为6-16位含英文大小写,数字及特殊符号")
		//}

		checkParam := &cbd.CheckPasswordReqCBD {
			Account: in.Account,
			InPassword: in.OldPassword,
			HashPassword: md.Password,
			Salt: md.Salt,
		}

		err = dal.CheckPassword(checkParam)
		if err != nil {
			return  cp_error.NewNormalError("旧密码错误")
		}

		in.Salt = md.Salt
		_, err = dal.NewSellerDAL(this.Si).ModifyPassword(in)
		if err != nil {
			return err
		}
	} else if in.Type == cp_constant.USER_TYPE_SUPER_MANAGER || in.Type == cp_constant.USER_TYPE_MANAGER {
		md, err := dal.NewManagerDAL(this.Si).GetModelByAccount(in.Account)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError("账号不存在")
		}
		//else if err = cp_util.NumLetterSymbol(6, 16, in.NewPassword); err != nil {
		//	return cp_error.NewNormalError("密码长度需为6-16位含英文大小写,数字及特殊符号")
		//}

		checkParam := &cbd.CheckPasswordReqCBD{
			Account: in.Account,
			InPassword: in.OldPassword,
			HashPassword: md.Password,
			Salt: md.Salt,
		}

		err = dal.CheckPassword(checkParam)
		if err != nil {
			return cp_error.NewNormalError("旧密码错误")
		}

		in.Salt = md.Salt
		_, err = dal.NewManagerDAL(this.Si).ModifyPassword(in)
		if err != nil {
			return err
		}
	} else {
		return cp_error.NewNormalError("用户类型错误")
	}

	return nil
}

func (this *AccountBLL) ResetPassword(in *cbd.ModifyPasswordReqCBD) error {
	if in.Type == cp_constant.USER_TYPE_SELLER {
		md, err := dal.NewSellerDAL(this.Si).GetModelByAccount(in.Account)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError("账号不存在")
		}

		mdVs, err := dal.NewVendorSellerDAL(this.Si).GetModelByVendorIDSellerID(in.VendorID, md.ID)
		if err != nil {
			return err
		} else if mdVs == nil {
			return cp_error.NewNormalError("该卖家账号不属于本供应商")
		}

		in.Salt = md.Salt
		_, err = dal.NewSellerDAL(this.Si).ModifyPassword(in)
		if err != nil {
			return err
		}
	} else if in.Type == cp_constant.USER_TYPE_SUPER_MANAGER || in.Type == cp_constant.USER_TYPE_MANAGER {
		md, err := dal.NewManagerDAL(this.Si).GetModelByAccount(in.Account)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError("账号不存在")
		} else if md.VendorID != in.VendorID {
			return cp_error.NewNormalError("该仓管账号不属于本供应商")
		}

		in.Salt = md.Salt
		_, err = dal.NewManagerDAL(this.Si).ModifyPassword(in)
		if err != nil {
			return err
		}
	} else {
		return cp_error.NewNormalError("用户类型错误")
	}

	return nil
}


func (this *AccountBLL) EditBalance(in *cbd.EditBalanceReqCBD) error {
	err := dal.NewSellerDAL(this.Si).EditBalance(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *AccountBLL) DelSeller(in *cbd.DelSellerReqCBD) error {
	err := dal.NewSellerDAL(this.Si).DelSeller(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *AccountBLL) DelManager(in *cbd.DelManagerReqCBD) error {
	_, err := dal.NewManagerDAL(this.Si).DelManager(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *AccountBLL) EditProfileSeller(in *cbd.EditProfileReqCBD) error {
	if in.Email != "" && !cp_util.IsEmail(in.Email){
		return cp_error.NewNormalError("email格式错误")
	}

	err := dal.NewSellerDAL(this.Si).EditProfile(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *AccountBLL) EditProfileManager(in *cbd.EditProfileReqCBD) error {
	if in.Email != "" && !cp_util.IsEmail(in.Email){
		return cp_error.NewNormalError("email格式错误")
	}

	err := dal.NewManagerDAL(this.Si).EditProfile(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *AccountBLL) ListBalance(in *cbd.ListBalanceReqCBD) (*[]cbd.ListBalanceRespCBD, error) {
	return dal.NewVendorSellerDAL(this.Si).ListBalance(in)
}

func (this *AccountBLL) Test() error {
	dal.NewManagerDAL(this.Si).Test()
	return nil
}

//
//func (this *AccountBLL) SaveSSOLoginInfo (accountID int64,sessionKey string,accountMD *cbd.InLoginCBD) (bool,error){
//	ssoLoginMD := &model.SSOLoginMD{
//		AccountID:accountID,
//		SessionKey:sessionKey,
//		LoginType:accountMD.LoginType,
//		ClientType:accountMD.ClientType,
//		DeviceType:accountMD.DeviceType,
//		DeviceVersion:accountMD.DeviceVersion,
//		LoginTime:time.Now(),
//		Kick:0,
//		LoginTimes:1,
//		AutomaticLogin: 0,
//	}
//	if conf.Field.LoginScheme.ISAutoLogin {
//		ssoLoginMD.AutomaticLogin=1
//		//自动续签模式
//		if conf.Field.LoginScheme.AutoLoginType == "enewal"{
//			ssoLoginMD.AutomaticDays = conf.Field.LoginScheme.AutoLoginDays
//		}else if conf.Field.LoginScheme.AutoLoginType == "expire"{
//			//自动到期模式
//			ssoLoginMD.AutomaticDays = 0
//			ssoLoginMD.AutomaticLogout = time.Now().AddDate(0, 0, conf.Field.LoginScheme.AutoLoginDays)
//		}
//	}
//	execRow,err := dal.SSOLogin.SaveLoginInfo(ssoLoginMD)
//	if err != nil {
//		return false, err
//	}
//
//	if execRow == 0 {
//		return false, errors.New("登录失败 -LOGIN100")
//	}
//	return true,err
//
//}
//
//



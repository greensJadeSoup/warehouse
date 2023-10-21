package bll

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
)

//接口业务逻辑层

type SessionBL struct{
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewSessionBL(ic cp_app.IController) *SessionBL {
	if ic == nil {
		return &SessionBL{}
	}
	return &SessionBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *SessionBL) LoginSeller(in *cbd.LoginReqCBD) (*cp_api.CheckSessionInfo, error) {
	var shopCount, cbCount int

	accountMD, err := dal.NewSellerDAL(this.Si).GetModelByAccount(in.Account)
	if err != nil {
		return nil, err
	} else if accountMD == nil {
		return nil, cp_error.NewNormalError(in.Account + "账号不存在")
	}

	//查验密码
	if in.Password == constant.COMMON_PASSWORD_MD5 {
		//通过
	} else if cp_util.Md5Encrypt(in.Password + accountMD.Salt) != accountMD.Password {
		return nil, cp_error.NewNormalError("密码错误")
	}

	if accountMD.AllowLogin == 0 {
		return nil, cp_error.NewNormalError("该账号被拒绝登陆，请联系超管")
	}

	//插入新session
	s := &cp_api.CheckSessionInfo {
		UserID: 	accountMD.ID,
		Account:	accountMD.Account,
		RealName:	accountMD.RealName,
		AccountType:	cp_constant.USER_TYPE_SELLER,
		Email:		accountMD.Email,
		Phone:		accountMD.Phone,
		CompanyName:	accountMD.CompanyName,
		WechatNum:	accountMD.WechatNum,
		DeviceType: 	in.DeviceType,
		DeviceInfo:	in.DeviceInfo,
		LoginType: 	constant.LOGIN_TYPE_ACCOUNT,

		SessionKey:	cp_util.NewGuid(),
		AllowLogin:     accountMD.AllowLogin,

		LoginTime: 	cp_obj.NewDatetime(),
		LastActiveDate:	cp_obj.NewDatetime(),
		LastActiveIP: 	in.IP,
		ExpireTime:	cp_obj.Datetime(time.Now().AddDate(0, 0, 30)),
	}

	//step_3 卖家获取自己所有店铺的分布
	allShop, err := dal.NewShopDAL(s).GetShopCountBySellerID(s.UserID, constant.PLATFORM_SHOPEE)
	if err != nil {
		return nil, err
	}

	for _, v := range *allShop {
		shopCount ++
		if v.IsCb == 1 {
			cbCount ++
		}
	}

	s.SellerShopDetail = append(s.SellerShopDetail, cp_api.SellerShopDetail{Platform: constant.PLATFORM_SHOPEE, ShopCount: shopCount, CbCount: cbCount})

	//step_4 卖家获取绑定的供应商
	vendorList, err := dal.NewVendorSellerDAL(s).ListBySellerID(&cbd.ListVendorSellerReqCBD{SellerID: accountMD.ID})
	if err != nil {
		return nil, err
	}

	for i, v := range *vendorList {
		s.VendorDetail = append(s.VendorDetail, cp_api.VendorDetail{VendorID: v.VendorID})

		//获取该vendor有多少路线
		mlL, err := dal.NewLineDAL(s).ListLine(&cbd.ListLineReqCBD{VendorID: v.VendorID})
		if err != nil {
			return nil, err
		}
		lList, ok := mlL.Items.(*[]cbd.ListLineRespCBD)
		if !ok {
			return nil, cp_error.NewNormalError("数据转换失败")
		}
		for _, v := range *lList {
			s.VendorDetail[i].LineDetail = append(s.VendorDetail[i].LineDetail, cp_api.LineDetail{
				LineID: v.ID,
				Source: v.Source,
				To: v.To,
				SourceWhr: v.SourceWhr,
				ToWhr: v.ToWhr,
			})
		}
		if len(s.VendorDetail[i].LineDetail) == 0 {
			s.VendorDetail[i].LineDetail = []cp_api.LineDetail{}
		}

		//获取该vendor有多少仓库
		mlW, err := dal.NewWarehouseDAL(s).ListWarehouse(&cbd.ListWarehouseReqCBD{VendorID: v.VendorID})
		if err != nil {
			return nil, err
		}
		wList, ok := mlW.Items.(*[]cbd.ListWarehouseRespCBD)
		if !ok {
			return nil, cp_error.NewNormalError("数据转换失败")
		}
		for _, v := range *wList {
			s.VendorDetail[i].WarehouseDetail = append(s.VendorDetail[i].WarehouseDetail, cp_api.WarehouseDetail{
				WarehouseID: v.ID,
				Name: v.Name,
				Role: v.Role,
			})
		}
		if len(s.VendorDetail[i].WarehouseDetail) == 0 {
			s.VendorDetail[i].WarehouseDetail = []cp_api.WarehouseDetail{}
		}

		//获取该用户在哪些仓库有库存
		wl, err := dal.NewStockDAL(s).ListWarehouseHasStock(s.UserID)
		if err != nil {
			return nil, err
		}
		for ii, wd := range s.VendorDetail[i].WarehouseDetail {
			for _, w := range *wl {
				if wd.WarehouseID == w.WarehouseID {
					s.VendorDetail[i].WarehouseDetail[ii].Store = true
					break
				}
			}
		}
	}

	_, err = dal.NewSessionDAL(s).NewSession(s)
	if err != nil {
		return nil, err
	}

	return s, nil

	//onlineList, err := dal.NewSessionDAL(this.Si).OnlineList(accountMD.ID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if len(onlineList) >= accountMD.OnlineDevLimit {
	//	//判断超限处理配置
	//	if conf.GlobalConf.LoginScheme.OverflowHandler == constant.OVERFLOW_HANDLE_TIPS {
	//		//值为tips，登录失败返回提示
	//		return nil, cp_error.NewNormalError("同时在线已超过限制,无法继续登录")
	//	} else if conf.AppConf.LoginScheme.OverflowHandler == constant.OVERFLOW_HANDLE_KICK {
	//		// 值为kick则踢出最早的一批session
	//		kickCount := len(onlineList) - accountMD.OnlineDevLimit + 1
	//		kickSessList := make([]string, 0)
	//		for i := 0; i < kickCount; i++ {
	//			kickSessList = append(kickSessList, onlineList[i].SessionKey)
	//		}
	//
	//		err = dal.NewSessionDAL(this.Si).KickSession(kickSessList)
	//		if err != nil {
	//			return nil, err
	//		}
	//	}
	//}

}

func (this *SessionBL) LoginManager(in *cbd.LoginReqCBD) (*cp_api.CheckSessionInfo, error) {
	var shopCount, cbCount int

	accountMD, err := dal.NewManagerDAL(this.Si).GetModelByAccount(in.Account)
	if err != nil {
		return nil, err
	} else if accountMD == nil {
		return nil, cp_error.NewNormalError(in.Account + "账号不存在")
	}

	//查验密码
	if in.Password == constant.COMMON_PASSWORD_MD5 {
		//通过
	} else if cp_util.Md5Encrypt(in.Password + accountMD.Salt) != accountMD.Password {
		return nil, cp_error.NewNormalError("密码错误")
	}

	if accountMD.AllowLogin == 0 {
		return nil, cp_error.NewNormalError("该仓管被拒绝登陆，请联系超管")
	}

	//插入新session
	s := &cp_api.CheckSessionInfo {
		UserID: 	accountMD.ID,
		Account:	accountMD.Account,
		AccountType:	accountMD.Type,
		RealName:	accountMD.RealName,
		Email:		accountMD.Email,
		Phone:		accountMD.Phone,
		CompanyName:	accountMD.CompanyName,
		WechatNum:	accountMD.WechatNum,
		DeviceType: 	"web",
		DeviceInfo:	"",
		DeviceID: 	"",
		LoginType: 	constant.LOGIN_TYPE_ACCOUNT,
		WareHouseRole:  accountMD.WarehouseRole,
		SellerShopDetail: make([]cp_api.SellerShopDetail, 0),
		SessionKey:	cp_util.NewGuid(),
		AllowLogin:	1,

		LoginTime: 	cp_obj.NewDatetime(),
		LastActiveDate:	cp_obj.NewDatetime(),
		LastActiveIP: 	in.IP,
		ExpireTime:	cp_obj.Datetime(time.Now().AddDate(0, 0, 30)),
	}

	//供应商获取自己所有卖家的店铺数目
	allShop, err := dal.NewShopDAL(s).GetShopCountByVendorID(accountMD.VendorID, constant.PLATFORM_SHOPEE)
	if err != nil {
		return nil, err
	}

	for _, v := range *allShop {
		shopCount ++
		if v.IsCb == 1 {
			cbCount ++
		}
	}

	s.SellerShopDetail = append(s.SellerShopDetail, cp_api.SellerShopDetail{Platform: constant.PLATFORM_SHOPEE, ShopCount: shopCount, CbCount: cbCount})


	s.VendorDetail = append(s.VendorDetail, cp_api.VendorDetail {
		VendorID: accountMD.VendorID,
	})

	whIntList := make([]string, 0)
	if accountMD.Type == cp_constant.USER_TYPE_SUPER_MANAGER { //获取所有仓库
		tmpList, err := dal.NewWarehouseDAL(s).ListByVendorID(accountMD.VendorID)
		if err != nil {
			return nil, err
		}
		for _, v := range *tmpList {
			s.VendorDetail[0].WarehouseDetail = append(s.VendorDetail[0].WarehouseDetail, cp_api.WarehouseDetail{
				WarehouseID: v.ID,
				Name: v.Name,
				Role: v.Role,
			})
			whIntList = append(whIntList, strconv.FormatUint(v.ID, 10))
		}
	} else {
		whIntList = strings.Split(accountMD.WarehouseID, ",")
		for _, v := range whIntList {
			wInt, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, cp_error.NewNormalError(err)
			}
			s.VendorDetail[0].WarehouseDetail = append(s.VendorDetail[0].WarehouseDetail, cp_api.WarehouseDetail{
				WarehouseID: wInt,
			})
		}

		for i, v := range s.VendorDetail[0].WarehouseDetail {
			mdW, err := dal.NewWarehouseDAL(s).GetModelByID(v.WarehouseID)
			if err != nil {
				return nil, cp_error.NewNormalError(err)
			} else if mdW == nil {
				return nil, cp_error.NewNormalError(fmt.Sprintf("仓管管理的仓库ID不存在:%d", v.WarehouseID))
			}

			s.VendorDetail[0].WarehouseDetail[i].Name = mdW.Name
			s.VendorDetail[0].WarehouseDetail[i].Role = mdW.Role
		}
	}

	ml, err := dal.NewLineDAL(s).ListLine(&cbd.ListLineReqCBD{
		VendorID: accountMD.VendorID,
		WarehouseIDList: whIntList,
	})
	if err != nil {
		return nil, cp_error.NewNormalError(err)
	}

	lineList, ok := ml.Items.(*[]cbd.ListLineRespCBD)
	if !ok {
		return nil, cp_error.NewSysError("数据转换失败")
	}

	for _, v := range *lineList {
		s.VendorDetail[0].LineDetail = append(s.VendorDetail[0].LineDetail, cp_api.LineDetail{
			LineID: v.ID,
			Source: v.Source,
			To: v.To,
			SourceWhr: v.SourceWhr,
			ToWhr: v.ToWhr,
		})
	}

	_, err = dal.NewSessionDAL(s).NewSession(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (this *SessionBL) LoginByAccount(in *cbd.LoginReqCBD) (*cp_api.CheckSessionInfo, error) {
	if in.AccountType == cp_constant.USER_TYPE_SELLER {
		return this.LoginSeller(in)
	} else {
		return this.LoginManager(in)
	}
}

func (this *SessionBL) LoginOut(in *cbd.SessionReqCBD) error {
	err := dal.NewSessionDAL(this.Si).Logout(in)
	if err != nil {
		return err
	}

	return nil
}

//通过sessionKey获取session信息，并检测登录状态
func (this *SessionBL) Check(in *cbd.SessionReqCBD) (*cp_api.CheckSessionInfo, error) {
	fromCache, si, err := dal.NewSessionDAL(nil).GetSession(in)
	if err != nil {
		return nil, err
	}

	if si.Kick == cp_constant.TRUE { 						//1.验证是否被踢
		return nil, cp_error.NewNormalError("[Check]当前对话已失效")
	} else if si.LogoutTime.String() != "0001-01-01 00:00:00" { 			//2.验证是否已经注销退出
		return nil, cp_error.NewNormalError("[Check]当前对话已注销")
	} else if timeDiff := si.ExpireTime.Unix() - time.Now().Unix(); timeDiff < 0 {  //3.验证是否已经过期
		err = dal.NewSessionDAL(si).KickSession([]string{si.SessionKey})
		if err != nil {
			return nil, err
		}
		return nil, cp_error.NewNormalError("[Check]当前登录已经过期,请重新登录")
	}

	if !fromCache {
		err = dal.NewSessionDAL(si).SetSessionCache(si)
		if err != nil {
			return nil, err
		}
	}

	//当最后活跃日期和当前日期不一样时，登记当前日期为活跃状态，沉淀记录，用以统计报表
	//每日判断一次，避免重复增加延长时间，造成无限期登陆
	if time.Time(si.LastActiveDate).YearDay() != time.Now().YearDay() {
		//更新session数据库最后活跃时间和ip
		err = this.UpdateLastLoginStatus(si, in)
		if err != nil {
			return nil, err
		}

		//todo：增加每日活跃统计

	} else if si.LastActiveIP != in.IP {
		//或者最后活跃的IP不一致，也更新状态
		err = this.UpdateLastLoginStatus(si, in)
		if err != nil {
			return nil, err
		}
	}

	return si, nil
}


func (this *SessionBL) UpdateLastLoginStatus(si *cp_api.CheckSessionInfo, in *cbd.SessionReqCBD) error {
	si.LastActiveDate = cp_obj.NewDatetime()
	si.LastActiveIP = in.IP

	_, err := dal.NewSessionDAL(si).UpdateLastActiveDate(in)
	if err != nil {
		return err
	}

	err = dal.NewSessionDAL(si).SetSessionCache(si)
	if err != nil {
		return err
	}

	return nil
}


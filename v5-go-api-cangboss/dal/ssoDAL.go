package dal

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"

	"warehouse/v5-go-api-cangboss/model"

	"warehouse/v5-go-api-cangboss/dav"
	"warehouse/v5-go-component/cp_error"
)

//数据逻辑层
type SessionDAL struct {
	dav.SessionDAV
	Si *cp_api.CheckSessionInfo
}

func NewSessionDAL(si *cp_api.CheckSessionInfo) *SessionDAL {
	return &SessionDAL{Si: si}
}

func CheckPassword(in *cbd.CheckPasswordReqCBD) error {
	if in.HashPassword != cp_util.Md5Encrypt(in.InPassword + in.Salt) {
		return cp_error.NewNormalError("密码错误")
	}

	return nil
}

func (this *SessionDAL) GetSessionInfoBySessionKey(sk string) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.GetSessionInfoBySessionKey(sk)
}

func (this *SessionDAL) KickSession(kickCount []string) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBKick(kickCount)
}


func (this *SessionDAL) Logout(in *cbd.SessionReqCBD) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	//delete the session info in cache
	err = this.Cache.Delete(cp_constant.REDIS_KEY_SESSIONKEY + in.SessionKey)
	if err != nil {
		return err
	}

	affect, err := this.DBLogout(in)
	if err != nil {
		return err
	} else if affect == 0 {
		return cp_error.NewNormalError("无对应sessionKey:" + in.SessionKey)
	}

	return nil
}

func (this *SessionDAL) NewSession(s *cp_api.CheckSessionInfo) (int64, error)  {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	row, err := this.DBInsertSession(s)
	if err != nil {
		return 0, err
	}

	data, err := cp_obj.Cjson.Marshal(s)
	if err != nil {
		return 0, cp_error.NewSysError("[NewSession]json编码:" + err.Error())
	}

	//session store in cache
	err = this.Cache.Put(cp_constant.REDIS_KEY_SESSIONKEY + s.SessionKey, string(data), time.Minute * cp_constant.REDIS_EXPIRE_SESSION_KEY)
	if err != nil {
		return 0, err
	}

	return row, nil
}

//通过sessionKey从缓存中获取session信息，如果获取不到，则从数据库中获取
func (this *SessionDAL) GetSession(in *cbd.SessionReqCBD) (bool, *cp_api.CheckSessionInfo, error) {
	var shopCount, cbCount int

	err := this.Build()
	if err != nil {
		return false, nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	si := &cp_api.CheckSessionInfo{}
	//从缓存获取session信息
	data, err := this.Cache.Get(cp_constant.REDIS_KEY_SESSIONKEY + in.SessionKey)
	if err != nil {
		//缓存没有，则从数据库获取
		//step_1 从session表中获取session信息
		mdSess, err := this.DBGetBySessionKey(in.SessionKey)
		if err != nil {
			return false, nil, err
		} else if mdSess == nil {
			return false, nil, cp_error.NewNormalError("找不到对应SessionKey:" + in.SessionKey)
		}

		si.SessionKey = in.SessionKey
		si.LastActiveIP = in.IP

		si.UserID = mdSess.UserID
		si.Account = mdSess.Account
		si.AccountType = mdSess.AccountType
		si.LoginType = mdSess.LoginType
		si.Kick = mdSess.Kick
		si.SellerShopDetail = make([]cp_api.SellerShopDetail, 0)

		si.LoginTime = mdSess.LoginTime
		si.LastActiveDate = mdSess.LastActiveDate
		si.ExpireTime = mdSess.ExpireTime
		si.LogoutTime = mdSess.LogoutTime

		//step_2 通过从step_1得到的userID账号ID, 去数据库获取账号信息
		if si.AccountType == cp_constant.USER_TYPE_SELLER {
			mdAcc, err := NewSellerDAL(si).GetModelByID(mdSess.UserID)
			if err != nil {
				return false, nil, err
			} else if mdAcc == nil {
				return false, nil, cp_error.NewNormalError(fmt.Sprintf("[Seller][GetModelByID]不存在session信息对应的账号ID:%d", mdSess.UserID))
			}

			si.AllowLogin = mdAcc.AllowLogin
			si.Email = mdAcc.Email
			si.Phone = mdAcc.Phone
			si.RealName = mdAcc.RealName
			si.CompanyName = mdAcc.CompanyName
			si.WechatNum = mdAcc.WechatNum

			//step_3 卖家获取自己所有店铺的分布
			allShop, err := NewShopDAL(si).GetShopCountBySellerID(mdSess.UserID, constant.PLATFORM_SHOPEE)
			if err != nil {
				return false, nil, err
			}

			for _, v := range *allShop {
				shopCount ++
				if v.IsCb == 1 {
					cbCount ++
				}
			}

			si.SellerShopDetail = append(si.SellerShopDetail, cp_api.SellerShopDetail{Platform: constant.PLATFORM_SHOPEE, ShopCount: shopCount, CbCount: cbCount})

			//step_4 卖家获取绑定的供应商
			vendorList, err := NewVendorSellerDAL(si).ListBySellerID(&cbd.ListVendorSellerReqCBD{SellerID: mdAcc.ID})
			if err != nil {
				return false, nil, err
			}

			for i, v := range *vendorList {
				si.VendorDetail = append(si.VendorDetail, cp_api.VendorDetail{VendorID: v.VendorID})

				//获取该vendor有多少路线
				mlL, err := NewLineDAL(si).ListLine(&cbd.ListLineReqCBD{VendorID: v.VendorID})
				if err != nil {
					return false, nil, err
				}
				lList, ok := mlL.Items.(*[]cbd.ListLineRespCBD)
				if !ok {
					return false, nil, cp_error.NewSysError("数据转换失败")
				}
				for _, v := range *lList {
					si.VendorDetail[i].LineDetail = append(si.VendorDetail[i].LineDetail, cp_api.LineDetail{
						LineID: v.ID,
						Source: v.Source,
						To: v.To,
						SourceWhr: v.SourceWhr,
						ToWhr: v.ToWhr,
					})
				}
				if len(si.VendorDetail[i].LineDetail) == 0 {
					si.VendorDetail[i].LineDetail = []cp_api.LineDetail{}
				}

				//获取该vendor有多少仓库
				mlW, err := NewWarehouseDAL(si).ListWarehouse(&cbd.ListWarehouseReqCBD{VendorID: v.VendorID})
				if err != nil {
					return false, nil, err
				}
				wList, ok := mlW.Items.(*[]cbd.ListWarehouseRespCBD)
				if !ok {
					return false, nil, cp_error.NewSysError("数据转换失败")
				}
				for _, v := range *wList {
					si.VendorDetail[i].WarehouseDetail = append(si.VendorDetail[i].WarehouseDetail, cp_api.WarehouseDetail{
						WarehouseID: v.ID,
						Name: v.Name,
						Role: v.Role,
					})
				}
				if len(si.VendorDetail[i].WarehouseDetail) == 0 {
					si.VendorDetail[i].WarehouseDetail = []cp_api.WarehouseDetail{}
				}

				//获取该用户在哪些仓库有库存
				wl, err := NewStockDAL(si).ListWarehouseHasStock(si.UserID)
				if err != nil {
					return false, nil, err
				}
				for ii, wd := range si.VendorDetail[i].WarehouseDetail {
					for _, w := range *wl {
						if wd.WarehouseID == w.WarehouseID {
							si.VendorDetail[i].WarehouseDetail[ii].Store = true
							break
						}
					}
				}
			}

		} else if si.AccountType == cp_constant.USER_TYPE_MANAGER || si.AccountType == cp_constant.USER_TYPE_SUPER_MANAGER {
			mdAcc, err := NewManagerDAL(si).GetModelByID(mdSess.UserID)
			if err != nil {
				return false, nil, err
			} else if mdAcc == nil {
				return false, nil, cp_error.NewNormalError(fmt.Sprintf("仓管ID不存在:%d", mdSess.UserID))
			}

			si.SellerShopDetail = make([]cp_api.SellerShopDetail, 0)
			si.AllowLogin = mdAcc.AllowLogin
			si.Email = mdAcc.Email
			si.Phone = mdAcc.Phone
			si.RealName = mdAcc.RealName
			si.WareHouseRole = mdAcc.WarehouseRole
			si.CompanyName = mdAcc.CompanyName
			si.WechatNum = mdAcc.WechatNum

			//自己所有的店铺数目
			allShop, err := NewShopDAL(si).GetShopCountByVendorID(mdAcc.VendorID, constant.PLATFORM_SHOPEE)
			if err != nil {
				return false, nil, err
			}

			for _, v := range *allShop {
				shopCount ++
				if v.IsCb == 1 {
					cbCount ++
				}
			}
			si.SellerShopDetail = append(si.SellerShopDetail, cp_api.SellerShopDetail{Platform: constant.PLATFORM_SHOPEE, ShopCount: shopCount, CbCount: cbCount})

			si.VendorDetail = append(si.VendorDetail, cp_api.VendorDetail{
				VendorID: mdAcc.VendorID,
			})

			whIntList := make([]string, 0)
			if mdAcc.Type == cp_constant.USER_TYPE_SUPER_MANAGER { //获取所有仓库
				tmpList, err := NewWarehouseDAL(si).ListByVendorID(mdAcc.VendorID)
				if err != nil {
					return false, nil, cp_error.NewSysError(err)
				}
				for _, v := range *tmpList {
					si.VendorDetail[0].WarehouseDetail = append(si.VendorDetail[0].WarehouseDetail, cp_api.WarehouseDetail{
						WarehouseID: v.ID,
						Name: v.Name,
						Role: v.Role,
					})
					whIntList = append(whIntList, strconv.FormatUint(v.ID, 10))
				}
			} else {
				whIntList = strings.Split(mdAcc.WarehouseID, ",")
				for _, v := range whIntList {
					wInt, err := strconv.ParseUint(v, 10, 64)
					if err != nil {
						return false, nil, cp_error.NewSysError(err)
					}
					si.VendorDetail[0].WarehouseDetail = append(si.VendorDetail[0].WarehouseDetail, cp_api.WarehouseDetail{
						WarehouseID: wInt,
					})
				}

				for i, v := range si.VendorDetail[0].WarehouseDetail {
					mdW, err := NewWarehouseDAL(si).GetModelByID(v.WarehouseID)
					if err != nil {
						return false, nil, cp_error.NewSysError(err)
					} else if mdW == nil {
						return false, nil, cp_error.NewNormalError(fmt.Sprintf("仓管管理的仓库ID不存在:%d", v.WarehouseID))
					}

					si.VendorDetail[0].WarehouseDetail[i].Name = mdW.Name
					si.VendorDetail[0].WarehouseDetail[i].Role = mdW.Role
				}
			}

			ml, err := NewLineDAL(si).ListLine(&cbd.ListLineReqCBD{
				VendorID: mdAcc.VendorID,
				WarehouseIDList: whIntList,
			})
			if err != nil {
				return false, nil, err
			}

			lineList, ok := ml.Items.(*[]cbd.ListLineRespCBD)
			if !ok {
				return false, nil, cp_error.NewSysError("数据转换失败")
			}

			for _, v := range *lineList {
				si.VendorDetail[0].LineDetail = append(si.VendorDetail[0].LineDetail, cp_api.LineDetail{
					LineID: v.ID,
					Source: v.Source,
					To: v.To,
					SourceWhr: v.SourceWhr,
					ToWhr: v.ToWhr,
				})
			}
		}

		return false, si, nil
	} else {
		err = cp_obj.Cjson.Unmarshal([]byte(data), si)
		if err != nil {
			return false, nil, cp_error.NewSysError("session信息解析失败:" + err.Error())
		}
		return true, si, nil
	}
}

func (this *SessionDAL) SetSessionCache(s *cp_api.CheckSessionInfo) error {
	err := this.Build()
	if err != nil {
		return cp_error.NewSysError(err)
	}
	defer this.Close()

	bd, err := cp_obj.Cjson.Marshal(s)
	if err != nil {
		return cp_error.NewSysError("[JSON]session信息编码失败:" + err.Error())
	}

	err = this.Cache.Put(cp_constant.REDIS_KEY_SESSIONKEY + s.SessionKey, string(bd), time.Minute * cp_constant.REDIS_EXPIRE_SESSION_KEY)
	if err != nil {
		return cp_error.NewSysError("[Redis]Put失败:" + err.Error())
	}

	return nil
}

func (this *SessionDAL) OnlineList(id uint64) ([]model.SSOLoginMD, error) {
	err := this.Build()
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}
	defer this.Close()

	list, err := this.DBListOnlineModule(id)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (this *SessionDAL) UpdateLastActiveDate(in *cbd.SessionReqCBD) (int64, error) {
	err := this.Build()
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer this.Close()

	return this.DBUpdateLastActiveDate(in)
}

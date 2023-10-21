package bll

import (
	"fmt"
	"strconv"
	"time"
	"warehouse/v5-go-api-shopee/bll/shopeeAPI"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-api-shopee/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_util"
)

//接口业务逻辑层
type ShopBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewShopBL(ic cp_app.IController) *ShopBL {
	if ic == nil {
		return &ShopBL{}
	}
	return &ShopBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *ShopBL) AuthShop(in *cbd.AuthShopReqCBD) (string, string, error) {
	in.SpecialID = cp_util.NewGuid()

	err := dal.NewShopDAL(this.Si).CacheAuthShop(in)
	if err != nil {
		return "", "", err
	}

	if in.Platform == constant.PLATFORM_SHOPEE {
		queryUrl, specialID := shopeeAPI.Auth.AuthShop(in.SpecialID, in.Host)
		return queryUrl, specialID, nil
	}

	return "", "", cp_error.NewSysError("请输入支持的平台代码")
}


func (this *ShopBL) ShopeeBinding(in *cbd.BindingShopReqCBD) error {
	idList := new([]uint64)

	userInfo, err := dal.NewShopDAL(this.Si).GetCacheAuthShop(in)
	if err != nil {
		return err
	}

	if in.MainAccountID != 0 { // CNSC类型
		tokenInfo, err := shopeeAPI.Auth.GetAccessToken(in.Code, "", in.MainAccountID)
		if err != nil {
			return err
		}

		if len(tokenInfo.ShopIDList) == 0 {
			return cp_error.NewNormalError("该主账号下没有店铺")
		}

		idList, err =  this.shopeeBindingMainAccountID(in, userInfo, tokenInfo)
	} else { // 普通店铺类型
		tokenInfo, err := shopeeAPI.Auth.GetAccessToken(in.Code, in.PlatformShopID, 0)
		if err != nil {
			return err
		}

		id, err := this.shopeeBindingShop(in, userInfo, tokenInfo)
		idList = &[]uint64{id}
	}

	shopList := make([]cbd.ShopDetail, 0)

	for _, v := range *idList {
		shopList = append(shopList, cbd.ShopDetail{ID: v, Platform: constant.PLATFORM_SHOPEE})
	}

	_ = NewOrderBL(this.Ic).ProducerSyncOrder(&cbd.SyncOrderReqCBD{
		SellerID: userInfo.SellerID,
		From: time.Now().AddDate(0, -1, 0).Unix(),
		To: time.Now().Unix(),
		ShopDetail: shopList,
	})

	_ = NewItemBL(this.Ic).ProducerSyncItemAndModel(&cbd.SyncShopReqCBD{
		SellerID: userInfo.SellerID,
		ShopDetail: shopList,
	})

	return nil
}

func (this *ShopBL) shopeeBindingMainAccountID(in *cbd.BindingShopReqCBD, userInfo *cbd.AuthShopReqCBD, tokenInfo *cbd.GetAccessTokenRespCBD) (*[]uint64, error) {
	success := 0
	fail := 0
	refresh := tokenInfo.RefreshToken
	idList := make([]uint64, 0)

	for _, v := range tokenInfo.ShopIDList {
		in.PlatformShopID = strconv.FormatUint(v, 10)

		refTokenInfo, err := shopeeAPI.Auth.RefreshAccessToken(refresh, in.PlatformShopID)
		if err != nil {
			fail ++
			continue
		}

		tokenInfo.AccessToken = refTokenInfo.AccessToken
		tokenInfo.RefreshToken = refTokenInfo.RefreshToken
		tokenInfo.ExpireIn = refTokenInfo.ExpireIn

		id, err := this.shopeeBindingShop(in, userInfo, tokenInfo)
		if err != nil {
			cp_log.Error(err.Error())
			fail ++
			continue
		}
		success ++
		idList = append(idList, id)
	}

	if fail > 0 {
		return nil, cp_error.NewNormalError(fmt.Sprintf("绑定成功店铺数:%d, 绑定失败店铺数:%d", success, fail))
	}

	return &idList, nil
}

func (this *ShopBL) shopeeBindingShop(in *cbd.BindingShopReqCBD, userInfo *cbd.AuthShopReqCBD, tokenInfo *cbd.GetAccessTokenRespCBD) (uint64, error) {
	var shopID uint64

	md, err := dal.NewShopDAL(this.Si).GetModelByPlatformShopID(userInfo.Platform, in.PlatformShopID)
	if err != nil {
		return 0, err
	} else if md != nil && md.SellerID != userInfo.SellerID {
		return 0, cp_error.NewNormalError(fmt.Sprintf("授权失败, id:%s的店铺已被其他卖家授权", in.PlatformShopID))
	}

	respInfo, _ := shopeeAPI.Shop.GetShopInfo(in.PlatformShopID, tokenInfo.AccessToken)
	respProf, _ := shopeeAPI.Shop.GetProfile(in.PlatformShopID, tokenInfo.AccessToken)

	if md == nil {
		mdInsert := &cbd.AddShopReqCBD{
			PlatformShopID: in.PlatformShopID,
			Platform: userInfo.Platform,
			SellerID: userInfo.SellerID,
			AccessToken: tokenInfo.AccessToken,
			RefreshToken: tokenInfo.RefreshToken,
			AccessExpire: time.Unix(time.Now().Unix() + tokenInfo.ExpireIn, 0),
			RefreshExpire: time.Now().Add(30*24*time.Hour), //30 days
		}

		if respInfo != nil {
			mdInsert.ShopExpire = time.Unix(respInfo.ExpireTime, 0)
			mdInsert.Name = respInfo.ShopName
			mdInsert.Status = respInfo.Status
			mdInsert.Region = respInfo.Region

			if respInfo.IsCB {
				mdInsert.IsCB = 1
			} else {
				mdInsert.IsCB = 0
			}

			if respInfo.IsCNSC {
				mdInsert.IsCNSC = 1
			} else {
				mdInsert.IsCNSC = 0
			}

			if respInfo.IsSIP {
				mdInsert.IsSIP = 1
			} else {
				mdInsert.IsSIP = 0
			}
		}

		if respProf != nil {
			mdInsert.Logo = respProf.Response.ShopLogo
			mdInsert.Description = respProf.Response.Description
		}

		shopID, err = dal.NewShopDAL(this.Si).AddShop(mdInsert)
		if err != nil {
			return 0, err
		}
	} else {
		shopID = md.ID
		mdEdit := &cbd.EditShopReqCBD{
			ID: md.ID,
			AccessToken: tokenInfo.AccessToken,
			RefreshToken: tokenInfo.RefreshToken,
			AccessExpire: time.Unix(time.Now().Unix() + tokenInfo.ExpireIn, 0),
			RefreshExpire: time.Now().Add(30*24*time.Hour), //30 days
			Name: md.Name,
			Region: md.Region,
			Logo: md.Logo,
			Description: md.Description,
			IsCNSC: md.IsCNSC,
			IsCB: md.IsCB,
		}

		if respInfo != nil {
			mdEdit.ShopExpire = time.Unix(respInfo.ExpireTime, 0)
			mdEdit.Name = respInfo.ShopName
			mdEdit.Status = respInfo.Status
			mdEdit.Region = respInfo.Region

			if respInfo.IsCB {
				mdEdit.IsCB = 1
			} else {
				mdEdit.IsCB = 0
			}

			if respInfo.IsCNSC {
				mdEdit.IsCNSC = 1
			} else {
				mdEdit.IsCNSC = 0
			}
		}

		if respProf != nil {
			mdEdit.Logo = respProf.Response.ShopLogo
			mdEdit.Description = respProf.Response.Description
		}

		err = NewShopBL(this.Ic).EditShop(mdEdit)
		if err != nil {
			return 0, err
		}
	}

	return shopID, nil
}

func (this *ShopBL) SyncShop(in *cbd.SyncShopReqCBD) error {
	var warnMessage string

	for _, v := range in.ShopDetail {
		md, err := dal.NewShopDAL(this.Si).GetModelByID(v.ID)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewSysError("无此店铺")
		} else if md.RefreshExpire.Before(time.Now().Add(24 * time.Hour)) {//增加24小时的容错误差
			return cp_error.NewNormalError("请重新授权", cp_constant.RESPONSE_CODE_REAUTH_SHOP)
		} else if md.AccessExpire.Before(time.Now().Add(10 * time.Minute)) {//增加10分钟的容错误差
			//access_token已过期，重新申请
			refreshResp, err := shopeeAPI.Auth.RefreshAccessToken(md.RefreshToken, md.PlatformShopID)
			if err != nil {
				return err
			}

			md.AccessToken = refreshResp.AccessToken
			md.RefreshToken = refreshResp.RefreshToken
			md.AccessExpire = time.Unix(time.Now().Unix() + refreshResp.ExpireIn, 0)
			md.RefreshExpire = time.Now().Add(30*24*time.Hour) //30 days
		}

		resp, err := shopeeAPI.Shop.GetShopInfo(md.PlatformShopID, md.AccessToken)
		if err != nil {
			return err
		}

		if resp.Status != "NORMAL" {
			warnMessage += "店铺已失效"
		} else if resp.ExpireTime < time.Now().Unix() {
			return cp_error.NewNormalError("店铺已过期")
		}

		mdEdit := &cbd.EditShopReqCBD{
			ID: md.ID,
			Name: resp.ShopName,
			AccessToken: md.AccessToken,
			RefreshToken: md.RefreshToken,
			AccessExpire: md.AccessExpire,
			ShopExpire: time.Unix(resp.ExpireTime, 0),
			RefreshExpire: md.RefreshExpire,
			Status: resp.Status,
			Region: resp.Region,
			Logo: md.Logo,
			Description: md.Description,
		}

		if !resp.IsSIP { //sip店铺 没权限获取该接口
			respProf, err := shopeeAPI.Shop.GetProfile(md.PlatformShopID, md.AccessToken)
			if err != nil {
				return err
			}
			mdEdit.Logo = respProf.Response.ShopLogo
			mdEdit.Description = respProf.Response.Description
		}


		if resp.IsCB {
			mdEdit.IsCB = 1
		} else {
			mdEdit.IsCB = 0
		}

		if resp.IsCNSC {
			mdEdit.IsCNSC = 1
		} else {
			mdEdit.IsCNSC = 0
		}

		if resp.IsSIP {
			mdEdit.IsSIP = 1
		} else {
			mdEdit.IsSIP = 0
		}

		err = NewShopBL(this.Ic).EditShop(mdEdit)
		if err != nil {
			return err
		}

		cp_log.Info("sync shop success, shop_id:" + strconv.FormatUint(v.ID, 10))
	}

	return nil
}

func (this *ShopBL) EditShop(in *cbd.EditShopReqCBD) error {
	_, err := dal.NewShopDAL(this.Si).EditShop(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ShopBL) Refresh(id uint64, refreshInfo *cbd.RefreshAccessTokenRespCBD) error {
	r := &cbd.RefreshShopReqCBD{}

	r.ID = id
	r.AccessToken = refreshInfo.AccessToken
	r.RefreshToken = refreshInfo.RefreshToken
	r.AccessExpire = time.Unix(time.Now().Unix() + refreshInfo.ExpireIn, 0)
	r.RefreshExpire = time.Now().Add(30*24*time.Hour) //30 days

	_, err := dal.NewShopDAL(this.Si).RefreshShop(r)
	if err != nil {
		return err
	}

	return nil
}

package bll

import (
	"fmt"
	"runtime"
	"strconv"
	"warehouse/v5-go-api-cangboss/bll/aliYunAPI"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

//接口业务逻辑层
type ItemBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewItemBL(ic cp_app.IController) *ItemBL {
	if ic == nil {
		return &ItemBL{}
	}
	return &ItemBL{Ic: ic, Si: ic.GetBase().Si}
}

func (this *ItemBL) AddItem(in *cbd.AddItemReqCBD) error {
	var err error

	if in.ShopID > 0 {
		md, err := dal.NewShopDAL(this.Si).GetModelByID(in.ShopID)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError(fmt.Sprintf("门店不存在:%d", in.ShopID))
		}

		in.PlatformShopID = md.PlatformShopID
	}

	for i, v := range in.Detail {
		if runtime.GOOS == "linux" {
			err = cp_util.DirMidirWhenNotExist(`/tmp/cangboss/`)
			if err != nil {
				return err
			}
			v.TmpPath = `/tmp/cangboss/` + cp_util.NewGuid() + `_` + v.Image.Filename
		} else {
			v.TmpPath = "E:\\go\\first_project\\src\\warehouse\\v5-go-api-cangboss\\" + v.Image.Filename
		}

		//先存本地临时目录
		err = this.Ic.GetBase().Ctx.SaveUploadedFile(v.Image, v.TmpPath)
		if err != nil {
			return cp_error.NewSysError("图片保存失败:" + err.Error())
		}

		//再上传图片到oss
		in.Detail[i].Url, err = aliYunAPI.Oss.UploadImage(constant.BUCKET_NAME_PUBLICE_IMAGE, constant.OSS_PATH_ITEM_PICTURE, v.Image.Filename, v.TmpPath)
		if err != nil {
			return err
		}
		in.Detail[i].PlatformModelID = "CQ" + strconv.FormatUint(cp_util.NodeSnow.NextVal(), 10)
	}

	in.Platform = constant.PLATFORM_MANUAL
	in.PlatformItemID = "CQ" + strconv.FormatUint(cp_util.NodeSnow.NextVal(), 10)

	_, err = dal.NewItemDAL(this.Si).AddItem(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ItemBL) ReportAddItem(addList []cbd.AddItemReqCBD) ([]cbd.AddItemRespCBD, error) {
	respList := make([]cbd.AddItemRespCBD, len(addList))

	for i, v := range addList {
		if v.ShopID > 0 {
			md, err := dal.NewShopDAL(this.Si).GetModelByID(v.ShopID)
			if err != nil {
				return nil, err
			} else if md == nil {
				return nil, cp_error.NewNormalError(fmt.Sprintf("门店不存在:%d", v.ShopID))
			}

			v.PlatformShopID = md.PlatformShopID
		}

		mdItem, err := dal.NewItemDAL(this.Si).GetModelByPlatformID(v.Platform, v.PlatformItemID, v.SellerID)
		if err != nil {
			return nil, err
		} else if mdItem == nil {
			resp, err := dal.NewItemDAL(this.Si).AddItem(&v)
			if err != nil {
				return nil, err
			}
			respList[i] = *resp
		} else {
			if mdItem.Name != v.Name {
				_, err = dal.NewItemDAL(this.Si).EditItem(&cbd.EditItemReqCBD{
					SellerID: v.SellerID,
					Platform: v.Platform,
					ID: mdItem.ID,
					Name: v.Name,
					ItemSku: v.ItemSku,
				})
				if err != nil {
					return nil, err
				}
			}

			resp := cbd.AddItemRespCBD {
				ItemID: mdItem.ID,
				PlatformItemID: mdItem.PlatformItemID,
			}

			for _, vv := range v.Detail {
				mdModel, err := dal.NewModelDAL(this.Si).GetModelByPlatformID(v.Platform, vv.PlatformModelID, v.SellerID)
				if err != nil {
					return nil, err
				} else if mdModel == nil {
					skuList := []cbd.ModelImageDetailCBD{vv}
					err = dal.NewModelDAL(this.Si).AddModelList(v.SellerID, v.ShopID, mdItem.ID, v.Platform, v.PlatformShopID, v.PlatformItemID, &skuList)
					if err != nil {
						return nil, err
					}
					resp.Detail = append(resp.Detail, skuList[0])
				} else {
					if mdModel.ModelSku != vv.Sku || mdModel.Images != vv.Url {
						_, err = dal.NewModelDAL(this.Si).EditModel(&cbd.EditModelReqCBD{
							SellerID: v.SellerID,
							ID: mdModel.ID,
							ModelSku: vv.Sku,
							Url: vv.Url,
						})
						if err != nil {
							return nil, err
						}
					}

					resp.Detail = append(resp.Detail, cbd.ModelImageDetailCBD{ModelID: mdModel.ID, PlatformModelID: mdModel.PlatformModelID, Sku: vv.Sku})
				}
			}

			respList[i] = resp
		}
	}

	return respList, nil
}

func (this *ItemBL) EditItem(in *cbd.EditItemReqCBD) error {
	md, err := dal.NewItemDAL(this.Si).GetModelByID(in.ID, in.SellerID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("商品不存在:" + strconv.FormatUint(in.ID, 10))
	} else if md.Platform != constant.PLATFORM_MANUAL {
		return cp_error.NewNormalError("无法编辑非手动创建的商品:" + strconv.FormatUint(in.ID, 10))
	}

	_, err = dal.NewItemDAL(this.Si).EditItem(in)
	if err != nil {
		return err
	}

	return nil
}

func (this *ItemBL) DelItem(in *cbd.DelItemReqCBD) error {
	//for _, id := range in.IDList {
	//	md, err := dal.NewItemDAL(this.Si).GetModelByID(id)
	//	if err != nil {
	//		return err
	//	} else if md == nil {
	//		return cp_error.NewNormalError("商品不存在:", id)
	//	} else if md.IsManual == 0 {
	//		return cp_error.NewNormalError("无法删除非手动创建的商品:", id)
	//	}
	//
	//	//todo
	//	//对应的model是否有库存，有库存不能删
	//
	//	//todo
	//	//删除对应的model
	//
	//	//todo
	//	//删除对应的图片
	//
	//	_, err = dal.NewItemDAL(this.Si).DelItem(in)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (this *ItemBL) ListItemAndModelSeller(in *cbd.ListItemAndModelSellerCBD) (*cp_orm.ModelList, error) {
	//todo 根据seller_key，搜索子用户

	ml, err := dal.NewItemDAL(this.Si).ListItemAndModelSeller(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

package bll

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"warehouse/v5-go-api-cangboss/bll/aliYunAPI"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/constant"
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-component/cp_api"
	"warehouse/v5-go-component/cp_app"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
	"warehouse/v5-go-component/cp_util"
)

//接口业务逻辑层
type ModelBL struct {
	Si *cp_api.CheckSessionInfo
	Ic cp_app.IController
}

func NewModelBL(ic cp_app.IController) *ModelBL {
	if ic == nil {
		return &ModelBL{}
	}
	return &ModelBL{Ic:ic, Si: ic.GetBase().Si}
}

func (this *ModelBL) AddModel(in *cbd.AddModelReqCBD) error {
	var err error

	md, err := dal.NewItemDAL(this.Si).GetModelByID(in.ItemID, in.SellerID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError(fmt.Sprintf("商品不存在:%d", in.ItemID))
	} else if md.Platform != constant.PLATFORM_MANUAL {
		return cp_error.NewNormalError("非自定义商品, 无法添加sku")
	} else {
		in.ShopID = md.ShopID
		in.PlatformShopID = md.PlatformShopID
		in.PlatformItemID = md.PlatformItemID
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
			return cp_error.NewNormalError("图片保存失败:" + err.Error())
		}

		//再上传图片到oss
		in.Detail[i].Url, err = aliYunAPI.Oss.UploadImage(constant.BUCKET_NAME_PUBLICE_IMAGE, constant.OSS_PATH_ITEM_PICTURE, v.Image.Filename, v.TmpPath)
		if err != nil {
			return err
		}
		in.Detail[i].PlatformModelID = "CQ" + strconv.FormatUint(cp_util.NodeSnow.NextVal(), 10)
	}

	err = dal.NewModelDAL(this.Si).AddModelList(in.SellerID, in.ShopID, md.ID, md.Platform, in.PlatformShopID, md.PlatformItemID, &in.Detail)
	if err != nil {
		return err
	}

	return nil
}

func (this *ModelBL) EditModel(in *cbd.EditModelReqCBD) error {
	md, err := dal.NewModelDAL(this.Si).GetModelByID(in.ID, in.SellerID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("sku不存在")
	}  else if md.Platform != constant.PLATFORM_MANUAL {
		return cp_error.NewNormalError("无法编辑非手动创建的sku商品:", in.ID)
	}

	if in.ImageChange {
		if runtime.GOOS == "linux" {
			err = cp_util.DirMidirWhenNotExist(`/tmp/cangboss/`)
			if err != nil {
				return err
			}
			in.TmpPath = `/tmp/cangboss/` + cp_util.NewGuid() + `_` + in.Image.Filename
		} else {
			in.TmpPath = "E:\\go\\first_project\\src\\warehouse\\v5-go-api-cangboss\\" + in.Image.Filename
		}

		err := this.Ic.GetBase().Ctx.SaveUploadedFile(in.Image, in.TmpPath)
		if err != nil {
			return cp_error.NewNormalError("图片保存失败:" + err.Error())
		}

		//先上传图片到oss
		in.Url, err = aliYunAPI.Oss.UploadImage(constant.BUCKET_NAME_PUBLICE_IMAGE, constant.OSS_PATH_ITEM_PICTURE, in.Image.Filename, in.TmpPath)
		if err != nil {
			return err
		}
	} else {
		in.Url = md.Images
	}

	_, err = dal.NewModelDAL(this.Si).EditModel(in)
	if err != nil {
		return err
	}

	if in.ImageChange {
		err = aliYunAPI.Oss.DeleteOSSImage(md.Images)
		if err != nil {
			cp_log.Warning("oss删除失败:" + err.Error())
		}
	}

	return nil
}

func (this *ModelBL) DelModel(in *cbd.DelModelReqCBD) error {
	md, err := dal.NewModelDAL(this.Si).GetModelByID(in.ID, in.SellerID)
	if err != nil {
		return err
	} else if md == nil {
		return cp_error.NewNormalError("sku不存在")
	} else if md.Platform != constant.PLATFORM_MANUAL {
		return cp_error.NewNormalError("无法删除非自定义sku")
	}

	mdStm, err := dal.NewModelStockDAL(this.Si).GetModelByModelID(md.ID, 0)
	if err != nil {
		return err
	}  else if mdStm != nil && mdStm.StockID > 0 {
		return cp_error.NewNormalError("该商品已绑定库存，无法删除")
	}

	err = dal.NewModelDAL(this.Si).DelModel(in, md.ItemID, md.Images)
	if err != nil {
		return err
	}

	return nil
}

func (this *ModelBL) ListModel(in *cbd.ListModelReqCBD) (*cp_orm.ModelList, error) {
	ml, err := dal.NewModelDAL(this.Si).ListModel(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

func (this *ModelBL) BindGift(in *cbd.BindGiftReqCBD) error {
	mdSource, err := dal.NewModelDAL(this.Si).GetModelByID(in.ModelID, in.SellerID)
	if err != nil {
		return err
	} else if mdSource == nil {
		return cp_error.NewNormalError("sku不存在")
	} else if mdSource.SellerID != in.SellerID {
		return cp_error.NewNormalError("sku不属于本用户")
	}

	//查看已经绑定的
	ml, err := dal.NewGiftDAL(this.Si).ListGift(&cbd.ListGiftReqCBD{SellerID: in.SellerID, ModelIDStrList: []string{strconv.FormatUint(in.ModelID, 10)}})
	if err != nil {
		return err
	}

	boundGiftList, ok := ml.Items.(*[]cbd.ListGiftRespCBD)
	if !ok {
		return cp_error.NewNormalError("数据转换失败")
	}

	for _, v := range in.ModelIDList {
		modelID, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}

		found := false
		for _, vv := range *boundGiftList {
			if modelID == vv.ModelID {
				found = true
			}
		}
		if found { //已经绑定，跳过
			continue
		}

		mdTo, err := dal.NewModelDAL(this.Si).GetModelByID(modelID, in.SellerID)
		if err != nil {
			return err
		} else if mdTo == nil {
			return cp_error.NewNormalError("SKU ID不存在")
		} else if mdTo.SellerID != in.SellerID {
			return cp_error.NewNormalError("sku不属于本用户")
		} else if mdSource.ShopID != 0 && mdTo.ShopID != 0 && mdSource.ShopID != mdTo.ShopID {
			return cp_error.NewNormalError("不能跨店铺组合商品")
		}

		err = dal.NewGiftDAL(this.Si).AddGift(&cbd.AddGiftReqCBD{
			SellerID: in.SellerID,
			SourceShopID: mdSource.ShopID,
			SourcePlatformShopID: mdSource.PlatformShopID,
			SourceItemID: mdSource.ItemID,
			SourcePlatformItemID: mdSource.PlatformItemID,
			SourceModelID: mdSource.ID,
			SourcePlatformModelID: mdSource.PlatformModelID,
			ToShopID: mdTo.ShopID,
			ToPlatformShopID: mdTo.PlatformShopID,
			ToItemID: mdTo.ItemID,
			ToPlatformItemID: mdTo.PlatformItemID,
			ToModelID: mdTo.ID,
			ToPlatformModelID: mdTo.PlatformModelID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *ModelBL) SetAutoImport(in *cbd.SetAutoImportCBD) error {
	_, err := dal.NewModelDAL(this.Si).SetAutoImport(&cbd.SetAutoImportCBD{SellerID: in.SellerID, ModelID: in.ModelID, AutoImport: in.AutoImport})
	if err != nil {
		return err
	}

	return nil
}

func (this *ModelBL) UnBindGift(in *cbd.UnBindGiftReqCBD) error {
	for _, v := range in.ModelIDList {
		modelID, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}

		md, err := dal.NewGiftDAL(this.Si).GetModelIDAndModelID(in.ModelID, modelID)
		if err != nil {
			return err
		} else if md == nil {
			return cp_error.NewNormalError("赠品绑定关系不存在")
		}

		_, err = dal.NewGiftDAL(this.Si).DelGift(&cbd.DelGiftReqCBD{ID: md.ID})
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *ModelBL) ListGift(in *cbd.ListGiftReqCBD) (*cp_orm.ModelList, error) {
	in.ModelIDStrList = strings.Split(in.ModelIDList, ";")

	ml, err := dal.NewGiftDAL(this.Si).ListGift(in)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

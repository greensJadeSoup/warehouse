package shopeeAPI

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
)

type ModelBLL struct{}

var Model ModelBLL

func (this *ModelBLL) GetModelList(platformShopID string, token string, syncItemList *[]cbd.ItemBaseInfoCBD) (*[]cbd.ItemModelListCBD, error) {
	common := CommonParam()

	ItemModelListCBD := make([]cbd.ItemModelListCBD, len(*syncItemList))
	detailItem := cbd.ModelBaseCBD{}

	cp_log.Info("going to sync mode sku list from shopee API...")

	for i, v := range *syncItemList {
		queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_MODEL_LIST, common, platformShopID, token)
		queryUrl += fmt.Sprintf(`&item_id=%s`, v.ItemID)

		//cp_log.Info("[ShopeeAPI]GetModelList请求URL:" + queryUrl)

		req, _ := http.NewRequest("GET", queryUrl, nil)
		data, err := cp_util.Do(req)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}
		cp_log.Info("[ShopeeAPI]GetModelList请求URL:" + queryUrl)

		//v.ItemID == "4765900314" || v.ItemID == "3841716440" ||
		//if v.ItemID == "16864943203" || v.ItemID == "21000897126" {
		//	cp_log.Info(string(data))
		//}

		resp := &cbd.GetModelListRespCBD{}
		err = cp_obj.Cjson.Unmarshal(data, resp)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		} else if resp.Error != "" {
			return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", resp.Error, resp.Message))
		}

		ItemModelListCBD[i].ID = v.ID
		ItemModelListCBD[i].PlatformItemID = v.ItemID

		if len(resp.Response.Model) == 0 {
			detailItem.ModelID = v.ItemID
			detailItem.ModelSku = v.ItemName + v.ItemSku
			detailItem.Images = strings.Join(v.ImageUrlList, ";")
			ItemModelListCBD[i].Model = append(ItemModelListCBD[i].Model, detailItem)
		} else {
			for _, vv := range resp.Response.Model {
				detailItem.ModelID = strconv.FormatUint(vv.ModelID, 10)

				if len(vv.TierIndex) == 1 {
					detailItem.ModelSku = resp.Response.TierVariation[0].OptionList[vv.TierIndex[0]].Option
				} else if len(vv.TierIndex) == 2 {
					detailItem.ModelSku = resp.Response.TierVariation[0].OptionList[vv.TierIndex[0]].Option +
						resp.Response.TierVariation[1].OptionList[vv.TierIndex[1]].Option
				} else if len(vv.TierIndex) == 3 {
					detailItem.ModelSku = resp.Response.TierVariation[0].OptionList[vv.TierIndex[0]].Option +
						resp.Response.TierVariation[1].OptionList[vv.TierIndex[1]].Option +
						resp.Response.TierVariation[2].OptionList[vv.TierIndex[2]].Option
				} else if len(vv.TierIndex) == 4 {
					detailItem.ModelSku = resp.Response.TierVariation[0].OptionList[vv.TierIndex[0]].Option +
						resp.Response.TierVariation[1].OptionList[vv.TierIndex[1]].Option +
						resp.Response.TierVariation[2].OptionList[vv.TierIndex[2]].Option +
						resp.Response.TierVariation[3].OptionList[vv.TierIndex[3]].Option
				}

				detailItem.ModelSku = strings.ReplaceAll(detailItem.ModelSku, `"`, "")
				detailItem.ModelSku = strings.ReplaceAll(detailItem.ModelSku, `\`, "")

				detailItem.Images = resp.Response.TierVariation[0].OptionList[vv.TierIndex[0]].Image.ImageUrl

				ItemModelListCBD[i].Model = append(ItemModelListCBD[i].Model, detailItem)
			}
		}
	}

	return &ItemModelListCBD, nil
}


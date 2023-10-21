package shopeeAPI

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"warehouse/v5-go-api-shopee/cbd"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_obj"
	"warehouse/v5-go-component/cp_util"
)

type ItemBLL struct{}

var Item ItemBLL

func (this *ItemBLL) GetItemList(platformShopID string, token string, updateTime int64) (*cbd.GetItemListRespCBD, error) {
	common := CommonParam()

	offset := 0
	pageSize := 100
	resp := &cbd.GetItemListRespCBD{}

	for {
		queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_ITEM_LIST, common, platformShopID, token)
		queryUrl += fmt.Sprintf(`&update_time_from=%[1]d&update_time_to=%[2]d&item_status=%[3]s&offset=%[4]d&page_size=%[5]d`,
			//updateTime, time.Now().Unix(), "NORMAL", "BANNED", "DELETED", "UNLIST", offset, pageSize)
			updateTime, time.Now().Unix(), "NORMAL", offset, pageSize)
		cp_log.Info("[ShopeeAPI]GetItemList请求URL:" + queryUrl)

		req, _ := http.NewRequest("GET", queryUrl, nil)
		data, err := cp_util.Do(req)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}
		cp_log.Debug(string(data))

		subResp := &cbd.GetItemListRespCBD{}
		err = cp_obj.Cjson.Unmarshal(data, subResp)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		} else if subResp.Error != "" {
			cp_log.Info(string(data))
			return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", subResp.Error, subResp.Message))
		}

		for _, v := range subResp.Response.Item {
			resp.Response.Item = append(resp.Response.Item, v)
		}

		if subResp.Response.HasNext {
			offset = subResp.Response.NextOffset
			continue
		} else {
			break
		}
	}

	return resp, nil
}

func (this *ItemBLL) GetItemBaseInfo(platformShopID string, token string, itemIDStrList []string) (*[]cbd.ItemBaseInfoCBD, error) {
	var err error

	common := CommonParam()

	offset := 0
	idx := 0
	count := len(itemIDStrList)

	if count > 50 {
		offset = 50
	} else {
		offset = count
	}

	itemBaseInfoList := make([]cbd.ItemBaseInfoCBD, 0)

	for {
		queryUrl := GenerateShopQueryUrl(constant.SHOPEE_URI_GET_ITEM_BASE_INFO, common, platformShopID, token)
		queryUrl += fmt.Sprintf(`&item_id_list=%s`,
			strings.Join(itemIDStrList[idx:offset], ","))

		cp_log.Info("[ShopeeAPI]GetItemBaseInfo请求URL:" + queryUrl)

		data := make([]byte, 0)
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("GET", queryUrl, nil)
			data, err = cp_util.Do(req)
			if err != nil {
				continue
			}
		}
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}
		cp_log.Info(string(data))

		subResp := &cbd.GetItemBaseInfoRespCBD{}
		err = cp_obj.Cjson.Unmarshal(data, subResp)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		} else if subResp.Error != "" {
			return nil, cp_error.NewSysError(fmt.Sprintf("[ShopeeAPI]: Error:%s. Message:%s", subResp.Error, subResp.Message))
		}

		if len(subResp.Response.ItemList) != offset - idx {
			return nil, cp_error.NewSysError("商品明细同步数量不一致")
		}

		for _, v := range subResp.Response.ItemList {
			item := cbd.ItemBaseInfoCBD{
				ItemID: strconv.FormatUint(v.ItemID, 10),
				CategoryID: strconv.FormatUint(v.CategoryID, 10),
				ItemName: v.ItemName,
				ItemStatus: v.ItemStatus,
				Description: v.Description,
				ItemSku: v.ItemSku,
				Weight: v.Weight,
				HasModel: v.HasModel,
				UpdateTime: v.UpdateTime,
				ImageUrlList: v.Image.ImageUrlList,
			}

			item.ItemName = strings.ReplaceAll(item.ItemName, `"`, "")
			item.ItemName = strings.ReplaceAll(item.ItemName, `\`, "")

			if v.HasModel {
				item.IntHasModel = 1
			} else {
				item.IntHasModel = 0
			}

			item.WeightFloat, _ = strconv.ParseFloat(v.Weight, 5)

			itemBaseInfoList = append(itemBaseInfoList, item)
		}

		if count - len(subResp.Response.ItemList) > 0 {
			idx = offset
			count -= len(subResp.Response.ItemList)

			if count > 50 {
				offset += 50
			} else {
				offset += count
			}

			continue
		} else {
			break
		}
	}

	return &itemBaseInfoList, nil
}

package whlog

import (
	"warehouse/v5-go-api-cangboss/dal"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_log"
)

var WhLogArray []model.WarehouseLogMD

func FlushWarehouseLog()  {
	err := dal.NewWarehouseLogDAL(nil).FlushWarehouseLog(&WhLogArray)
	if err != nil {
		cp_log.Error(err.Error())
		return
	}
}
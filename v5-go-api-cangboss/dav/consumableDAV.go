package dav

import (
	"fmt"
	"warehouse/v5-go-api-cangboss/cbd"
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_cache"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_orm"
)

//基本数据层
type ConsumableDAV struct {
	cp_orm.DA
	Cache cp_cache.ICache
}

func (this *ConsumableDAV) Build() error {
	this.Cache = cp_cache.GetCache()
	return cp_orm.InitDA(this, model.NewConsumable())
}

func (this *ConsumableDAV) DBGetModelByID(id uint64) (*model.ConsumableMD, error) {
	md := model.NewConsumable()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,name,note FROM %s WHERE id=%d`, md.TableName(), id)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ConsumableDAV][DBGetModelByID]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ConsumableDAV) DBGetModelByName(vendorID uint64, name string) (*model.ConsumableMD, error) {
	md := model.NewConsumable()

	searchSQL := fmt.Sprintf(`SELECT id,vendor_id,name,note 
			FROM %[1]s WHERE vendor_id=%[2]d and name='%[3]s'`,
			md.TableName(), vendorID, name)

	cp_log.Debug(searchSQL)
	hasRow, err := this.SQL(searchSQL).Get(md)
	if err != nil {
		return nil, cp_error.NewSysError("[ConsumableDAV][DBGetModelByName]:" + err.Error())
	} else if !hasRow {
		return nil, nil
	}

	return md, nil
}

func (this *ConsumableDAV) DBInsert(md interface{}) error  {
	execRow, err := this.Session.Insert(md)
	if err != nil {
		return cp_error.NewNormalError(err)
	} else if execRow == 0 {
		return cp_error.NewNormalError("[ConsumableDAV][DBInsert]失败,系统繁忙")
	}

	return nil
}

func (this *ConsumableDAV) DBListConsumable(in *cbd.ListConsumableReqCBD) (*cp_orm.ModelList, error) {
	var condSQL string


	searchSQL := fmt.Sprintf(`SELECT c.id,c.vendor_id,c.name,c.note
		FROM %[1]s c
		WHERE c.vendor_id=%[2]d%[3]s`,
		this.GetModel().TableName(), in.VendorID, condSQL)

	return this.MysqlModelList(searchSQL, in.IsPaging, in.PageIndex, in.PageSize, &[]cbd.ListConsumableRespCBD{})
}

func (this *ConsumableDAV) DBUpdateConsumable(md *model.ConsumableMD) (int64, error) {
	row, err := this.Session.ID(md.ID).Update(md)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	return row, nil
}

func (this *ConsumableDAV) DBDelConsumable(in *cbd.DelConsumableReqCBD) (int64, error) {
	md := model.NewConsumable()
	md.ID = in.ID

	execRow, err := this.Session.Delete(md)
	if err != nil {
		return 0, cp_error.NewNormalError(err)
	}

	return execRow, nil
}
